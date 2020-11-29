package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var topicCmd = &cobra.Command{
	Use:   "topic",
	Short: `
Run helpful topic commands
`,
}

func init() {
	topicCmd.AddCommand(listTopics)
	rootCmd.AddCommand(topicCmd)
}

var listTopics = &cobra.Command{
	Use:   "list",
	Short: "Displays all topics for a cluster",
	RunE: func(cmd *cobra.Command, args []string) error {
		topics, err := kafka.ListTopics()
		 if err != nil {
		 	return err
		 }
		 for _, topic := range topics {
		 	fmt.Println(topic)
		 }
		return nil
	},
}