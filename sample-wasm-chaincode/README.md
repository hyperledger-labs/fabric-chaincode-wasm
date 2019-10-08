# Compile source code to WASM

## RUST

We have a sample wasm binary under: ```sample-wasm-chaincode/chaincode_example02/rust/app_main.wasm```

 - Install rust toolchain and wasm-pack. [Official Setup guide here](https://rustwasm.github.io/book/game-of-life/setup.html) , or use these two commands:
     ```
    curl https://sh.rustup.rs -sSf | sh
    curl https://rustwasm.github.io/wasm-pack/installer/init.sh -sSf | sh
     ```
 - Create a ```Cargo.toml``` ([sample](https://github.com/kleash/wasmer-chaincode-test/blob/master/sample-wasm-chaincode/chaincode_example02/rust/Cargo.toml)) in root directory with this content:
    ```
    [package]
    name = "app_main"
    version = "0.1.0"
    authors = ["shubham aggarwal <ag.shubham94@gmail.com>"]
    edition = "2018"
    
    [lib]
    crate-type = ["cdylib"]
    
    [dependencies]
    wasm-bindgen = "0.2"
    ```
 - Create a src folder and place ```lib.rs``` in source forlder. Your directory structure should look like this:
     ```
     .
    ├── Cargo.toml
    └── src
        └── lib.rs
    ```
 - From root directory, give command ```wasm-pack build```. It will take some time for first time to download all dependencies. Once done you will receive a similar message in console
 ```Your wasm pkg is ready to publish at ./pkg.```
 - If successful, the wasm binary can be located at ```pkg/app_main_bg.wasm```


## C
We have a sample wasm binary under: ```sample-wasm-chaincode/chaincode_example02/c/app_main.wasm```

 - Install clang8 and llvm [Linux repository link](https://apt.llvm.org)
 - To compile C file to wasm, issue following command
```
 clang \
   --target=wasm32 \
   -O3 \
   -flto \
   -nostdlib \
   -Wl,--no-entry \
   -Wl,--export-all \
   -Wl,--lto-O3 \
   -o app_main.wasm \
   main.c
```