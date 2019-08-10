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
	//Global variables to be used by exported wasm functions
	chaincodeName string
	stub          shim.ChaincodeStubInterface
	args          []string
	result        []byte
}

//Index Names
var chaincodeStoreIndex = "chaincodeData"


// ResolveFunc defines a set of import functions that may be called within a WebAssembly module.
func (r *Resolver) ResolveFunc(module, field string) exec.FunctionImport {

	if debug {
		fmt.Printf("Resolve func: %s %s\n", module, field)
	}
	switch module {
	case "env":
		switch field {
		case "__print":
			return func(vm *exec.VirtualMachine) int64 {
				ptr := int(uint32(vm.GetCurrentFrame().Locals[0]))
				msgLen := int(uint32(vm.GetCurrentFrame().Locals[1]))
				msg := vm.Memory[ptr : ptr+msgLen]

				if debug {
					fmt.Printf("[app] print fn called with msg: %s\n", string(msg))
				}else{
					fmt.Print(string(msg))
				}

				return 0
			}
		case "__get_parameter":
			return func(vm *exec.VirtualMachine) int64 {
				paramNumber := int(uint32(vm.GetCurrentFrame().Locals[0]))
				ptrForResult := int(uint32(vm.GetCurrentFrame().Locals[1]))

				//Check if argument contains this many elements
				if len(r.args) < paramNumber {
					return -1
				}


				paramToReturn:=r.args[paramNumber]

				//Memory location for storing parameter
				result := vm.Memory[ptrForResult : ptrForResult+len(paramToReturn)]


				//Copying the getState parameter to above memory location
				copy(result,paramToReturn)

				if debug {
					fmt.Printf("[app] get parameter fn called with parameter number: %d , result: %s \n", paramNumber, paramToReturn)
				}
				//Returning length of parameter
				return int64(len(paramToReturn))
			}
		case "__get_state":
			return func(vm *exec.VirtualMachine) int64 {

				//Pointer and length for key
				ptr := int(uint32(vm.GetCurrentFrame().Locals[0]))
				msgLen := int(uint32(vm.GetCurrentFrame().Locals[1]))


				//Pointer for value to be returned
				ptr2 := int(uint32(vm.GetCurrentFrame().Locals[2]))

				msg := vm.Memory[ptr : ptr+msgLen]
				fmt.Printf("[app] getState fn called with msg: %s\n", string(msg))

				s := fmt.Sprintf("%s_%s",r.chaincodeName,msg)

				valueFromState, _ := r.stub.GetState(s)
				if valueFromState == nil {
					return -1
				}

				//Memory location for storing result of getState
				result := vm.Memory[ptr2 : ptr2+len(valueFromState)]

				//Copying the getState result to above memory location
				copy(result,valueFromState)

				if debug {
					fmt.Printf("[app] getState fn response is: %s\n", string(valueFromState))
				}

				//Returning length of value
				return int64(len(valueFromState))
			}
		case "__put_state":
			return func(vm *exec.VirtualMachine) int64 {

				//Pointer and length for key
				keyPtr := int(uint32(vm.GetCurrentFrame().Locals[0]))
				keyMsgLen := int(uint32(vm.GetCurrentFrame().Locals[1]))
				key := vm.Memory[keyPtr : keyPtr+keyMsgLen]


				//Pointer and length for value
				valuePtr := int(uint32(vm.GetCurrentFrame().Locals[2]))
				valueMsgLen := int(uint32(vm.GetCurrentFrame().Locals[3]))
				value := vm.Memory[valuePtr : valuePtr+valueMsgLen]

				if debug {
					fmt.Printf("[app] putState fn called with key: %s and value: %s\n", string(key), string(value))
				}

				s := fmt.Sprintf("%s_%s",r.chaincodeName,key)

				// Store the key, value in ledger
				err := r.stub.PutState(s, value)
				if err != nil {
					fmt.Printf(UKNOWN_ERROR, err.Error())
					return -1
				}
				return 0
			}
		case "__delete_state":
			return func(vm *exec.VirtualMachine) int64 {

				//Pointer and length for key
				ptr := int(uint32(vm.GetCurrentFrame().Locals[0]))
				msgLen := int(uint32(vm.GetCurrentFrame().Locals[1]))

				msg := vm.Memory[ptr : ptr+msgLen]

				if debug {
					fmt.Printf("[app] deleteState fn called with msg: %s\n", string(msg))
				}
				s := fmt.Sprintf("%s_%s",r.chaincodeName,msg)

				err := r.stub.DelState(s)
				if err != nil {
					return -1
				}

				//Returning length of value
				return 0
			}
		case "__return_result":
			return func(vm *exec.VirtualMachine) int64 {

				//Pointer and length for key
				ptr := int(uint32(vm.GetCurrentFrame().Locals[0]))
				msgLen := int(uint32(vm.GetCurrentFrame().Locals[1]))

				msg := vm.Memory[ptr : ptr+msgLen]

				if debug {
					fmt.Printf("[app] returnResult fn called with msg: %s\n", string(msg))
				}
				r.result = make([]byte, msgLen)
				copy(r.result,msg)
				//Returning length of value
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
	} else if function == "invoke" {
		// invoke a new wasm chaincode
		return t.invoke(stub, args)
	}else if function == "installedChaincodes" {
		// invoke a new wasm chaincode
		return t.installedChaincodes(stub, args)
	}

	return shim.Error("Invalid invoke function name. Expecting \"invoke\" \"create\" \"query\"")
}

func (t *WASMChaincode) installedChaincodes(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	// Get all chaincodes from the ledger
	installedChaincodeResultsIterator, err := stub.GetStateByPartialCompositeKey(chaincodeStoreIndex, []string{});
	if err != nil {
		return shim.Error(err.Error())
	}

	// Iterate through result set and get chaincode names
	var i int
	var installedChaincodeNamesList string = ""

	for i = 0; installedChaincodeResultsIterator.HasNext(); i++ {
		// Note that we don't get the value (2nd return variable), we'll just get the marble name from the composite key
		responseRange, err := installedChaincodeResultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}

		// get the color and name from color~name composite key
		objectType, compositeKeyParts, err := stub.SplitCompositeKey(responseRange.Key)
		if err != nil {
			return shim.Error(err.Error())
		}
		returnedChaincodeName := compositeKeyParts[0]
		fmt.Printf("- found a chaincode from index:%s name:%s\n", objectType, returnedChaincodeName)
		installedChaincodeNamesList += returnedChaincodeName
		installedChaincodeNamesList += "\n"

	}

	fmt.Printf("Invoke Response:%d\n", installedChaincodeNamesList)
	return shim.Success([]byte(installedChaincodeNamesList))
}

func (t *WASMChaincode) invoke(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) < 2 {
		return shim.Error("Incorrect number of arguments. Expecting chaincode name to invoke")
	}

	chaincodeName := args[0]
	funcToInvoke := args[1]

	//Initialize global variables for exported wasm functions
	r := Resolver{
		chaincodeName,stub,args[2:],nil,
	}

	// Get the state from the ledger
	ledgerChaincodeKey, _ := stub.CreateCompositeKey(chaincodeStoreIndex, []string{chaincodeName})
	Chaincodebytes, _ := stub.GetState(ledgerChaincodeKey)
	if Chaincodebytes == nil {
		jsonResp := "{\"Error\":\"No Chaincode for " + chaincodeName + "\"}"
		return shim.Error(jsonResp)
	}

	result := runWASM(Chaincodebytes, funcToInvoke,len(args)-2,&r)

	fmt.Printf("Invoke Response:%d\n", result)
	return txnResult(result,r.result)
}

// Store a new wasm chaincode in state. Receives chaincode name and wasm file encoded in hex
func (t *WASMChaincode) create(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("Create function")
	if len(args) < 2 {
		return shim.Error("Incorrect number of arguments. Expecting atleast 2 arguments")
	}

	chaincodeName := args[0]
	fmt.Println("Installing wasm chaincode: " + chaincodeName)

	//check if same chaincode name is already present
	chaincodeFromState, err := stub.GetState(chaincodeName)
	if chaincodeFromState != nil {
		return shim.Error(CHAINCODE_EXISTS)
	}

	//Decode the chaincode
	chaincodeHexEncoded := args[1]
	fmt.Println("Encoded wasm chaincode: " + chaincodeHexEncoded)
	chaincodeDecoded, err := hex.DecodeString(chaincodeHexEncoded)



	//Initialize global variables for exported wasm functions
	r := Resolver{
		chaincodeName,stub,args[2:],nil,
	}

	result := runWASM(chaincodeDecoded, "init", len(args)-2,&r)

	fmt.Printf("Init Response:%d\n", result)

	if result!=0 {
		return shim.Error("Chaincode init invocation failed")
	}

	// Store the chaincode in
	ledgerChaincodeKey, err := stub.CreateCompositeKey(chaincodeStoreIndex, []string{chaincodeName})
	err = stub.PutState(ledgerChaincodeKey, chaincodeDecoded)
	if err != nil {
		s := fmt.Sprintf(UKNOWN_ERROR, err.Error())
		return shim.Error(s)
	}
	return shim.Success([]byte("Success! Installed wasm chaincode"))
}

// query callback representing the query of a chaincode
func runWASM(Chaincodebytes []byte, funcToInvoke string, numberOfArgs int, r *Resolver) int64 {

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
	}, r, nil)

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
	result, err := vm.Run(entryID,int64(numberOfArgs))
	if err != nil {
		vm.PrintStackTrace()
		panic(err)
	}
	end := time.Now()

	fmt.Printf("return value = %d, duration = %v\n", result, end.Sub(start))

	return result
}

func txnResult(vmExecResult int64, resultGlobal []byte) pb.Response {

	//Shim.Error if result is negative
	if vmExecResult == -1 {

		//Check if some result is returned by wasm chaincode
		if resultGlobal ==nil{
			return shim.Error(strconv.FormatInt(vmExecResult, 10))
		}else {
			return shim.Error(string(resultGlobal))
		}
	}else{
		if resultGlobal ==nil{
			return shim.Success([]byte(strconv.FormatInt(vmExecResult, 10)))
		}else {
			return shim.Success(resultGlobal)
		}
	}
}

func main() {
	err := shim.Start(new(WASMChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
