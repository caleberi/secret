/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"os"
	"path/filepath"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
)

var encodingKey string

func init() {
	RootCmd.PersistentFlags().StringVarP(
		&encodingKey, "key", "k", "", "the key to use when encoding and decoding secrets")
}

func SecretPath() string {
	home, _ := homedir.Dir()
	return filepath.Join(home, ".secrets")
}

// rootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "secret",
	Short: "API Secret Book",
	Long:  `Secret is an CLI application for store api secrets `,
}

func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
