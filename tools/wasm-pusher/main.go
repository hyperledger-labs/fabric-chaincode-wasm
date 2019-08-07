/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
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

	file, err := ioutil.ReadFile("file.wasm")
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	fmt.Println(len(file))
	encodedFile := hex.EncodeToString(file)
	fmt.Println(len(encodedFile))

	response, err := client.Execute(channel.Request{
		ChaincodeID: chainCodeID,
		Fcn:         "create",
		Args: [][]byte{[]byte("balancewasm2"),
			[]byte(encodedFile), // wasm chaincode
			[]byte("account1"), []byte("100"),
			[]byte("account2"), []byte("1000"),
		},
	})

	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	fmt.Println(response)
}
