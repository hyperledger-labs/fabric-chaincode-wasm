package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"strconv"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"

	"github.com/perlin-network/life/exec"
	"github.com/perlin-network/life/wasm-validation"
)

const CHAINCODE_EXISTS = "{\"code\":101, \"reason\": \"chaincode exists with same name\"}"
const UKNOWN_ERROR = "{\"code\":301, \"reason\": \"uknown error : %s\"}"

type WASMChaincode struct {
}

// Resolver defines imports for WebAssembly modules ran in Life.
type Resolver struct {
	tempRet0 int64
}

var chaincodeNameGlobal string
var stubGlobal shim.ChaincodeStubInterface

// ResolveFunc defines a set of import functions that may be called within a WebAssembly module.
func (r *Resolver) ResolveFunc(module, field string) exec.FunctionImport {
	fmt.Printf("Resolve func: %s %s\n", module, field)
	switch module {
	case "env":
		switch field {
		case "__get_state":
			return func(vm *exec.VirtualMachine) int64 {
				ptr := int(uint32(vm.GetCurrentFrame().Locals[0]))
				msgLen := int(uint32(vm.GetCurrentFrame().Locals[1]))
				msg := vm.Memory[ptr : ptr+msgLen]
				fmt.Printf("[app] getState fn called with msg: %s\n", string(msg))

				s := fmt.Sprintf("%s_%s",chaincodeNameGlobal,msg)

				valueFromState, _ := stubGlobal.GetState(s)
				if valueFromState == nil {
					return -1
				}


				fmt.Printf("[app] getState fn response is: %s\n", string(valueFromState))


				//How to return the result of get state to wasm here?
				return 0
			}
		case "__put_state":
			return func(vm *exec.VirtualMachine) int64 {
				keyPtr := int(uint32(vm.GetCurrentFrame().Locals[0]))
				keyMsgLen := int(uint32(vm.GetCurrentFrame().Locals[1]))
				key := vm.Memory[keyPtr : keyPtr+keyMsgLen]

				valuePtr := int(uint32(vm.GetCurrentFrame().Locals[2]))
				valueMsgLen := int(uint32(vm.GetCurrentFrame().Locals[3]))
				value := vm.Memory[valuePtr : valuePtr+valueMsgLen]

				fmt.Printf("[app] putState fn called with key: %s and value: %s\n", string(key), string(value))


				s := fmt.Sprintf("%s_%s",chaincodeNameGlobal,key)
				// Store the chaincode in ledger
				err := stubGlobal.PutState(s, value)
				if err != nil {
					fmt.Printf(UKNOWN_ERROR, err.Error())
					return -1
				}
				return 0
			}
		default:
			panic(fmt.Errorf("unknown field: %s", field))
		}
	default:
		panic(fmt.Errorf("unknown module: %s", module))
	}
}

// ResolveGlobal defines a set of global variables for use within a WebAssembly module.
func (r *Resolver) ResolveGlobal(module, field string) int64 {
	fmt.Printf("Resolve global: %s %s\n", module, field)
	switch module {
	case "env":
		switch field {
		case "__constant_variable":
			return 424
		default:
			panic(fmt.Errorf("unknown field: %s", field))
		}
	default:
		panic(fmt.Errorf("unknown module: %s", module))
	}
}
func (t *WASMChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("Init invoked")
	return shim.Success(nil)
}

func (t *WASMChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("Invoke function")
	function, args := stub.GetFunctionAndParameters()
	fmt.Printf("Invoke function %s with args %v", function, args)

	if function == "create" {
		// Create a new wasm chaincode
		return t.create(stub, args)
	} else if function == "query" {
		// query a wasm chaincode
		return t.query(stub, args)
	}

	return shim.Error("Invalid invoke function name. Expecting \"invoke\" \"create\" \"query\"")
}

// Store a new wasm chaincode in state. Receives chaincode name and wasm file encoded in hex
func (t *WASMChaincode) create(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("Create function")
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	ChaincodeName := args[0]

	chaincodeNameGlobal = ChaincodeName
	stubGlobal = stub

	fmt.Println("Installing wasm chaincode: " + ChaincodeName)

	//check if same chaincode name is already present
	chaincodeFromState, err := stub.GetState(ChaincodeName)
	if chaincodeFromState != nil {
		return shim.Error(CHAINCODE_EXISTS)
	}

	//Decode the chaincode
	ChaincodeHexEncoded := args[1]
	fmt.Println("Encoded wasm chaincode: " + ChaincodeHexEncoded)
	ChaincodeDecoded, err := hex.DecodeString(ChaincodeHexEncoded)


	result := runWASM(ChaincodeDecoded, "init");

	fmt.Printf("Init Response:%s\n", result)

	// Store the chaincode in ledger
	err = stub.PutState(ChaincodeName, ChaincodeDecoded)
	if err != nil {
		s := fmt.Sprintf(UKNOWN_ERROR, err.Error())
		return shim.Error(s)
	}
	return shim.Success([]byte("Success! Installed wasm chaincode"))
}

// query callback representing the query of a chaincode
func (t *WASMChaincode) query(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting chaincode name to query")
	}

	ChaincodeName := args[0]
	funcToInvoke := args[1]
	chaincodeNameGlobal = ChaincodeName
	stubGlobal = stub

	// Get the state from the ledger
	Chaincodebytes, _ := stub.GetState(ChaincodeName)
	if Chaincodebytes == nil {
		jsonResp := "{\"Error\":\"No Chaincode for " + ChaincodeName + "\"}"
		return shim.Error(jsonResp)
	}

	//How to send any parameter to query function?
	result := runWASM(Chaincodebytes, funcToInvoke);



	fmt.Printf("Query Response:%s\n", result)
	return shim.Success([]byte(strconv.FormatInt(result, 10)))
}

// query callback representing the query of a chaincode
func runWASM(Chaincodebytes []byte, funcToInvoke string) int64 {

	//entryFunctionFlag := flag.String("entry", funcToInvoke, "entry function name")
	//noFloatingPointFlag := flag.Bool("no-fp", false, "disable floating point")
	flag.Parse()

	validator, err := wasm_validation.NewValidator()
	if err != nil {
		panic(err)
	}
	err = validator.ValidateWasm(Chaincodebytes)
	if err != nil {
		panic(err)
	}

	// Instantiate a new WebAssembly VM with a few resolved imports.
	vm, err := exec.NewVirtualMachine(Chaincodebytes, exec.VMConfig{
		DefaultMemoryPages:   128,
		DefaultTableSize:     65536,
		DisableFloatingPoint: false,
	}, new(Resolver), nil)

	if err != nil {
		panic(err)
	}

	// Get the function ID of the entry function to be executed.
	entryID, ok := vm.GetFunctionExport(funcToInvoke)
	if !ok {
		fmt.Printf("Entry function %s not found; starting from 0.\n", funcToInvoke)
		entryID = 0
	}

	start := time.Now()

	// Run the WebAssembly chaincode's entry function.
	result, err := vm.Run(entryID)
	if err != nil {
		vm.PrintStackTrace()
		panic(err)
	}
	end := time.Now()

	fmt.Printf("return value = %d, duration = %v\n", result, end.Sub(start))

	return result
}

func main() {
	err := shim.Start(new(WASMChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
