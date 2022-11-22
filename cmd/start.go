/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"liarslie/pkg/client"
	"os"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Starts the game of liars lie.",
	Long:  `Spawns the number of agents and sets waits for commands.`,
	Run: func(cmd *cobra.Command, args []string) {
		value, err := cmd.Flags().GetInt("value")
		if err != nil {
			fmt.Sprintf("Cannot Start Agents. Missing: %s", err.Error())
			os.Exit(1)
		}
		maxValue, err := cmd.Flags().GetInt("max-value")
		if err != nil {
			fmt.Sprintf("Cannot Start Agents. Missing: %s", err.Error())
			os.Exit(1)
		}
		numAgents, err := cmd.Flags().GetInt("num-agents")
		if err != nil {
			fmt.Sprintf("Cannot Start Agents. Missing: %s", err.Error())
			os.Exit(1)
		}
		liarRatio, err := cmd.Flags().GetFloat32("liar-ratio")
		if err != nil {
			fmt.Sprintf("Cannot Start Agents. Missing: %s", err.Error())
			os.Exit(1)
		}
		client.StartClient(cmd, value, maxValue, numAgents, liarRatio)
	},
}

func init() {
	RootCmd.AddCommand(startCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	startCmd.PersistentFlags().Int("value", 5, "The real value for the game.")
	startCmd.PersistentFlags().Int("max-value", 100, "The max value for liar agents.")
	startCmd.PersistentFlags().Int("num-agents", 10, "The max number of agents to spawn.")
	startCmd.PersistentFlags().Float32("liar-ratio", 0.33, "The % ratio of liars in the game.")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// startCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
