package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"log"
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
		cgroups, err := kafka.ListConsumerGroups()
		 if err != nil {
		 	return err
		 }

		 for _, cgroup := range *cgroups {
			 cgroupJSON, err := json.MarshalIndent(cgroup, "", "  ")
			 if err != nil {
				 log.Println("unable to marshall json")
				 return err
			 }

			 fmt.Printf("%s\n", string(cgroupJSON))
		 }
		return nil
	},
}