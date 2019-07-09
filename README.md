# WASM Chaincode

VM runtime selected based on evaluation of development activity and resources available online:
 - [wasmer](https://github.com/wasmerio/go-ext-wasm)
 - [life](https://github.com/perlin-network/life)


## Strategy

 - **How to store wasm chaincodes**
 Store as byte array in state against chaincode name
 - **How wasm chaincodes will interact with peer without a shim layer?**
 WASM supports importing functions from host modules. So our wasmcc chaincode will expose shim layer functions as wasm exported functions which can then in turn will be imported by chaincode developers in their wasm chaincode
 - **How to instantiate or invoke wasm chaincode functions?**
 Host modules can call wasm functions directly. We will define some standard definations for init and invoke functions, which need to be implemented by wasm chaincode developers and exported from their modules.

### Strategy of storing WASM Chaincodes
![Wasm storage](https://github.com/kleash/wasmer-chaincode-test/blob/master/docs/gliffy/wasmcc-create.png)


### Strategy of WASM invocation
![Wasm storage](https://github.com/kleash/wasmer-chaincode-test/blob/master/docs/gliffy/wasmcc-invoke.png)

## Functionality tested

 - WASM invocation using wasmer : [Link](https://github.com/kleash/wasmer-chaincode-test/tree/master/wasmer/1)
 - WASM import/export functions using wasmer : [Link](https://github.com/kleash/wasmer-chaincode-test/tree/master/wasmer/2)
 - WASM import/export functions using life : [Link](https://github.com/kleash/wasmer-chaincode-test/tree/master/life/testone)

## Current Progress

Tried a dummy smart contract to integrate both of the above WASM VMs
 - Contains ```create``` and ```query``` functions
 - ```create``` will accept wasm chaincode name and hex encoded wasm chaincode
 - then it stores the chaincode in state
 - ```query``` acceptes wasm chaincode name and function to invoke
 - ```query``` retrieves the wasm chaincode bytes from state and execute it in wasm vm
 - ```query``` dynamically invokes the function from wasm chaincode whose name it accepted as a parameter initialy


### WASMER Issues

WASMER Chaincode: [link](https://github.com/kleash/wasmer-chaincode-test/tree/master/smartcontract/wasmercontractone)

WASMER need some C libraries and go modules for successful build. It's throwing following error at chaincode instantiation:
```
2019-06-20 01:45:29.309 UTC [chaincodeCmd] checkChaincodeCmdParams -> INFO 001 Using default escc
2019-06-20 01:45:29.309 UTC [chaincodeCmd] checkChaincodeCmdParams -> INFO 002 Using default vscc
Error: could not assemble transaction, err proposal response was not successful, error code 500, msg error starting container: error starting container: Failed to generate platform-specific docker build: Error returned from build: 2 "# github.com/chaincode/awesomeProject/vendor/github.com/wasmerio/go-ext-wasm/wasmer
/usr/bin/ld: cannot find -lwasmer_runtime_c_api
collect2: error: ld returned 1 exit status
```

**Probable Cause:** wasmer_runtime_c_api.so is not being packaged as part of install chaincode.

**Next Step:** Somehow include the c library

## life Issues


life Chaincode: [link](https://github.com/kleash/wasmer-chaincode-test/tree/master/smartcontract/lifecontractthree)

**Successfully installed wasmcc and invoked a wasm chaincode**

**1. Tool to convert a file to string**:
Navigate to ``tools/file-encoder``

Invoke: ``./encoder``

Pass the absolute or relative path of your webassembly module/chaincode. For example, you can pass sample wasm in this repository: ``life/testthree/app_main.wasm``

The tool will copy the encoded chaincode to your clipboard and also display it over console.


**2. Install wasmcc chaincode:**
```
peer chaincode install -n wasmcc -v 1.0 -p github.com/chaincode/lifecontractthree
```

and
```
CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/users/Admin@org2.example.com/msp CORE_PEER_ADDRESS=peer0.org2.example.com:9051 CORE_PEER_LOCALMSPID="Org2MSP" CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt peer chaincode install -n wasmcc -v 1.0 -p github.com/chaincode/lifecontractthree
```

**3. Instantiate wasmcc chaincode:**
```
peer chaincode instantiate -o orderer.example.com:7050 --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C mychannel -n wasmcc -v 1.0 -c '{"Args":["init"]}' -P "AND ('Org1MSP.peer','Org2MSP.peer')"
```

**4. Install a wasm chaincode**

 - 1st argument is create function i.e. to create a new webassembly chaincode.
 - 2nd argument is chaincode name.
 - 3rd argument is encoded file string which we got in 1st step.


```
peer chaincode invoke -o orderer.example.com:7050 --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C mychannel -n wasmcc --peerAddresses peer0.org1.example.com:7051 --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt --peerAddresses peer0.org2.example.com:9051 --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt -c '{"Args":["create","balancewasm","0061736d010000000191808080000360047f7f7f7f0060027f7f006000017f02a5808080000203656e760b5f5f7075745f7374617465000003656e760b5f5f6765745f73746174650001038380808000020202058380808000010011079f8080800003066d656d6f7279020004696e697400020b6765745f62616c616e636500030ab480808000022400418080c0004108418880c00041031000418b80c0004108419380c0004102100041000b0d00418080c0004108100141000b0b9e808080000100418080c0000b156163636f756e74313130306163636f756e7432313000b980808000046e616d6501ae8080800004000b5f5f7075745f7374617465010b5f5f6765745f73746174650204696e6974030b6765745f62616c616e636500fe808080000970726f64756365727302086c616e677561676501045275737404323031380c70726f6365737365642d6279030572757374631d312e33352e30202833633233356435363020323031392d30352d3230290677616c72757305302e382e300c7761736d2d62696e6467656e12302e322e3437202861316663323730663229"]}'
```
**5. Invoke wasm chaincode:**

 - 2nd argument is chaincode name.
 - 3rd argument is function name of the webassembly module we want to invoke.


```
peer chaincode query -C mychannel -n wasmcc -c '{"Args":["query","balancewasm","get_balance"]}'
```

**Next Step:** 
Fix all TODOs in go and rust files, namely:
 - RUST/WASM: ``TODO:`` Return response of get_balance function. (How to return a string from wasm function!)
 - GO : ``TODO:`` Send parameters to any wasm function. As of now, params is accepted as ...int64
 - GO: ``TODO:`` Return the result of get state to wasm from get state wrapper.

All the above TODOs are present in code also for better understanding

~~Create getState and putState wrapper so it can be executed by a wasm module~~


~~**Successfully Instantiated**~~

~~Install a wasm chaincode and try to execute it~~

~~life Chaincode: [link](https://github.com/kleash/wasmer-chaincode-test/tree/master/smartcontract/lifecontractone)~~


~~**It's throwing following error at chaincode instantiation:**~~
```
connecto/input/src/github.com/connecto/lifecontractone/vendor/github.com/perlin-network/life/compiler/module.go:147:31: too many arguments in call to disasm.Disassemble
	have (wasm.Function, *wasm.Module)
	want ([]byte)
connecto/input/src/github.com/connecto/lifecontractone/vendor/github.com/perlin-network/life/compiler/module.go:151:39: cannot use d (type []disasm.Instr) as type *disasm.Disassembly in argument to NewSSAFunctionCompiler
connecto/input/src/github.com/connecto/lifecontractone/vendor/github.com/perlin-network/life/compiler/module.go:223:31: too many arguments in call to disasm.Disassemble
	have (wasm.Function, *wasm.Module)
	want ([]byte)
connecto/input/src/github.com/connecto/lifecontractone/vendor/github.com/perlin-network/life/compiler/module.go:227:39: cannot use d (type []disasm.Instr) as type *disasm.Disassembly in argument to NewSSAFunctionCompiler
```

~~**Next Step:** Issue opened at github of life: https://github.com/perlin-network/life/issues/82~~