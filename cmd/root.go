package cmd

import (
	"github.com/reederc42/yt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

var rootCmd = &cobra.Command{
	Use: "yt",
	Short: "yt implements inheritance and components in YAML",
	PreRunE: viperBindFlags,
	RunE: rootCmdEntry,
	Args: cobra.MaximumNArgs(1),
}

func init() {
	rootCmd.Flags().StringP("input", "i", os.Stdin.Name(), "input file")
	rootCmd.Flags().StringP("output", "o", os.Stdout.Name(), "output file")
	rootCmd.Flags().StringP("query", "q", "",
		"document query (overwrites argument)")
	rootCmd.Flags().BoolP("silence-usage", "s", false,
		"silences usage on error")
	rootCmd.Flags().BoolP("json", "j", false, "output as JSON")
}

func rootCmdEntry(cmd *cobra.Command, args []string) error {
	if viper.GetBool("silence-usage") {
		cmd.SilenceUsage = true
		cmd.SilenceErrors = true
	}

	i, d, err := getInput(viper.GetString("input"))
	if err != nil {
		return err
	}
	o, err := getOutput(viper.GetString("output"))
	if err != nil {
		return err
	}
	writeDir, _ := os.Getwd()
	_ = os.Chdir(d)
	iValue, err := ioutil.ReadAll(i)
	if err != nil {
		return err
	}
	v, err := yt.Compile(iValue)
	if err != nil {
		return err
	}
	var qry string
	if len(args) > 0 {
		qry = args[0]
	}
	if q := viper.GetString("query"); q != "" {
		qry = q
	}
	if qry != "" {
		v, err = yt.Query(v, qry)
		if err != nil {
			return err
		}
	}
	_ = os.Chdir(writeDir)
	if !viper.GetBool("json") {
		err = yt.WriteYAML(v, o)
	} else {
		err = yt.WriteJSON(v, o)
	}
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

func getInput(input string) (io.Reader, string, error) {
	if input == os.Stdin.Name() {
		return os.Stdin, "", nil
	}

	i, err := os.Open(input)
	dir := filepath.Dir(input)
	return i, dir, err
}

func getOutput(output string) (io.Writer, error) {
	if output == os.Stdout.Name() {
		return os.Stdout, nil
	}

	return os.Create(output)
}
