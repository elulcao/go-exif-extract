package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var cfgFile string
var vpr = viper.New()

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:          "go-exif-extractor",
	Short:        "Extracts the exif data from images in a directory and subdirectories",
	Long:         `Extracts the exif data from images in a directory and subdirectories`,
	SilenceUsage: true,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.AddCommand(extractCmd)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config.yaml)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		vpr.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".askGPT" (without extension).
		vpr.AddConfigPath(".")
		vpr.AddConfigPath(home)
		vpr.SetConfigType("yaml")
		vpr.SetConfigName(".config")
	}

	vpr.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := vpr.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "using config file:", vpr.ConfigFileUsed())

		// Add the config file flags to the cobra commands.
		viperToCobraFlags(rootCmd, vpr)
		viperToCobraFlags(extractCmd, vpr)
	}
}

// viperToCobraFlags copies the values from viper to cobra flags.
func viperToCobraFlags(cmd *cobra.Command, vpr *viper.Viper) {
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		if vpr.IsSet(f.Name) {
			_ = cmd.Flags().Set(f.Name, vpr.GetString(f.Name))
		}
	})
}
