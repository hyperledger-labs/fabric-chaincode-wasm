# WASM Pusher

This tool is based on fabric-go-sdk. It can be used to directly install a wasm chaincode.

WASMCC supports wasm chaincode in three formats i.e. .wasm binary, zip file containing wasm binary and hex encoded wasm file. WASM pusher can install wasm chaincode in any of the above formats.

```
go build
./wasm-pusher -n balancewasm -w ../../sample-wasm-chaincode/chaincode_example02/rust/app_main.wasm -u User1 -a a,100,b,100
```