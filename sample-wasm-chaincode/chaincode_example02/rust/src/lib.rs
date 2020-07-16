use std::convert::TryFrom;
use std::str;

extern "C" {
    fn __print(msg: *const u8, len: usize) -> i64;
    fn __get_parameter(paramNumber: usize, result: *const u8) -> i64;
    fn __get_parameter_size(paramNumber: usize) -> i64;
    fn __get_state(msg: *const u8, len: usize, value: *const u8) -> i64;
    fn __get_state_size(msg: *const u8, len: usize) -> i64;
    fn __put_state(key: *const u8, key_len: usize, value: *const u8, value_len: usize) -> i64;
    fn __delete_state(msg: *const u8, len: usize) -> i64;
    fn __return_result(msg: *const u8, len: usize) -> i64;
}

/// Calls host function to return invocation result.
/// This result will be returned as transaction response to user.
fn return_result(msg: *const u8, len: usize) -> i64 {
    return unsafe { __return_result(msg, len) };
}

/// Calls host function to retrieve transaction parameter
fn get_parameter(param_number: usize, result: *const u8) -> i64 {
    let result_key_len = unsafe { __get_parameter(param_number, result) };

    if result_key_len < 0 {
        let error_msg = ("Unable to retrieve transaction parameter").as_bytes();
        print(error_msg.as_ptr(), error_msg.len());
        return_result(error_msg.as_ptr(), error_msg.len());
        return -1;
    }

    return result_key_len;
}

/// Calls host function to log a message. It will help in debugging.
fn print(msg: *const u8, len: usize) -> i64 {
    return unsafe { __print(msg, len) };
}


/// Init function accepts 4 transaction parameters i.e. two account names and corresponding balances.
/// It tries to store these both accounts in ledger.
#[no_mangle]
pub extern "C" fn init(args: i64) -> i64 {
    if args != 4 {
        let error_msg = "ERROR! Incorrect number of arguments. Expecting 4".as_bytes();
        print(error_msg.as_ptr(), error_msg.len());
        return_result(error_msg.as_ptr(), error_msg.len());
        return -1;
    }


    let first_param_len = unsafe { __get_parameter_size(0) }; //Account 1 name
    let second_param_len = unsafe { __get_parameter_size(1) }; //Account 1 balance
    let third_param_len = unsafe { __get_parameter_size(2) }; // Account 2 name
    let fourth_param_len = unsafe { __get_parameter_size(3) }; // Account 2 balance

    let first_account_name_length:usize = usize::try_from(first_param_len).unwrap();
    let first_account = vec![0; first_account_name_length];

    let second_account_name_length:usize = usize::try_from(third_param_len).unwrap();
    let second_account = vec![0; second_account_name_length];

    // Account's asset holdings
    let first_account_balance_length:usize = usize::try_from(second_param_len).unwrap();
    let first_account_balance = vec![0; first_account_balance_length];

    let second_account_balance_length:usize = usize::try_from(fourth_param_len).unwrap();
    let second_account_balance = vec![0; second_account_balance_length];


    //transaction parameter one as first account name
    get_parameter(0, first_account.as_ptr());

    //transaction parameter two as asset balance of first account
    get_parameter(1, first_account_balance.as_ptr());

    //Validate and convert balance to integer
    let the_bytes = &first_account_balance[0..first_account_balance_length];
    let the_string = str::from_utf8(the_bytes).expect("not UTF-8");
    let _the_number: u64 = the_string.parse().expect("not a number");

    //transaction parameter three as second account name
    get_parameter(2, second_account.as_ptr());

    //transaction parameter four as asset balance of second account
    get_parameter(3, second_account_balance.as_ptr());

    //Validate and convert balance to integer
    let the_bytes = &second_account_balance[0..second_account_balance_length];
    let the_string = str::from_utf8(the_bytes).expect("not UTF-8");
    let _the_number: u64 = the_string.parse().expect("not a number");


    let put_state_result = unsafe { __put_state(first_account.as_ptr(), first_account_name_length, first_account_balance.as_ptr(), first_account_balance_length) };
    if put_state_result == -1 {
        let error_msg = "ERROR! Unable to insert first account to state".as_bytes();
        return_result(error_msg.as_ptr(), error_msg.len());
        return -1;
    }


    let put_state_result = unsafe { __put_state(second_account.as_ptr(), second_account_name_length, second_account_balance.as_ptr(), second_account_balance_length) };
    if put_state_result == -1 {
        let error_msg = "ERROR! Unable to insert second account to state".as_bytes();
        return_result(error_msg.as_ptr(), error_msg.len());
        return -1;
    }

    return 0;
}

/// Invoke function accepts three transaction parameters i.e. from account, to account, units to be transferred.
/// It retrieves the balance of both accounts from state, updates the balance and store it in state.
#[no_mangle]
pub extern "C" fn invoke(args: i64) -> i64 {
    if args != 3 {
        let error_msg = "ERROR! Incorrect number of arguments. Expecting 3".as_bytes();
        print(error_msg.as_ptr(), error_msg.len());
        return -1;
    }

    // Entities
    let first_param_len = unsafe { __get_parameter_size(0) }; // Account 1 name
    let first_account_name_len:usize = usize::try_from(first_param_len).unwrap();
    let first_account = vec![0; first_account_name_len];

    let second_param_len = unsafe { __get_parameter_size(1) }; // Account 1 name
    let second_account_len:usize = usize::try_from(second_param_len).unwrap();
    let second_account = vec![0; second_account_len];

    //Transaction amount
    let txn_amount: u64;

    //get from account
    get_parameter(0, first_account.as_ptr());

    //get to account
    get_parameter(1, second_account.as_ptr());

    //get amount to transfer
    let third_param_len = unsafe { __get_parameter_size(2) }; // amount to transfer
    let txn_amount_length:usize = usize::try_from(third_param_len).unwrap();
    let txn_amount_ptr = vec![0; txn_amount_length];
    get_parameter(2, txn_amount_ptr.as_ptr());
    //Validate and convert amount to integer
    let the_bytes = &txn_amount_ptr[0..txn_amount_length];
    let the_string = str::from_utf8(the_bytes).expect("not UTF-8");
    txn_amount = the_string.parse().expect("not a number");


    //Get asset balance of A and B
    let from_account_balance: u64;
    let to_account_balance: u64;

    //get from account balance
    let account_balance_state_len = unsafe { __get_state_size(first_account.as_ptr(), first_account_name_len) }; // amount from transfer
    if account_balance_state_len == -1 {
        let error_msg = "ERROR! from account not found".as_bytes();
        print(error_msg.as_ptr(), error_msg.len());
        return_result(error_msg.as_ptr(), error_msg.len());
        return -1;
    }
    let account_bal_length:usize = usize::try_from(account_balance_state_len).unwrap();
    let get_result = vec![0; account_bal_length];
    unsafe { __get_state(first_account.as_ptr(), first_account_name_len, get_result.as_ptr()) };

    // convert byte array -> string -> integer
    let the_bytes = &get_result[0..account_bal_length];
    let the_string = str::from_utf8(the_bytes).expect("not UTF-8");
    from_account_balance = the_string.parse().expect("not a number");

    //get to account balance
    let account_balance_state_len = unsafe { __get_state_size(second_account.as_ptr(), second_account_len) }; // amount to transfer
    if account_balance_state_len == -1 {
        let error_msg = "ERROR! to account not found".as_bytes();
        print(error_msg.as_ptr(), error_msg.len());
        return_result(error_msg.as_ptr(), error_msg.len());
        return -1;
    }
    let account_bal_length:usize = usize::try_from(account_balance_state_len).unwrap();
    let get_result = vec![0; account_bal_length];
    unsafe { __get_state(second_account.as_ptr(), second_account_len, get_result.as_ptr()) };

    let the_bytes = &get_result[0..account_bal_length];
    let the_string = str::from_utf8(the_bytes).expect("not UTF-8");
    to_account_balance = the_string.parse().expect("not a number");


    //validate from account balance
    if from_account_balance < txn_amount {
        let error_msg = "ERROR! insufficient units in from account".as_bytes();
        print(error_msg.as_ptr(), error_msg.len());
        return_result(error_msg.as_ptr(), error_msg.len());
        return -1;
    }

    // Perform the execution
    let from_account_balance_updated = from_account_balance - txn_amount;
    let to_account_balance_updated = to_account_balance + txn_amount;

    let updated_bal_msg = format!("{}{}{}{}", "Aval = ", from_account_balance_updated, " Bval = ", to_account_balance_updated);
    let bal_msg = updated_bal_msg.as_bytes();
    print(bal_msg.as_ptr(), bal_msg.len());

    // update from account balance to the ledger
    let string_from_account_balance_updated = from_account_balance_updated.to_string();
    let bytes_from_account_balance_updated = string_from_account_balance_updated.as_bytes();
    let put_state_result = unsafe { __put_state(first_account.as_ptr(), first_account_name_len, bytes_from_account_balance_updated.as_ptr(), bytes_from_account_balance_updated.len()) };

    if put_state_result == -1 {
        let error_msg = "Unable to update from account balance".as_bytes();
        print(error_msg.as_ptr(), error_msg.len());
        return_result(error_msg.as_ptr(), error_msg.len());
        return -1;
    }

    // update to account balance to the ledger
    let string_to_account_balance_updated = to_account_balance_updated.to_string();
    let bytes_to_account_balance_updated = string_to_account_balance_updated.as_bytes();
    let put_state_result = unsafe { __put_state(second_account.as_ptr(), second_account_len, bytes_to_account_balance_updated.as_ptr(), bytes_to_account_balance_updated.len()) };

    if put_state_result == -1 {
        let error_msg = "Unable to update to account balance".as_bytes();
        print(error_msg.as_ptr(), error_msg.len());
        return_result(error_msg.as_ptr(), error_msg.len());
        return -1;
    }


    let success_msg = "Successfully transferred account balance".as_bytes();
    return_result(success_msg.as_ptr(), success_msg.len());
    return 0;
}

/// Query function accepts one transaction parameter i.e. account name.
/// It retrieves account balance from state and returns that as function response.
#[no_mangle]
pub extern "C" fn query(args: i64) -> i64 {
    if args != 1 {
        let error_msg = "ERROR! Incorrect number of arguments. Expecting name of the person to query".as_bytes();
        print(error_msg.as_ptr(), error_msg.len());
        return_result(error_msg.as_ptr(), error_msg.len());
        return -1;
    }

    //Declare an array with expected value size
    let account_balance = [0; 24];

    let account_name = [0; 24];
    let account_name_len;

    //parameter one
    let result_key_len = get_parameter(0, account_name.as_ptr());
    account_name_len = usize::try_from(result_key_len).unwrap();

    //get account balance from state
    let len_get_state = unsafe { __get_state(account_name.as_ptr(), account_name_len, account_balance.as_ptr()) };
    if len_get_state == -1 {
        let error_msg = "ERROR! account not found".as_bytes();
        print(error_msg.as_ptr(), error_msg.len());
        return_result(error_msg.as_ptr(), error_msg.len());
        return -1;
    }
    let account_balance_state_len = usize::try_from(len_get_state);
    return_result(account_balance.as_ptr(), account_balance_state_len.unwrap());
    return 0;
}

/// Delete function accepts one transaction parameter i.e. account name.
/// It tries to delete the account from state.
#[no_mangle]
pub extern "C" fn delete(args: i64) -> i64 {
    if args != 1 {
        let error_msg = "ERROR! Incorrect number of arguments. Expecting 1".as_bytes();
        print(error_msg.as_ptr(), error_msg.len());
        return_result(error_msg.as_ptr(), error_msg.len());
        return -1;
    }

    let account_name = [0; 24];
    let account_name_length;

    //parameter one
    let result_key_len = get_parameter(0, account_name.as_ptr());
    account_name_length = usize::try_from(result_key_len).unwrap();

    //delete state
    let delete_state_result = unsafe { __delete_state(account_name.as_ptr(), account_name_length) };

    if delete_state_result == -1 {
        let error_msg = "Failed to delete state".as_bytes();
        print(error_msg.as_ptr(), error_msg.len());
        return_result(error_msg.as_ptr(), error_msg.len());
        return -1;
    }

    let success_msg = "Success! Account deleted".as_bytes();
    return_result(success_msg.as_ptr(), success_msg.len());
    return 0;
}