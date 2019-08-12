/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/ledger"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/pkg/errors"
)

func main() {
	sdk, err := fabsdk.New(config.FromFile("./first-network.yaml"))
	if err != nil {
		fmt.Println(errors.WithMessage(err, "failed to create SDK"))
		os.Exit(-1)
	}
	defer sdk.Close()

	user := "User1"
	org := "Org1"
	channelName := "mychannel"
	chainCodeID := "wasmcc"

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

	fmt.Println("Enter wasm file path : ")
	var filepath string
	fmt.Scanf("%s", &filepath)

	fmt.Println("Trying to read file : " + filepath)

	file, err := ioutil.ReadFile(filepath)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	fmt.Printf("file length : %d \n", len(file))
	encodedFile := hex.EncodeToString(file)
	fmt.Printf("encoded file length : %d \n", len(encodedFile))

	//Get chaincode name from user
	fmt.Println("Please enter wasm chaincode name : ")
	var wasmChaincodeName string
	fmt.Scanf("%s", &wasmChaincodeName)

	//Get init parameters from user
	fmt.Println("Please enter all parameters for your init function separated by new line. Hit enter twice for exit: ")

	var args [][]byte

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		input := scanner.Text()

		if input == ""{
			break
		}

		args = append(args, []byte(input))
	}


	fmt.Println("Executing transaction")

	response, err := client.Execute(channel.Request{
		ChaincodeID: chainCodeID,
		Fcn:         "create",
		Args: args,
	})

	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	fmt.Println(response)
}
