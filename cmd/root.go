package cmd

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"text/template"

	"github.com/reederc42/yt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

var rootCmd = &cobra.Command{
	Use:   "yt [query]",
	Short: "yt implements inheritance and components in YAML",
	RunE:  rootCmdEntry,
	//Args is the first function run after parsing flags
	//  before this, cobra.OnInitialize(y ...func()) is called
	//  because this command requires arguments, flags are bound to viper here
	Args: func(cmd *cobra.Command, args []string) error {
		if err := viperBindFlags(cmd, args); err != nil {
			return err
		}
		if err := silenceCmdUsage(cmd, args); err != nil {
			return err
		}
		err := cobra.MaximumNArgs(1)(cmd, args)
		return err
	},
	Version: "v0.0.1",
}

func init() {
	rootCmd.Flags().StringP("input", "f", os.Stdin.Name(), "input file")
	rootCmd.Flags().StringP("insert", "i", "", "insert value: sets [query] as"+
		" value (YAML string)")
	rootCmd.Flags().BoolP("json", "j", false, "output as JSON")
	rootCmd.Flags().BoolP("no-compile", "n", false, "do not compile input")
	rootCmd.Flags().StringP("output", "o", os.Stdout.Name(), "output file")
	rootCmd.Flags().StringP("query", "q", "",
		"document query (overwrites query argument)")
	rootCmd.Flags().BoolP("silence-usage", "s", false,
		"silences usage on error")
	rootCmd.Flags().StringP("template", "t", "", "outputs template processed "+
		"with yt document")
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
	var documentValue interface{}
	if viper.GetBool("no-compile") {
		err := yaml.Unmarshal(iValue, &documentValue)
		if err != nil {
			return err
		}
	} else {
		v, err := yt.Compile(iValue, map[string]bool{})
		if err != nil {
			return err
		}
		documentValue = v
	}
	var qry string
	if len(args) > 0 {
		qry = args[0]
	}
	if q := viper.GetString("query"); q != "" {
		qry = q
	}
	insertRaw := viper.GetString("insert")
	if insertRaw != "" {
		var insertObj interface{}
		err = yaml.Unmarshal([]byte(insertRaw), &insertObj)
		if err != nil {
			return err
		}
		documentValue, err = yt.Insert(documentValue, insertObj, qry)
		if err != nil {
			return err
		}
	} else if qry != "" {
		documentValue, err = yt.Query(documentValue, qry, map[string]bool{})
		if err != nil {
			return err
		}
	}
	_ = os.Chdir(writeDir)
	if tplFile := viper.GetString("template"); tplFile == "" {
		if !viper.GetBool("json") {
			err = yt.WriteYAML(documentValue, o)
		} else {
			err = yt.WriteJSON(documentValue, o)
		}
	} else {
		tplRaw, err := ioutil.ReadFile(tplFile)
		if err != nil {
			return err
		}
		tpl := template.Must(template.New("").Parse(string(tplRaw)))
		err = yt.WriteTemplate(documentValue, tpl, o)
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

func silenceCmdUsage(cmd *cobra.Command, _ []string) error {
	if viper.GetBool("silence-usage") {
		cmd.SilenceUsage = true
		cmd.SilenceErrors = true
	}
	return nil
}
