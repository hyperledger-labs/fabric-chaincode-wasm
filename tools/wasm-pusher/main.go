/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"encoding/hex"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/ledger"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/pkg/errors"
)

var cfgFile, user, org, channelName, wasmfile, chaincode string
var args []string

func main() {

	chainCodeID := "wasmcc"

	err := readConfigVar()
	if err == nil {
		return
	}

	sdk, err := fabsdk.New(config.FromFile(cfgFile))
	if err != nil {
		fmt.Println(errors.WithMessage(err, "failed to create SDK"))
		os.Exit(-1)
	}
	defer sdk.Close()

	clientChannelContext := sdk.ChannelContext(channelName, fabsdk.WithUser(user), fabsdk.WithOrg(org))
	// client for interacting directly with the ledger
	ledger, err := ledger.New(clientChannelContext)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	bci, err := ledger.QueryInfo()
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	fmt.Println("starting block height:", bci.BCI.Height)

	// client for sending invoke
	client, err := channel.New(clientChannelContext)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	fmt.Println("Trying to read file : " + wasmfile)

	file, err := ioutil.ReadFile(wasmfile)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	fmt.Printf("file length : %d \n", len(file))
	encodedFile := hex.EncodeToString(file)
	fmt.Printf("encoded file length : %d \n", len(encodedFile))

	var txnargs [][]byte

	//Add wasm chaincode name to arguments
	txnargs = append(txnargs, []byte(chaincode))

	if args != nil {
		for _, arg := range args {
			txnargs = append(txnargs, []byte(arg))
		}
	}

	fmt.Println("Executing transaction")

	response, err := client.Execute(channel.Request{
		ChaincodeID: chainCodeID,
		Fcn:         "create",
		Args:        txnargs,
	})

	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	fmt.Println(response)
}

func readConfigVar() error {
	var cmd = &cobra.Command{
		Use:   "wasm-pusher",
		Short: "WASM pusher is a tool used to install webassembly chaincodes in Hyperledger fabric using fabric-chaincode-wasm",
		Long: `This tool is based on fabric-go-sdk. It can be used to directly install a wasm chaincode on Hyperledger fabric.
WASMCC supports wasm chaincode in three formats i.e. .wasm binary, zip file containing wasm binary and hex encoded wasm file.
WASM pusher can install wasm chaincode in any of the above formats.`,
		RunE: func(command *cobra.Command, args []string) error {
			cfgFile = viper.GetString("configfile")
			if cfgFile == "" {
				return fmt.Errorf("configfile flag is required")
			}
			user = viper.GetString("user")
			if user == "" {
				return fmt.Errorf("user flag is required")
			}
			org = viper.GetString("org")
			if org == "" {
				return fmt.Errorf("org flag is required")
			}
			channelName = viper.GetString("channelName")
			if channelName == "" {
				return fmt.Errorf("channelName flag is required")
			}
			wasmfile = viper.GetString("wasmfile")
			if wasmfile == "" {
				return fmt.Errorf("wasmfile flag is required")
			}
			chaincode = viper.GetString("chaincode")
			if chaincode == "" {
				return fmt.Errorf("chaincode flag is required")
			}
			args = viper.GetStringSlice("args")
			return nil
		},
	}

	viper.AutomaticEnv()
	flags := cmd.Flags()
	flags.StringVarP(&cfgFile, "configfile", "c", "./first-network.yaml", "fabric config file")
	viper.BindPFlag("configfile", flags.Lookup("configfile"))
	flags.StringVarP(&user, "user", "u", "User1", "User identity to use")
	viper.BindPFlag("user", flags.Lookup("user"))
	flags.StringVarP(&org, "org", "o", "Org1", "Organization of user")
	viper.BindPFlag("org", flags.Lookup("org"))
	flags.StringVarP(&channelName, "channelName", "c", "mychannel", "channel on which wasmcc is installed")
	viper.BindPFlag("channelName", flags.Lookup("channelName"))
	flags.StringVarP(&user, "wasmfile", "w", "", "wasm chaincode filepath")
	viper.BindPFlag("wasmfile", flags.Lookup("wasmfile"))
	flags.StringVarP(&chaincode, "chaincode", "cc", "", "wasm chaincode name")
	viper.BindPFlag("chaincode", flags.Lookup("chaincode"))
	flags.StringSliceVarP(&args, "args", "a", nil, "arguments for init fn of wasm chaincode")
	viper.BindPFlag("args", flags.Lookup("args"))

	return cmd.Execute()
}
