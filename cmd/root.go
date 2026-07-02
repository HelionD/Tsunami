/*
Copyright © 2026 NAME HERE tsunami.project.dev@gmail.com
*/
package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string

	rootCmd = &cobra.Command{
		Use:   "Tsunami",
		Short: "TSUNAMI",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return initializeConfig(cmd)
		},
	}
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default locations: ., $HOME/.Tsunami/)")
}

func initializeConfig(cmd *cobra.Command) error {
	// 1. Set up Viper to use environment variables.
	viper.SetEnvPrefix("TSUNAMI")
	// Allow for nested keys in environment variables (e.g. `TSUNAMI_DATABASE_HOST`)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "*", "-", "*"))
	viper.AutomaticEnv()

	// 2. Handle the configuration file.
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Search for a config file in default locations.
		home, err := os.UserHomeDir()
		// Only panic if we can't get the home directory.
		cobra.CheckErr(err)

		// Search for a config file with the name "config" (without extension).
		viper.AddConfigPath(".")
		viper.AddConfigPath(home + "/.Tsunami")
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
	}

	// 3. Read the configuration file.
	// If a config file is found, read it in. We use a robust error check
	// to ignore "file not found" errors, but panic on any other error.
	if err := viper.ReadInConfig(); err != nil {
		// It's okay if the config file doesn't exist.
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if !errors.As(err, &configFileNotFoundError) {
			return err
		}
	}

	// 4. Bind Cobra flags to Viper.
	// This is the magic that makes the flag values available through Viper.
	// It binds the full flag set of the command passed in.
	err := viper.BindPFlags(cmd.Flags())
	if err != nil {
		return err
	}

	// This is an optional but useful step to debug your config.
	fmt.Println("Configuration initialized. Using config file:", viper.ConfigFileUsed())
	return nil
}
