extern "C" {
    fn __get_state(msg: *const u8, len: usize);
}

#[no_mangle]
pub extern "C" fn app_main() -> i32 {
    let message = "key1!".as_bytes();

    unsafe {
        __life_log(message.as_ptr(), message.len());
    }

    return 0;
}