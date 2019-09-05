use std::convert::TryFrom;
use std::str;



extern "C" {
    fn __print(msg: *const u8, len: usize) -> i64;
    fn __get_parameter(paramNumber: usize, result: *const u8) -> i64;
    fn __get_state(msg: *const u8, len: usize,value: *const u8) -> i64;
    fn __put_state(key: *const u8, keylen: usize, value: *const u8, valuelen: usize) -> i64;
    fn __delete_state(msg: *const u8, len: usize) -> i64;
    fn __return_result(msg: *const u8, len: usize) -> i64;
}


#[no_mangle]
pub extern "C" fn init(args: i64) -> i32 {

    if args != 4 {
        let s0 = "ERROR! Incorrect number of arguments. Expecting 4".as_bytes();
        unsafe {
            let _result = __print(s0.as_ptr(), s0.len());
        }
        return -1;
    }

    // Entities
    let a = [0; 24];
    let a_length;

    let b = [0; 24];
    let b_length;

    // Asset holdings
    let aval  = [0; 24];
    let aval_length;

    let bval  = [0; 24];
    let bval_length;



    unsafe {

        //parameter one
        let result_key_len = __get_parameter(0, a.as_ptr());
        a_length = usize::try_from(result_key_len).unwrap();

        //parameter two
        let result_key_len = __get_parameter(1, aval.as_ptr());
        aval_length = usize::try_from(result_key_len).unwrap();
        let the_bytes = &aval[0..aval_length];
        let the_string = str::from_utf8(the_bytes).expect("not UTF-8");
        let _the_number: u64 = the_string.parse().expect("not a number");

        //let s = String::from_utf8_lossy(&aval);
        //let parse_result = i32::from_str(&s).is_ok();

       /* if the_number ==100 {
            let s0 = "ERROR! Expecting integer value for asset holding \n".as_bytes();
            let _result = __print(s0.as_ptr(), s0.len());
            return -1;
        }
*/
        //parameter three
        let result_key_len = __get_parameter(2, b.as_ptr());
        b_length = usize::try_from(result_key_len).unwrap();


        //parameter four
        let result_key_len = __get_parameter(3, bval.as_ptr());
        bval_length = usize::try_from(result_key_len).unwrap();
        let the_bytes = &bval[0..bval_length];
        let the_string = str::from_utf8(the_bytes).expect("not UTF-8");
        let _the_number: u64 = the_string.parse().expect("not a number");

      /*  let s = String::from_utf8_lossy(&bval);
        let parse_result = i32::from_str(&s).is_ok();

        if parse_result {
            let s0 = "ERROR! Expecting integer value for asset holding".as_bytes();
            let _result = __print(s0.as_ptr(), s0.len());
            return -1;
        }
*/

        let result1 = __put_state(a.as_ptr(), a_length,aval.as_ptr(), aval_length);

        if result1 == -1 {
            return -1;
        }


        let result1 = __put_state(b.as_ptr(), b_length,bval.as_ptr(), bval_length);

        if result1 == -1 {
            return -1;
        }
    }

/*    let s = String::from_utf8_lossy(&aval);
    let parse_result = s.parse::<i32>().is_ok();

    unsafe {
        let s0 = s.as_bytes();
        let _result = __print(s0.as_ptr(), s0.len());
    }*/

    return 0;
}


#[no_mangle]
pub extern "C" fn invoke(args: i64) -> i32 {


    if args != 3 {
        let s0 = "ERROR! Incorrect number of arguments. Expecting 3".as_bytes();
        unsafe {
            let _result = __print(s0.as_ptr(), s0.len());
        }
        return -1;
    }

    // Entities
    let a = [0; 24];
    let a_length;

    let b = [0; 24];
    let b_length;

    //Transaction amount
    let txn_amount : u64;

    //Get parameters
    unsafe {

        //parameter one
        let result_key_len = __get_parameter(0, a.as_ptr());
        a_length = usize::try_from(result_key_len).unwrap();

        //parameter two
        let result_key_len = __get_parameter(1, b.as_ptr());
        b_length = usize::try_from(result_key_len).unwrap();

        //parameter three
        let txn_amount_ptr = [0; 24];
        let result_key_len = __get_parameter(2, txn_amount_ptr.as_ptr());
        let bval_length = usize::try_from(result_key_len).unwrap();
        let the_bytes = &txn_amount_ptr[0..bval_length];
        let the_string = str::from_utf8(the_bytes).expect("not UTF-8");
        txn_amount = the_string.parse().expect("not a number");
/*
        if !parse_result {
            let s0 = "ERROR! Invalid transaction amount, expecting a integer value".as_bytes();
            let _result = __print(s0.as_ptr(), s0.len());
            return -1;
        }
        txn_amount = s.parse::<i32>().unwrap();*/
    }

    //Get asset balance of A and B
    let a_val : u64;
    let b_val : u64;


    unsafe {

        //asset one
        let get_result = [0; 24];
        let result_key_len = __get_state(a.as_ptr(), a_length,get_result.as_ptr());
        if result_key_len==-1{
            let s0 = "ERROR! Entity not found".as_bytes();
            let _result = __print(s0.as_ptr(), s0.len());
            return -1;
        }
        // convert byte array -> string -> integer
        let bval_length = usize::try_from(result_key_len).unwrap();
        let the_bytes = &get_result[0..bval_length];
        let the_string = str::from_utf8(the_bytes).expect("not UTF-8");
        a_val = the_string.parse().expect("not a number");


        //asset two
        let get_result = [0; 24];
        let result_key_len = __get_state(b.as_ptr(), b_length,get_result.as_ptr());
        if result_key_len==-1{
            let s0 = "ERROR! Entity not found".as_bytes();
            let _result = __print(s0.as_ptr(), s0.len());
            return -1;
        }

        let bval_length = usize::try_from(result_key_len).unwrap();
        let the_bytes = &get_result[0..bval_length];
        let the_string = str::from_utf8(the_bytes).expect("not UTF-8");
        b_val = the_string.parse().expect("not a number");
    }


    // Perform the execution
    let a_val = a_val - txn_amount;
    let b_val = b_val + txn_amount;

    let ttemp = format!("{}{}{}{}", "Aval = ", a_val," Bval = ",b_val);
    let t = ttemp.as_bytes();

    // Write the state back to the ledger
    unsafe {
        __print(t.as_ptr(), t.len());

        let stemp =a_val.to_string();
        let s = stemp.as_bytes();
        let result1 = __put_state(a.as_ptr(), a_length,s.as_ptr(), s.len());

        if result1==-1 {
            return -1;
        }

        let stemp =b_val.to_string();
            let s = stemp.as_bytes();
        let result1 = __put_state(b.as_ptr(), b_length,s.as_ptr(), s.len());

        if result1==-1 {
            return -1;
        }
    }

    return 0;
}


#[no_mangle]
pub extern "C" fn query(args: i64) -> i32 {

    if args != 1 {
        let s0 = "ERROR! Incorrect number of arguments. Expecting name of the person to query".as_bytes();
        unsafe {
            let _result = __print(s0.as_ptr(), s0.len());
        }
        return -1;
    }

    //Declare an array with expected value size
    let get_result = [0; 24];

    let a = [0; 24];
    let a_length;

    unsafe {
       //parameter one
       let result_key_len = __get_parameter(0, a.as_ptr());
       a_length = usize::try_from(result_key_len).unwrap();

       //get state
       let len_get_state = __get_state(a.as_ptr(), a_length,get_result.as_ptr());

        let _b = usize::try_from(len_get_state);
        let _result = __print(get_result.as_ptr(), _b.unwrap());


       __return_result(get_result.as_ptr(), get_result.len());
    }
    return 0;
}



#[no_mangle]
pub extern "C" fn delete(args: i64) -> i32 {

    if args != 1 {
        let s0 = "ERROR! Incorrect number of arguments. Expecting 1".as_bytes();
        unsafe {
            let _result = __print(s0.as_ptr(), s0.len());
        }
        return -1;
    }

    //Declare an array with expected value size
    let get_result = [0; 24];

    let a = [0; 24];
    let a_length;

    unsafe {
        //parameter one
        let result_key_len = __get_parameter(0, a.as_ptr());
        a_length = usize::try_from(result_key_len).unwrap();

        //get state
        let result = __delete_state(a.as_ptr(), a_length);

        if result==-1{
            let s0 = "Failed to delete state".as_bytes();
            __return_result(get_result.as_ptr(), get_result.len());
            return -1;
        }
    }
    return 0;
}