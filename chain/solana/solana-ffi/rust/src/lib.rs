extern crate libc;
use arrayref::array_refs;
use digest::Digest;
use ren_bridge::{
    instruction::{burn, initialize, initialize_token, mint},
    state::RenBridge,
};
use renvm_sig::{RenVM, RenVmMsgBuilder};
use solana_client::{rpc_client::RpcClient, rpc_config::RpcSendTransactionConfig};
use solana_sdk::{
    commitment_config::CommitmentConfig,
    program_pack::Pack,
    pubkey::Pubkey,
    signature::{read_keypair_file, Signer},
    transaction::Transaction,
};
use spl_associated_token_account::{create_associated_token_account, get_associated_token_address};
use spl_token::instruction::burn_checked;
use std::{
    ffi::{CStr, CString},
    str::FromStr,
};

mod util;

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

#[no_mangle]
pub extern "C" fn ren_bridge_initialize(
    keypair_path: *const libc::c_char,
    rpc_url: *const libc::c_char,
    authority_pointer: *const u8,
) -> *const libc::c_char {
    // Solana default signer and fee payer.
    let buf_name = unsafe { CStr::from_ptr(keypair_path).to_bytes() };
    let keypair_path = String::from_utf8(buf_name.to_vec()).unwrap();
    let payer = read_keypair_file(&keypair_path).unwrap();

    // Initialize client.
    let buf_name = unsafe { CStr::from_ptr(rpc_url).to_bytes() };
    let rpc_url = String::from_utf8(buf_name.to_vec()).unwrap();
    let rpc_client = RpcClient::new(rpc_url);
    let commitment_config = CommitmentConfig::single_gossip();
    let (recent_blockhash, _, _) = rpc_client
        .get_recent_blockhash_with_commitment(commitment_config)
        .unwrap()
        .value;

    // Construct the RenVM authority 20-bytes Ethereum compatible address.
    let authority_slice =
        unsafe { std::slice::from_raw_parts(authority_pointer as *const u8, 20usize) };
    let mut authority = [0u8; 20usize];
    authority.copy_from_slice(authority_slice);

    // Find derived address that will hold RenBridge's state.
    let (ren_bridge_account_id, _) =
        Pubkey::find_program_address(&[b"RenBridgeState"], &ren_bridge::id());

    // Build and sign the initialize transaction.
    let mut tx = Transaction::new_with_payer(
        &[initialize(
            &ren_bridge::id(),
            &payer.pubkey(),
            &ren_bridge_account_id,
            authority,
        )
        .unwrap()],
        Some(&payer.pubkey()),
    );
    tx.sign(&[&payer], recent_blockhash);

    // Broadcast transaction.
    let signature = rpc_client
        .send_transaction_with_config(
            &tx,
            RpcSendTransactionConfig {
                preflight_commitment: Some(commitment_config.commitment),
                ..RpcSendTransactionConfig::default()
            },
        )
        .unwrap();

    CString::new(signature.to_string()).unwrap().into_raw()
}

#[no_mangle]
pub extern "C" fn ren_bridge_initialize_token(
    keypair_path: *const libc::c_char,
    rpc_url: *const libc::c_char,
    selector: *const libc::c_char,
) -> *const libc::c_char {
    // Solana default signer and fee payer.
    let buf_name = unsafe { CStr::from_ptr(keypair_path).to_bytes() };
    let keypair_path = String::from_utf8(buf_name.to_vec()).unwrap();
    let payer = read_keypair_file(&keypair_path).unwrap();

    // Initialize client.
    let buf_name = unsafe { CStr::from_ptr(rpc_url).to_bytes() };
    let rpc_url = String::from_utf8(buf_name.to_vec()).unwrap();
    let rpc_client = RpcClient::new(rpc_url);
    let commitment_config = CommitmentConfig::single_gossip();
    let (recent_blockhash, _, _) = rpc_client
        .get_recent_blockhash_with_commitment(commitment_config)
        .unwrap()
        .value;

    // Get selector hash.
    let buf_name = unsafe { CStr::from_ptr(selector).to_bytes() };
    let selector = String::from_utf8(buf_name.to_vec()).unwrap();
    let mut hasher = sha3::Keccak256::new();
    hasher.update(selector.as_bytes());
    let selector_hash: [u8; 32] = hasher.finalize().into();

    // Derived address that will be the token mint.
    let (token_mint_id, _) = Pubkey::find_program_address(&[&selector_hash[..]], &ren_bridge::id());

    // Build and sign the initialize transaction.
    let mut tx = Transaction::new_with_payer(
        &[initialize_token(
            &ren_bridge::id(),
            &payer.pubkey(),
            &token_mint_id,
            &spl_token::id(),
            selector_hash,
        )
        .unwrap()],
        Some(&payer.pubkey()),
    );
    tx.sign(&[&payer], recent_blockhash);

    // Broadcast transaction.
    let signature = rpc_client
        .send_transaction_with_config(
            &tx,
            RpcSendTransactionConfig {
                preflight_commitment: Some(commitment_config.commitment),
                ..RpcSendTransactionConfig::default()
            },
        )
        .unwrap();

    CString::new(signature.to_string()).unwrap().into_raw()
}

#[no_mangle]
pub extern "C" fn ren_bridge_initialize_account(
    keypair_path: *const libc::c_char,
    rpc_url: *const libc::c_char,
    selector: *const libc::c_char,
) -> *const libc::c_char {
    // Solana default signer and fee payer.
    let buf_name = unsafe { CStr::from_ptr(keypair_path).to_bytes() };
    let keypair_path = String::from_utf8(buf_name.to_vec()).unwrap();
    let payer = read_keypair_file(&keypair_path).unwrap();

    // Initialize client.
    let buf_name = unsafe { CStr::from_ptr(rpc_url).to_bytes() };
    let rpc_url = String::from_utf8(buf_name.to_vec()).unwrap();
    let rpc_client = RpcClient::new(rpc_url);
    let commitment_config = CommitmentConfig::single_gossip();
    let (recent_blockhash, _, _) = rpc_client
        .get_recent_blockhash_with_commitment(commitment_config)
        .unwrap()
        .value;

    // Get selector hash.
    let buf_name = unsafe { CStr::from_ptr(selector).to_bytes() };
    let selector = String::from_utf8(buf_name.to_vec()).unwrap();
    let mut hasher = sha3::Keccak256::new();
    hasher.update(selector.as_bytes());
    let selector_hash: [u8; 32] = hasher.finalize().into();

    // Derived address that will be the token mint.
    let (token_mint_id, _) = Pubkey::find_program_address(&[&selector_hash[..]], &ren_bridge::id());

    // Build and sign transaction.
    let mut tx = Transaction::new_with_payer(
        &[create_associated_token_account(
            &payer.pubkey(),
            &payer.pubkey(),
            &token_mint_id,
        )],
        Some(&payer.pubkey()),
    );
    tx.sign(&[&payer], recent_blockhash);

    // Broadcast transaction.
    let signature = rpc_client
        .send_transaction_with_config(
            &tx,
            RpcSendTransactionConfig {
                preflight_commitment: Some(commitment_config.commitment),
                ..RpcSendTransactionConfig::default()
            },
        )
        .unwrap();

    CString::new(signature.to_string()).unwrap().into_raw()
}

#[no_mangle]
pub extern "C" fn ren_bridge_get_burn_count(rpc_url: *const libc::c_char) -> u64 {
    // Initialize client.
    let buf_name = unsafe { CStr::from_ptr(rpc_url).to_bytes() };
    let rpc_url = String::from_utf8(buf_name.to_vec()).unwrap();
    let rpc_client = RpcClient::new(rpc_url);

    // Fetch account data.
    let (ren_bridge_account_id, _) =
        Pubkey::find_program_address(&[b"RenBridgeState"], &ren_bridge::id());
    let ren_bridge_account_data = rpc_client.get_account_data(&ren_bridge_account_id).unwrap();
    let ren_bridge_state = RenBridge::unpack_unchecked(&ren_bridge_account_data).unwrap();

    ren_bridge_state.burn_count + 1
}

#[no_mangle]
pub extern "C" fn ren_bridge_mint(
    keypair_path: *const libc::c_char,
    rpc_url: *const libc::c_char,
    authority_secret: *const libc::c_char,
    selector: *const libc::c_char,
    amount: u64,
) -> *const libc::c_char {
    // Solana default signer and fee payer.
    let buf_name = unsafe { CStr::from_ptr(keypair_path).to_bytes() };
    let keypair_path = String::from_utf8(buf_name.to_vec()).unwrap();
    let payer = read_keypair_file(&keypair_path).unwrap();

    // RenVM authority secret.
    let buf_name = unsafe { CStr::from_ptr(authority_secret).to_bytes() };
    let authority_secret = String::from_utf8(buf_name.to_vec()).unwrap();
    let renvm = RenVM::from_str(&authority_secret).unwrap();
    let renvm_authority = renvm.address();

    // Initialize client.
    let buf_name = unsafe { CStr::from_ptr(rpc_url).to_bytes() };
    let rpc_url = String::from_utf8(buf_name.to_vec()).unwrap();
    let rpc_client = RpcClient::new(rpc_url);
    let commitment_config = CommitmentConfig::single_gossip();
    let (recent_blockhash, _, _) = rpc_client
        .get_recent_blockhash_with_commitment(commitment_config)
        .unwrap()
        .value;

    // Get selector hash.
    let buf_name = unsafe { CStr::from_ptr(selector).to_bytes() };
    let selector = String::from_utf8(buf_name.to_vec()).unwrap();
    let mut hasher = sha3::Keccak256::new();
    hasher.update(selector.as_bytes());
    let selector_hash: [u8; 32] = hasher.finalize().into();

    // Derived address that will be the token mint.
    let (ren_bridge_account_id, _) =
        Pubkey::find_program_address(&[b"RenBridgeState"], &ren_bridge::id());
    let (token_mint_id, _) = Pubkey::find_program_address(&[&selector_hash[..]], &ren_bridge::id());
    let (mint_authority_id, _) =
        Pubkey::find_program_address(&[&token_mint_id.to_bytes()], &ren_bridge::id());
    let associated_token_account = get_associated_token_address(&payer.pubkey(), &token_mint_id);

    // Construct RenVM mint message and sign it.
    let renvm_mint_msg = RenVmMsgBuilder::default()
        .amount(amount)
        .to(associated_token_account.to_bytes())
        .s_hash(selector_hash)
        .build()
        .unwrap();
    let msg_hash = renvm_mint_msg.get_digest().unwrap();
    let renvm_sig = renvm.sign(&renvm_mint_msg).unwrap();
    let (sig_r, sig_s, sig_v) = array_refs![&renvm_sig, 32, 32, 1];
    let (mint_log_account_id, _) =
        Pubkey::find_program_address(&[&msg_hash[..]], &ren_bridge::id());
    let mut tx = Transaction::new_with_payer(
        &[
            mint(
                &ren_bridge::id(),
                &payer.pubkey(),
                &ren_bridge_account_id,
                &token_mint_id,
                &associated_token_account,
                &mint_log_account_id,
                &mint_authority_id,
                &spl_token::id(),
            )
            .unwrap(),
            util::mint_secp_instruction(
                sig_r,
                sig_s,
                u8::from_le_bytes(*sig_v),
                &util::encode_msg(
                    renvm_mint_msg.amount,
                    &renvm_mint_msg.s_hash,
                    &renvm_mint_msg.to,
                    &renvm_mint_msg.p_hash,
                    &renvm_mint_msg.n_hash,
                ),
                renvm_authority[..].to_vec(),
            ),
        ],
        Some(&payer.pubkey()),
    );
    tx.sign(&[&payer], recent_blockhash);

    // Broadcast transaction.
    let signature = rpc_client
        .send_transaction_with_config(
            &tx,
            RpcSendTransactionConfig {
                preflight_commitment: Some(commitment_config.commitment),
                ..RpcSendTransactionConfig::default()
            },
        )
        .unwrap();

    CString::new(signature.to_string()).unwrap().into_raw()
}

#[no_mangle]
pub extern "C" fn ren_bridge_burn(
    keypair_path: *const libc::c_char,
    rpc_url: *const libc::c_char,
    selector: *const libc::c_char,
    burn_count: u64,
    burn_amount: u64,
    recipient_pointer: *const u8,
) -> *const libc::c_char {
    // Solana default signer and fee payer.
    let buf_name = unsafe { CStr::from_ptr(keypair_path).to_bytes() };
    let keypair_path = String::from_utf8(buf_name.to_vec()).unwrap();
    let payer = read_keypair_file(&keypair_path).unwrap();

    // Initialize client.
    let buf_name = unsafe { CStr::from_ptr(rpc_url).to_bytes() };
    let rpc_url = String::from_utf8(buf_name.to_vec()).unwrap();
    let rpc_client = RpcClient::new(rpc_url);
    let commitment_config = CommitmentConfig::single_gossip();
    let (recent_blockhash, _, _) = rpc_client
        .get_recent_blockhash_with_commitment(commitment_config)
        .unwrap()
        .value;

    // Get selector hash.
    let buf_name = unsafe { CStr::from_ptr(selector).to_bytes() };
    let selector = String::from_utf8(buf_name.to_vec()).unwrap();
    let mut hasher = sha3::Keccak256::new();
    hasher.update(selector.as_bytes());
    let selector_hash: [u8; 32] = hasher.finalize().into();

    // Derived address that will be the token mint.
    let (ren_bridge_account_id, _) =
        Pubkey::find_program_address(&[b"RenBridgeState"], &ren_bridge::id());
    let (token_mint_id, _) = Pubkey::find_program_address(&[&selector_hash[..]], &ren_bridge::id());
    let associated_token_account = get_associated_token_address(&payer.pubkey(), &token_mint_id);

    // Construct the 25-bytes address of release recipient of the underlying assets.
    let recipient_slice =
        unsafe { std::slice::from_raw_parts(recipient_pointer as *const u8, 25usize) };
    let mut release_recipient = [0u8; 25usize];
    release_recipient.copy_from_slice(recipient_slice);

    let (burn_log_account_id, _) =
        Pubkey::find_program_address(&[&burn_count.to_le_bytes()[..]], &ren_bridge::id());
    let mut tx = Transaction::new_with_payer(
        &[
            burn_checked(
                &spl_token::id(),
                &associated_token_account,
                &token_mint_id,
                &payer.pubkey(),
                &[],
                burn_amount,
                9u8,
            )
            .unwrap(),
            burn(
                &ren_bridge::id(),
                &payer.pubkey(),
                &associated_token_account,
                &ren_bridge_account_id,
                &token_mint_id,
                &burn_log_account_id,
                release_recipient,
            )
            .unwrap(),
        ],
        Some(&payer.pubkey()),
    );
    tx.sign(&[&payer], recent_blockhash);

    // Broadcast transaction.
    let signature = rpc_client
        .send_transaction_with_config(
            &tx,
            RpcSendTransactionConfig {
                preflight_commitment: Some(commitment_config.commitment),
                ..RpcSendTransactionConfig::default()
            },
        )
        .unwrap();

    CString::new(signature.to_string()).unwrap().into_raw()
}
