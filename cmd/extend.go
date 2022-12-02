package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"liarslie/pkg/client"
	"os"
)

// extendCmd represents the start command
var extendCmd = &cobra.Command{
	Use:   "extend",
	Short: "Adds a new agents to the network",
	Long:  `Recreates the network with a new Agent`,
	Run: func(cmd *cobra.Command, args []string) {
		value, err := cmd.Flags().GetInt("value")
		if err != nil {
			fmt.Printf("Cannot Extend Network. Missing: %s", err.Error())
			os.Exit(1)
		}
		maxValue, err := cmd.Flags().GetInt("max-value")
		if err != nil {
			fmt.Printf("Cannot Extend Network. Missing: %s", err.Error())
			os.Exit(1)
		}
		numAgents, err := cmd.Flags().GetInt("num-agents")
		if err != nil {
			fmt.Printf("Cannot Extend Network. Missing: %s", err.Error())
			os.Exit(1)
		}
		liarRatio, err := cmd.Flags().GetFloat32("liar-ratio")
		if err != nil {
			fmt.Printf("Cannot Extend Network. Missing: %s", err.Error())
			os.Exit(1)
		}
		client.ExtendNetwork(value, maxValue, numAgents, liarRatio)
	},
}

func init() {
	RootCmd.AddCommand(extendCmd)
	extendCmd.PersistentFlags().Int("value", 5, "The real value for the game.")
	extendCmd.PersistentFlags().Int("max-value", 100, "The max value for liar agents.")
	extendCmd.PersistentFlags().Int("num-agents", 10, "The max number of agents to spawn.")
	extendCmd.PersistentFlags().Float32("liar-ratio", 0.33, "The % ratio of liars in the game.")
}
