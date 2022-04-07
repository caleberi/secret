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
var setCmd = &cobra.Command{
	Use:   "set",
	Short: "Sets a secret in your secret storage ❇️",
	Run: func(cmd *cobra.Command, args []string) {
		v := vault.FileVault(encodingKey, SecretPath())
		key, value := args[0], args[1]
		err := v.SetKeyInFileVault(key, value)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Value was successfully set 🐘")
	},
}

func init() {
	RootCmd.AddCommand(setCmd)
}
