/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/spf13/cobra"
	"liarslie/pkg/client"
)

// stopCmd represents the start command
var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stops the game of liars lie.",
	Long:  `Stops the game and exits.`,
	Run: func(cmd *cobra.Command, args []string) {

		client.StopClient()
	},
}

func init() {
	RootCmd.AddCommand(stopCmd)

}
