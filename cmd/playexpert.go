/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"liarslie/pkg/client"
	"log"
)

// playExpertCmd represents the start command
var playExpertCmd = &cobra.Command{
	Use:   "playexpert",
	Short: "Plays a round of liarslie in expert mode.",
	Long:  `Queries the agents to determines the real value V.`,
	Run: func(cmd *cobra.Command, args []string) {
		flags := cmd.Flags()
		numAgents, err := flags.GetInt("num-agents")
		if err != nil {
			fmt.Printf("Cannot Play expert mode: %s", err.Error())
			return
		}
		liarRatio, err := cmd.Flags().GetFloat32("liar-ratio")
		if err != nil {
			fmt.Printf("Cannot Play expert mode: %s", err.Error())
			return
		}
		_, err = client.PlayExpert(client.CurrentClient.AgentsFullNetwork, numAgents, liarRatio)
		if err != nil {
			log.Fatal(err.Error())
		}
	},
}

func init() {
	RootCmd.AddCommand(playExpertCmd)
	playExpertCmd.PersistentFlags().Int("num-agents", 10, "The max number of agents to spawn.")
	playExpertCmd.PersistentFlags().Float32("liar-ratio", 0.33, "The % ratio of liars in the game.")
}
