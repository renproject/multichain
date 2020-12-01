extern crate libc;
use solana_sdk::{
    pubkey::Pubkey,
};
use std::{
    ffi::{CStr, CString},
    str::FromStr,
};

#[no_mangle]
pub extern "C" fn hello(name: *const libc::c_char) -> *const libc::c_char {
    let buf_name = unsafe { CStr::from_ptr(name).to_bytes() };
    let str_name = String::from_utf8(buf_name.to_vec()).unwrap();
    let greeting = format!("Hello {}!", str_name);
    CString::new(greeting).unwrap().into_raw()
}

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

// #[no_mangle]
// pub extern "C" fn ren_bridge_initialize(
//     keypair_path: *const libc::c_char,
//     rpc_url: *const libc::c_char,
//     authority_pointer: *const u8,
// ) -> *const libc::c_char {
//     // Solana default signer and fee payer.
//     let buf_name = unsafe { CStr::from_ptr(keypair_path).to_bytes() };
//     let keypair_path = String::from_utf8(buf_name.to_vec()).unwrap();
//     let payer = read_keypair_file(&keypair_path).unwrap();
//
//     // Initialize client.
//     let buf_name = unsafe { CStr::from_ptr(rpc_url).to_bytes() };
//     let rpc_url = String::from_utf8(buf_name.to_vec()).unwrap();
//     let rpc_client = RpcClient::new(rpc_url);
//     let commitment_config = CommitmentConfig::single_gossip();
//     let (recent_blockhash, _, _) = rpc_client
//         .get_recent_blockhash_with_commitment(commitment_config)
//         .unwrap()
//         .value;
//
//     // Construct the RenVM authority 20-bytes Ethereum compatible address.
//     let authority_slice =
//         unsafe { std::slice::from_raw_parts(authority_pointer as *const u8, 20usize) };
//     let mut authority = [0u8; 20usize];
//     authority.copy_from_slice(authority_slice);
//
//     // Find derived address that will hold RenBridge's state.
//     let (ren_bridge_account_id, _) =
//         Pubkey::find_program_address(&[b"RenBridgeState"], &ren_bridge::id());
//
//     // Build and sign the initialize transaction.
//     let mut tx = Transaction::new_with_payer(
//         &[initialize(
//             &ren_bridge::id(),
//             &payer.pubkey(),
//             &ren_bridge_account_id,
//             authority,
//         )
//         .unwrap()],
//         Some(&payer.pubkey()),
//     );
//     tx.sign(&[&payer], recent_blockhash);
//
//     // Broadcast transaction.
//     let signature = rpc_client
//         .send_transaction_with_config(
//             &tx,
//             RpcSendTransactionConfig {
//                 preflight_commitment: Some(commitment_config.commitment),
//                 ..RpcSendTransactionConfig::default()
//             },
//         )
//         .unwrap();
//
//     CString::new(signature.to_string()).unwrap().into_raw()
// }
//
// #[no_mangle]
// pub extern "C" fn ren_bridge_mint() {
//     unimplemented!();
// }
//
// #[no_mangle]
// pub extern "C" fn ren_bridge_burn() {
//     unimplemented!();
// }
