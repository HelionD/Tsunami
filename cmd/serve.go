/*
Copyright © 2026 NAME HERE tsunami.project.dev@gmail.com

*/
package cmd

import (
	"fmt"
	"github.com/spf13/viper"
	"github.com/spf13/cobra"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Starts the server",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		port := viper.GetInt("port")
		fmt.Printf("Serve command called on port %d\n", port)
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	serveCmd.Flags().Int("port", 8080, "Port to run the server on")

	viper.BindPFlag("port", serveCmd.Flags().Lookup("port"))

}
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")}
