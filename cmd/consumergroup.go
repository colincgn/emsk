package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
)

var consumergroupCommand = &cobra.Command{
	Use:   "consumergroup",
	Short: `list and describe consumer groups on a cluster.`,
}

func init() {
	consumergroupCommand.AddCommand(listConsumerGroups)
	rootCmd.AddCommand(consumergroupCommand)
}

var listConsumerGroups = &cobra.Command{
	Use:   "list",
	Short: "Displays all consumer groups for a cluster",
	RunE: func(cmd *cobra.Command, args []string) error {

		consumergroups, err := kafka.ListConsumerGroups()
		if err != nil {
			return err
		}
		cgroups := *consumergroups
		cgroupJSON, err := json.MarshalIndent(cgroups, "", "  ")
		fmt.Println(string(cgroupJSON))
		return nil
	},
}
