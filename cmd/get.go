/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"secret/vault"

	"github.com/spf13/cobra"
)

// setCmd represents the set command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Gets a secret in your secret storage 📍",
	Run: func(cmd *cobra.Command, args []string) {
		v := vault.FileVault(encodingKey, SecretPath())
		key := args[0]
		value, err := v.GetKeyInFileVault(key)
		if err != nil {
			panic(err)
		}
		fmt.Printf("%s =  %s\n", key, value)
	},
}

func init() {
	RootCmd.AddCommand(getCmd)
}
