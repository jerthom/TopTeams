package cmd

import (
	"log"
	"os"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"topTeams/dota"
)

// rootCmd represents the base command
var rootCmd = &cobra.Command{
	Use:   "topTeams",
	Short: "A CLI that gets the top N pro DOTA teams",
	Long: `This CLI fetches teams for all of the ProPlayers from the opendota API. 
			It calculates an 'Experience' value for each player, defined by the number of seconds in that player's history.
			Team Experience is then calculated by the sum of all valid players' experience.
			The Top N teams are selected by TeamId, and then sorted by Team Experience before being written out to the specified file location.`,
	Run: func(cmd *cobra.Command, args []string) {
		n, _ := cmd.Flags().GetInt("numTeams")
		if n < 0 {
			log.Fatalf("Invalid value for numTeams; value must be >= zero")
		}

		teams, err := dota.TopTeams(n)
		if err != nil {
			log.Fatalf("Failed to fetch from opendota API: %v", err)
		}

		b, err := yaml.Marshal(teams)
		if err != nil {
			log.Fatalf("Failed to convert to yaml: %v", err)
		}

		o, _ := cmd.Flags().GetString("outputFile")
		err = os.WriteFile(o, b, 0644)
		if err != nil {
			log.Fatalf("Failed to write to output file: %v", err)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

// init is where we configure the flags for our cli command
func init() {
	rootCmd.Flags().IntP("numTeams", "n", 5, "number of top teams to fetch")
	rootCmd.Flags().StringP("outputFile", "o", "output.yaml", "output file location")
}
