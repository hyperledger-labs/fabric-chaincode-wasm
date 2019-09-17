/*
Copyright © 2019 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFile, user, org, channel, wasmfile, chaincode, args string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "wasm-pusher",
	Short: "WASM pusher is a tool used to install webassembly chaincodes in Hyperledger fabric using fabric-chaincode-wasm",
	Long: `This tool is based on fabric-go-sdk. It can be used to directly install a wasm chaincode on Hyperledger fabric.

WASMCC supports wasm chaincode in three formats i.e. .wasm binary, zip file containing wasm binary and hex encoded wasm file.
WASM pusher can install wasm chaincode in any of the above formats.`,

	RunE: func(cmd *cobra.Command, args []string) error {
		err := validateFlags()
		if err != nil {
			return err
		}
		return nil
	},
}

func validateFlags() error {

	viper.AutomaticEnv()
	flags := rootCmd.Flags()
	flags.StringVarP(&cfgFile, "configfile", "c", "./first-network.yaml", "fabric config file")
	viper.BindPFlag("configfile", flags.Lookup("configfile"))

	return nil
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.wasm-pusher.yaml)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".wasm-pusher" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".wasm-pusher")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
