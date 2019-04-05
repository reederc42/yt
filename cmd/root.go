package cmd

import (
	"github.com/reederc42/yt/yt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io"
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
	rootCmd.Flags().StringP("query", "q", "", "document query")
	rootCmd.Flags().BoolP("silence-usage", "s", false,
		"silences usage on error")
}

func rootCmdEntry(cmd *cobra.Command, _ []string) error {
	if viper.GetBool("silence-usage") {
		cmd.SilenceUsage = true
		cmd.SilenceErrors = true
	}

	i, err := getInput(viper.GetString("input"))
	if err != nil {
		return err
	}
	o, err := getOutput(viper.GetString("output"))
	if err != nil {
		return err
	}
	v, err := yt.Compile(i)
	if err != nil {
		return err
	}
	query := viper.GetString("query")
	if query != "" {
		v, err = yt.Query(v, query)
		if err != nil {
			return err
		}
	}
	err = yt.Write(v, o)
	return err
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil && rootCmd.SilenceErrors {
		println("error: " + err.Error())
	}
}

func viperBindFlags(cmd *cobra.Command, _ []string) error {
	return viper.BindPFlags(cmd.Flags())
}

func getInput(input string) (io.Reader, error) {
	if input == os.Stdin.Name() {
		return os.Stdin, nil
	}

	return os.Open(input)
}

func getOutput(output string) (io.Writer, error) {
	if output == os.Stdout.Name() {
		return os.Stdout, nil
	}

	return os.Create(output)
}
