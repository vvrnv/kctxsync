package cmd

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/spf13/cobra"
	"k8s.io/client-go/tools/clientcmd"
)

var kubeconfigPathGet string

// getCmd defines the command to display Kubernetes contexts with certificate expiration in days, hours, and minutes
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "List Kubernetes contexts with certificate expiration in days, hours, and minutes",
	Run: func(cmd *cobra.Command, args []string) {
		// Set default kubeconfig path if not provided
		if kubeconfigPathGet == "" {
			kubeconfigPathGet = filepath.Join(os.Getenv("HOME"), ".kube", "config")
		}

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

		// Print sorted context names with certificate expiration
		fmt.Println("List of available Kubernetes contexts (sorted alphabetically):")
		for _, contextName := range contextNames {
			// Get the context details
			context := config.Contexts[contextName]
			userName := context.AuthInfo

			// Find the user associated with this context
			user, exists := config.AuthInfos[userName]
			if !exists || user.ClientCertificateData == nil || len(user.ClientCertificateData) == 0 {
				fmt.Printf("- %s (no certificate found)\n", contextName)
				continue
			}

			// Parse the PEM-encoded certificate
			block, _ := pem.Decode(user.ClientCertificateData)
			if block == nil || block.Type != "CERTIFICATE" {
				fmt.Printf("- %s (invalid certificate)\n", contextName)
				continue
			}

			cert, err := x509.ParseCertificate(block.Bytes)
			if err != nil {
				fmt.Printf("- %s (failed to parse certificate: %v)\n", contextName, err)
				continue
			}

			// Calculate time remaining until expiration using time.Until
			timeRemaining := time.Until(cert.NotAfter)

			// Format the time remaining as Xd Yh Zm
			days := int(timeRemaining.Hours() / 24)
			hours := int(timeRemaining.Hours()) % 24
			minutes := int(timeRemaining.Minutes()) % 60

			// Format the output
			if timeRemaining > 0 {
				fmt.Printf("- %s (certificate time expiration: %dd %dh %dm)\n", contextName, days, hours, minutes)
			} else {
				fmt.Printf("- %s (certificate expired)\n", contextName)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
	// Define the config flag (-c or --config) to pass the kubeconfig path
	getCmd.Flags().StringVarP(&kubeconfigPathGet, "config", "c", "", "Path to the kubeconfig file")
}
