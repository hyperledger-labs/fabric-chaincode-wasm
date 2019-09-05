package main

import (
	"flag"
	"encoding/hex"
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
				return 0
			}
		case "__put_state":
			return func(vm *exec.VirtualMachine) int64 {
				ptr := int(uint32(vm.GetCurrentFrame().Locals[0]))
				msgLen := int(uint32(vm.GetCurrentFrame().Locals[1]))
				msg := vm.Memory[ptr : ptr+msgLen]
				fmt.Printf("[app] putState fn called with msg: %s\n", string(msg))
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

	fmt.Println("Installing wasm chaincode: "+ChaincodeName)

	//check if same chaincode name is already present
	chaincodeFromState,err := stub.GetState(ChaincodeName)
	if chaincodeFromState != nil {
		return shim.Error(CHAINCODE_EXISTS)
	}

	//Decode the chaincode
	ChaincodeHexEncoded := args[1]
	fmt.Println("Encoded wasm chaincode: "+ChaincodeHexEncoded)
	ChaincodeDecoded, err := hex.DecodeString(ChaincodeHexEncoded)

	// Store the chaincode in ledger
	err = stub.PutState(ChaincodeName, ChaincodeDecoded)
	if err != nil {
		s := fmt.Sprintf(UKNOWN_ERROR, err.Error())
		return shim.Error(s)
	}
	fmt.Println("Success! Installed wasm chaincode")

	return shim.Success(nil)
}

// query callback representing the query of a chaincode
func (t *WASMChaincode) query(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting chaincode name to query")
	}

	A := args[0]
	funcToInvoke := args[1]

	// Get the state from the ledger
	Chaincodebytes, _ := stub.GetState(A)
	if Chaincodebytes == nil {
		jsonResp := "{\"Error\":\"No Chaincode for " + A + "\"}"
		return shim.Error(jsonResp)
	}

	entryFunctionFlag := flag.String("entry", funcToInvoke, "entry function name")
	noFloatingPointFlag := flag.Bool("no-fp", false, "disable floating point")
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
		DisableFloatingPoint: *noFloatingPointFlag,
	}, new(Resolver), nil)

	if err != nil {
		panic(err)
	}

	// Get the function ID of the entry function to be executed.
	entryID, ok := vm.GetFunctionExport(*entryFunctionFlag)
	if !ok {
		fmt.Printf("Entry function %s not found; starting from 0.\n", *entryFunctionFlag)
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

	fmt.Printf("Query Response:%s\n", result)
	return shim.Success([]byte(strconv.FormatInt(result, 10)))
}




func main() {
	err := shim.Start(new(WASMChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
