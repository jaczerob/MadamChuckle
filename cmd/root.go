package cmd

import (
	"github.com/spf13/cobra"
)

var (
	verbose bool

	rootCmd = &cobra.Command{
		Use:   "madamchuckle",
		Short: "a Discord Bot made for Toontown Rewritten Discord servers",
		Long: `A Discord Bot made for Toontown Rewritten Discord servers, based off 
		my favorite Toontown shopkeep, Madam Chuckle!`,
	}
)

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initCfg)
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "set trace logging output")
}
