package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

//nolint:exhaustivestruct,gochecknoglobals
var rootCmd = &cobra.Command{
	Use:   "tlagent",
	Short: "Timelapse Agent",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
