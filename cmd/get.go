package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/spf13/cobra"
	"k8s.io/client-go/tools/clientcmd"
)

var kubeconfigPathGet string

// getCmd defines the command to display Kubernetes contexts
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "list k8s contexts",
	Run: func(cmd *cobra.Command, args []string) {
		// Path to the kubeconfig file (default is in the home directory)
		// Set default kubeconfig path if not provided
		if kubeconfigPathGet == "" {
			kubeconfigPathGet = filepath.Join(os.Getenv("HOME"), ".kube", "config")
		}

		//kubeconfig := filepath.Join(homedir.HomeDir(), ".kube", "config")

		// Load the configuration from the kubeconfig file
		config, err := clientcmd.LoadFromFile(kubeconfigPathGet)
		if err != nil {
			fmt.Printf("Error loading kubeconfig file: %v\n", err)
			os.Exit(1)
		}

		// Collect context names in a slice
		var contextNames []string
		for contextName := range config.Contexts {
			contextNames = append(contextNames, contextName)
		}

		// Sort the context names alphabetically
		sort.Strings(contextNames)

		// Print sorted context names
		fmt.Println("List of available Kubernetes contexts (sorted alphabetically):")
		for _, contextName := range contextNames {
			fmt.Println("- " + contextName)
		}
	},
}

func init() {
	rootCmd.AddCommand(getCmd)

	// Define the config flag (-c or --config) to pass the kubeconfig path
	getCmd.Flags().StringVarP(&kubeconfigPathGet, "config", "c", "", "Path to the kubeconfig file")
}
