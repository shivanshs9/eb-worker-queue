package cmd

import (
	"os"

	app "github.com/shivanshs9/eb-worker-queue/pkg"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "sqsd",
	Short: "Start the SQS worker, polling and forwarding the messages via HTTP requests",
	Run:   runWorker,
}

var debug bool
var options *app.AppOptions

func init() {
	options = new(app.AppOptions)
	rootCmd.PersistentFlags().BoolVarP(&debug, "verbose", "v", false, "verbose logging")
	rootCmd.Flags().StringVarP(&options.QueueUrl, "queueUrl", "q", "", "Provide the queue URL (required)")
	rootCmd.MarkFlagRequired("queueUrl")

	rootCmd.Flags().IntVarP(&options.MaxBufferedMessages, "maxJobs", "m", 10, "Provide the limit of messages to receive (max. 10)")
	rootCmd.Flags().StringVarP(&options.DefaultHttpPath, "httpPath", "p", "/", "Provide the HTTP Path of the API to hit with POST request of the job")

	rootCmd.Flags().StringVarP(&options.ApiHost, "host", "a", "http://localhost:80", "Provide the Host on which API is listening")
}

func runWorker(cmd *cobra.Command, args []string) {
	log := logrus.New()
	log.Info("Verbose logging enabled")
	if debug {
		log.SetLevel(logrus.DebugLevel)
	}
	if options.MaxBufferedMessages > 10 {
		options.MaxBufferedMessages = 10
	}
	app.StartApp(options, log)
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
