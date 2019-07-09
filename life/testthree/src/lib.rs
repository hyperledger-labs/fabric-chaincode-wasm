extern "C" {
    fn __get_state(msg: *const u8, len: usize);
    fn __put_state(key: *const u8, keylen: usize, value: *const u8, valuelen: usize);
}

#[no_mangle]
pub extern "C" fn init() -> i32 {
    let key = "account1".as_bytes();
    let value = "100".as_bytes();
    let key2 = "account2".as_bytes();
    let value2 = "10".as_bytes();

    unsafe {
        __put_state(key.as_ptr(), key.len(),value.as_ptr(), value.len());
    }


    unsafe {
        __put_state(key2.as_ptr(), key2.len(),value2.as_ptr(), value2.len());
    }

    //TODO: Return a success or fail message from here
    return 0;
}

#[no_mangle]
pub extern "C" fn get_balance() -> i32 {
    let message = "account1".as_bytes();

    unsafe {
        __get_state(message.as_ptr(), message.len());
    }

    //TODO: Return response of get state from here
    return 0;
}