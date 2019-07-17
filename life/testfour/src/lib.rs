use std::convert::TryFrom;

extern "C" {
    fn __print(msg: *const u8, len: usize) -> i64;
    fn __get_state(msg: *const u8, len: usize,value: *const u8) -> i64;
    fn __put_state(key: *const u8, keylen: usize, value: *const u8, valuelen: usize) -> i64;
}

#[no_mangle]
pub extern "C" fn init() -> i32 {
    let key1 = "account1".as_bytes();
    let value1 = "INR 100".as_bytes();
    let key2 = "account2".as_bytes();
    let value2 = "10".as_bytes();

    unsafe {
        let result1 = __put_state(key1.as_ptr(), key1.len(),value1.as_ptr(), value1.len());

        if result1==-1 {
            return -1;
        }

        let result2 = __put_state(key2.as_ptr(), key2.len(),value2.as_ptr(), value2.len());


        if result2==-1 {
            return -1;
        }
    }

    return 0;
}


#[no_mangle]
pub extern "C" fn get_balance() -> i32 {
    let message = "account1".as_bytes();

    //Declare an array with expected value size
    let get_result = [0; 24];


    unsafe {
       let len_get_state = __get_state(message.as_ptr(), message.len(),get_result.as_ptr());
  //     let n_us : usize::try_from(len_get_state) =unwrap();


        let _b = len_get_state as usize;
        let _b = usize::try_from(len_get_state);
        let _result = __print(message.as_ptr(), message.len());
        let _result = __print(get_result.as_ptr(), _b.unwrap());

    }

    //TODO: Return response of get state from here
    return 0;
}