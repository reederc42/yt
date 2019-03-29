package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

var rootCmd = &cobra.Command{
	Use: "yt",
	Short: "yt implements inheritance and components in YAML",
	PreRunE: viperBindFlags,
	RunE: rootCmdEntry,
}

func init() {
	rootCmd.Flags().StringP("input", "i", os.Stdin.Name(), "input file")
	rootCmd.Flags().StringP("output", "o", os.Stdout.Name(), "output file")
}

func rootCmdEntry(cmd *cobra.Command, _ []string) error {
	return nil
}

func Execute() {
	_ = rootCmd.Execute()
}

func viperBindFlags(cmd *cobra.Command, _ []string) error {
	return viper.BindPFlags(cmd.Flags())
}
