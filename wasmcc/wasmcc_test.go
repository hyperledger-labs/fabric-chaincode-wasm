package main

import (
	"encoding/hex"

	//	"encoding/hex"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"io/ioutil"
	"os"
	"strconv"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Tests for wasmcc simple asset transfer", func() {

	status200 := int32(200)
	var payload []byte
	var account1InitBal []byte
	var account2InitBal []byte

	stub := shim.NewMockStub("testingStub", new(WASMChaincode))
	BeforeSuite(func() {
		stub.MockInit("000", nil)
	})

	Describe("Sample asset wasm chaincode on wasmcc", func() {
		Context("wasm chaincode installation", func() {
			It("wasm binary file chaincode creation should be success", func() {
				result := stub.MockInvoke("000",
					[][]byte{[]byte("create"),
						[]byte("balancewasm-wasm"),
						[]byte(ReadAssetTransferWASM()),
						[]byte("account1"),
						[]byte("100"),
						[]byte("account2"),
						[]byte("1000")})
				Expect(result.Status).Should(Equal(status200))
			})
			It("wasm zip file chaincode creation should be success", func() {
				result := stub.MockInvoke("000",
					[][]byte{[]byte("create"),
						[]byte("balancewasm-zip"),
						[]byte(ReadAssetTransferWASMZip()),
						[]byte("account1"),
						[]byte("100"),
						[]byte("account2"),
						[]byte("1000")})
				Expect(result.Status).Should(Equal(status200))
			})
			It("wasm hex encoded chaincode creation should be success", func() {
				result := stub.MockInvoke("000",
					[][]byte{[]byte("create"),
						[]byte("balancewasm"),
						[]byte(ReadAssetTransferWASMHex()),
						[]byte("account1"),
						[]byte("100"),
						[]byte("account2"),
						[]byte("1000")})
				Expect(result.Status).Should(Equal(status200))
			})
			Specify("WASM chaincode should be installed successfully", func() {
				result := stub.MockInvoke("000",
					[][]byte{[]byte("installedChaincodes")})
				payload = []byte(result.Payload)
				Expect(payload).Should(Equal([]byte("balancewasm\nbalancewasm-wasm\nbalancewasm-zip\n")))
			})
		})
		Context("account1 is created with some balance", func() {
			It("account1 should exist", func() {
				result := stub.MockInvoke("000",
					[][]byte{[]byte("execute"),
						[]byte("balancewasm"),
						[]byte("query"),
						[]byte("account1")})

				account1InitBal = []byte(result.Payload)
				Expect(result.Status).Should(Equal(status200))
			})
			It("account1 balance should be same as getState", func() {
				account1Bal, _ := stub.GetState("balancewasm_account1")
				Expect(string(account1Bal)).Should(Equal(string(account1InitBal)))
			})
			Specify("account1 balance should be 100", func() {
				Expect(string("100")).Should(Equal(string(account1InitBal)))
			})
		})
		Context("account2 is created with some balance", func() {
			It("account2 should exist", func() {
				result := stub.MockInvoke("000",
					[][]byte{[]byte("execute"),
						[]byte("balancewasm"),
						[]byte("query"),
						[]byte("account2")})

				account2InitBal = []byte(result.Payload)
				Expect(result.Status).Should(Equal(status200))
			})
			It("account2 balance should be same as getState", func() {
				account2Bal, _ := stub.GetState("balancewasm_account2")
				Expect(string(account2Bal)).Should(Equal(string(account2InitBal)))
			})
			Specify("account2 balance should be 1000", func() {
				Expect(string("1000")).Should(Equal(string(account2InitBal)))
			})
		})
		Context("transfer 10 units from account2 to account2", func() {
			It("transfer should be successful", func() {
				result := stub.MockInvoke("000",
					[][]byte{[]byte("execute"),
						[]byte("balancewasm"),
						[]byte("invoke"),
						[]byte("account2"),
						[]byte("account1"),
						[]byte("10")}).Status
				Expect(result).Should(Equal(status200))
			})
			Specify("Account 1 and account 2 balance should be updated to new balance", func() {
				result := stub.MockInvoke("000",
					[][]byte{[]byte("execute"),
						[]byte("balancewasm"),
						[]byte("query"),
						[]byte("account1")})
				payload = []byte(result.Payload)

				newBalAcc1, _ := strconv.Atoi(string(payload))

				newExpectedBalAcc1, _ := strconv.Atoi(string(account1InitBal))
				newExpectedBalAcc1 = newExpectedBalAcc1 + 10
				Expect(newExpectedBalAcc1).Should(Equal(newBalAcc1))

				result = stub.MockInvoke("000",
					[][]byte{[]byte("execute"),
						[]byte("balancewasm"),
						[]byte("query"),
						[]byte("account2")})
				payload = []byte(result.Payload)

				newBalAcc2, _ := strconv.Atoi(string(payload))

				newExpectedBalAcc2, _ := strconv.Atoi(string(account2InitBal))
				newExpectedBalAcc2 = newExpectedBalAcc2 - 10
				Expect(newExpectedBalAcc2).Should(Equal(newBalAcc2))
			})
		})
	})
})

func ReadAssetTransferWASMZip() []byte {

	file, err := ioutil.ReadFile("../sample-wasm-chaincode/chaincode_example02/rust/app_main.zip")
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	//encodedFile := hex.EncodeToString(file)
	return file
}

func ReadAssetTransferWASM() []byte {

	file, err := ioutil.ReadFile("../sample-wasm-chaincode/chaincode_example02/rust/app_main.wasm")
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	//encodedFile := hex.EncodeToString(file)
	return file
}

func ReadAssetTransferWASMHex() []byte {

	file, err := ioutil.ReadFile("../sample-wasm-chaincode/chaincode_example02/rust/app_main.wasm")
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	encodedFile := hex.EncodeToString(file)
	return []byte(encodedFile)
}
