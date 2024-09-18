package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var (
	Tag    string
	Os     string
	Arch   string
	Date   string
	Commit string
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version information of kctxsync",
	Long:  `All software has versions. This is kctxsync's`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("kctxsync version %s %s/%s %s %s\n", Tag, Os, Arch, Date, Commit)
	},
}
