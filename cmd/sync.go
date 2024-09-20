package cmd

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api" // Importing api

	"github.com/spf13/cobra"
)

var kubeconfigPathSync string
var sshUser string
var sshHost string // New flag for specifying the SSH host
var syncAll bool   // Flag for syncing all contexts

// syncCmd defines the command to sync kubeconfig with a remote server
var syncCmd = &cobra.Command{
	Use:   "sync [context name]",
	Short: "Sync local kubeconfig with remote server",
	Args:  cobra.MaximumNArgs(1), // Maximum one argument: context name (can be omitted)
	Run: func(cmd *cobra.Command, args []string) {
		// Set default kubeconfig path if not provided
		if kubeconfigPathSync == "" {
			kubeconfigPathSync = filepath.Join(os.Getenv("HOME"), ".kube", "config")
		}

		// Load local kubeconfig
		localConfig, err := clientcmd.LoadFromFile(kubeconfigPathSync)
		if err != nil {
			fmt.Printf("Error loading local kubeconfig file from %s: %v\n", kubeconfigPathSync, err)
			os.Exit(1)
		}

		// If syncAll flag is set, sync all contexts
		if syncAll {
			for contextName := range localConfig.Contexts {
				fmt.Printf("Syncing context: %s\n", contextName)
				syncContext(contextName, localConfig)
			}
		} else {
			// If no context is provided, select the first available context
			var contextName string
			if len(args) == 0 {
				if len(localConfig.Contexts) == 1 {
					// If there is only one context, select it
					for ctx := range localConfig.Contexts {
						contextName = ctx
						break
					}
					fmt.Printf("No context provided, using the only available context: %s\n", contextName)
				} else {
					fmt.Println("Error: No context provided and there are multiple contexts in the config. Please choose a context.")
					os.Exit(1)
				}
			} else {
				contextName = args[0]
			}
			// Sync the single context
			syncContext(contextName, localConfig)
		}
	},
}

// syncContext performs the sync operation for a single context
func syncContext(contextName string, localConfig *api.Config) {
	// Find the context in the local kubeconfig
	context, ok := localConfig.Contexts[contextName]
	if !ok {
		fmt.Printf("Context '%s' not found in local kubeconfig\n", contextName)
		os.Exit(1)
	}

	// Find the cluster and user associated with the context
	clusterName := context.Cluster
	userName := context.AuthInfo

	cluster, ok := localConfig.Clusters[clusterName]
	if !ok {
		fmt.Printf("Cluster '%s' not found in local kubeconfig\n", clusterName)
		os.Exit(1)
	}

	user, ok := localConfig.AuthInfos[userName]
	if !ok {
		fmt.Printf("User '%s' not found in local kubeconfig\n", userName)
		os.Exit(1)
	}

	// Use the SSH host from the flag if provided, otherwise use the server from the kubeconfig
	serverHost := sshHost
	if serverHost == "" {
		// Extract the server URL from kubeconfig (remove https:// if exists)
		serverURL := strings.TrimPrefix(cluster.Server, "https://")
		// Remove port if present
		serverHost = strings.Split(serverURL, ":")[0]
		if serverHost == "" {
			fmt.Printf("Cluster for context '%s' does not have a valid server URL\n", contextName)
			os.Exit(1)
		}
	}

	// Get the remote kubeconfig via SSH
	fmt.Printf("Connecting to remote server %s to fetch kubeconfig...\n", serverHost)
	remoteKubeconfigData, err := getRemoteKubeconfig(serverHost)
	if err != nil {
		fmt.Printf("Failed to get remote kubeconfig: %v\n", err)
		os.Exit(1)
	}

	// Parse the remote kubeconfig
	remoteConfig, err := clientcmd.Load([]byte(remoteKubeconfigData))
	if err != nil {
		fmt.Printf("Error parsing remote kubeconfig: %v\n", err)
		os.Exit(1)
	}

	// Update certificates and keys in the local kubeconfig
	updateNeeded := false

	for _, remoteCluster := range remoteConfig.Clusters {
		if !bytes.Equal(cluster.CertificateAuthorityData, remoteCluster.CertificateAuthorityData) {
			fmt.Println("Updating certificate-authority-data...")
			cluster.CertificateAuthorityData = remoteCluster.CertificateAuthorityData
			updateNeeded = true
		}
	}

	for _, remoteUser := range remoteConfig.AuthInfos {
		if !bytes.Equal(user.ClientCertificateData, remoteUser.ClientCertificateData) {
			fmt.Println("Updating client-certificate-data...")
			user.ClientCertificateData = remoteUser.ClientCertificateData
			updateNeeded = true
		}
		if !bytes.Equal(user.ClientKeyData, remoteUser.ClientKeyData) {
			fmt.Println("Updating client-key-data...")
			user.ClientKeyData = remoteUser.ClientKeyData
			updateNeeded = true
		}
	}

	// Save the updated local kubeconfig
	if updateNeeded {
		fmt.Println("Updating local kubeconfig file...")
		err = clientcmd.WriteToFile(*localConfig, kubeconfigPathSync)
		if err != nil {
			fmt.Printf("Error writing updated kubeconfig file: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Local kubeconfig updated successfully.")
	} else {
		fmt.Println("Local kubeconfig is already up-to-date.")
	}
}

// getRemoteKubeconfig fetches the remote kubeconfig file via SSH
func getRemoteKubeconfig(serverHost string) (string, error) {
	sshCommand := exec.Command("ssh", fmt.Sprintf("%s@%s", sshUser, serverHost), "cat ~/.kube/config")

	var stdout, stderr bytes.Buffer
	sshCommand.Stdout = &stdout
	sshCommand.Stderr = &stderr

	err := sshCommand.Run()
	if err != nil {
		return "", fmt.Errorf("SSH command failed: %v, %s", err, stderr.String())
	}

	return stdout.String(), nil
}

func init() {
	rootCmd.AddCommand(syncCmd)

	// Define the config flag (-c or --config) to pass the kubeconfig path
	syncCmd.Flags().StringVarP(&kubeconfigPathSync, "config", "c", "", "Path to the kubeconfig file")
	// Define the user flag (-u or --user) to specify SSH user
	syncCmd.Flags().StringVarP(&sshUser, "user", "u", "root", "Username for SSH connection")
	// Define the ssh flag (-s or --ssh) to specify the SSH host manually
	syncCmd.Flags().StringVarP(&sshHost, "ssh", "s", "", "Specify SSH host (if not provided, will use the server from kubeconfig)")
	// Define the all flag (-a or --all) to sync all contexts
	syncCmd.Flags().BoolVarP(&syncAll, "all", "a", false, "Sync all contexts from the kubeconfig file")
}
