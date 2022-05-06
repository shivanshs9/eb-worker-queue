package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:       "sqsd queueUrl",
	Short:     "Start the SQS worker, polling and forwarding the messages via HTTP requests",
	Run:       runWorker,
	ValidArgs: []string{"queueUrl"},
	Args:      cobra.OnlyValidArgs,
}

func runWorker(cmd *cobra.Command, args []string) {
	fmt.Println(args)
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
