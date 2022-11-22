/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/spf13/cobra"
)

// playCmd represents the start command
var playCmd = &cobra.Command{
	Use:   "play",
	Short: "Plays a round of liarslie.",
	Long:  `Queries the agents to determines the real value V.`,
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func init() {
	RootCmd.AddCommand(playCmd)
}
