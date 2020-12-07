extern crate libc;
use solana_sdk::{
    pubkey::Pubkey,
};
use std::{
    ffi::{CStr, CString},
    str::FromStr,
};

#[no_mangle]
pub extern "C" fn unique_pubkey() -> *const libc::c_char {
    let pubkey = Pubkey::new_unique();
    let pubkey = pubkey.to_string();
    CString::new(pubkey).unwrap().into_raw()
}

#[no_mangle]
pub extern "C" fn program_derived_address(
    seeds_pointer: *const u8,
    seeds_size: libc::size_t,
    program: *const libc::c_char,
) -> *const libc::c_char {
    let seeds =
        unsafe { std::slice::from_raw_parts(seeds_pointer as *const u8, seeds_size as usize) };

    let buf_name = unsafe { CStr::from_ptr(program).to_bytes() };
    let program_str = String::from_utf8(buf_name.to_vec()).unwrap();
    let program_id = Pubkey::from_str(&program_str).unwrap();

    let (derived_address, _) = Pubkey::find_program_address(&[seeds], &program_id);

    CString::new(derived_address.to_string())
        .unwrap()
        .into_raw()
}
