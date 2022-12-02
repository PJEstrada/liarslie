/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"liarslie/pkg/client"
	"os"
)

// killCmd represents the start command
var killCmd = &cobra.Command{
	Use:   "kill",
	Short: "Kill an agent from the network.",
	Long:  `Kills the agent with the given ID`,
	Run: func(cmd *cobra.Command, args []string) {
		id, err := cmd.Flags().GetString("id")
		if err != nil {
			fmt.Println(fmt.Sprintf("Cannot Kill Agent: %s", err.Error()))
			os.Exit(1)
		}
		agentID, err := uuid.Parse(id)
		if err != nil {
			fmt.Println(fmt.Sprintf("Cannot Kill agent: %s", err.Error()))
			os.Exit(1)
		}
		client.KillAgent(agentID)
	},
}

func init() {
	RootCmd.AddCommand(killCmd)
	killCmd.PersistentFlags().String("id", "", "The id of the agent to kill.")
}
