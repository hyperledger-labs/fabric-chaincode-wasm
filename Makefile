.PHONY: all build unit-test clean

all: build unit-test clean

build: bin/wasmcc bin/wasm-pusher

clean:
	rm -rf bin/

unit-test: sample-wasm-chaincode/chaincode_example02/rust/app_main.zip
	cd wasmcc && go test -v 

.PHONY: bin/wasmcc bin/wasm-pusher
bin/wasmcc:
	mkdir -p bin/
	cd wasmcc/ && go build -o ../bin/wasmcc
	rm bin/wasmcc # check that it builds, chaincode not meant to be used directly

bin/wasm-pusher:
	mkdir -p bin/
	cd tools/wasm-pusher && go build -o ../../bin/wasm-pusher

sample-wasm-chaincode/chaincode_example02/rust/app_main.zip: sample-wasm-chaincode/chaincode_example02/rust/app_main.wasm
	zip sample-wasm-chaincode/chaincode_example02/rust/app_main.zip sample-wasm-chaincode/chaincode_example02/rust/app_main.wasm
