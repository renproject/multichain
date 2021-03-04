use gateway::types::{
    RenVmMintMessage, Secp256k1InstructionData, RENVM_MINT_MESSAGE_SIZE, RENVM_MINT_SECP_DATA_SIZE,
};
use solana_sdk::{
    instruction::Instruction,
    program_pack::Pack,
    secp256k1_instruction::{
        SecpSignatureOffsets, HASHED_PUBKEY_SERIALIZED_SIZE, SIGNATURE_OFFSETS_SERIALIZED_SIZE,
    },
};

/// Constructs a Secp256k1 instruction for RenVM mint to be verified by Solana Secp256k1 program.
pub fn mint_secp_instruction(
    sig_r: &[u8; 32],
    sig_s: &[u8; 32],
    sig_v: u8,
    message_arr: &[u8],
    eth_addr: Vec<u8>,
) -> Instruction {
    // Assert that the data we've received is of the correct size.
    assert_eq!(message_arr.len(), RENVM_MINT_MESSAGE_SIZE);
    assert_eq!(
        eth_addr.len() + sig_r.len() + sig_s.len() + 1 + message_arr.len(),
        RENVM_MINT_SECP_DATA_SIZE - 1,
    );

    // Allocate appropriate size for our instruction data.
    let mut instruction_data = vec![];
    let data_start = 1 + SIGNATURE_OFFSETS_SERIALIZED_SIZE;
    let total_size = data_start + RENVM_MINT_SECP_DATA_SIZE;
    instruction_data.resize(total_size, 0);

    // Calculate the offsets for a single ECDSA signature.
    let num_signatures = 1;
    instruction_data[0] = num_signatures;
    let eth_addr_offset = data_start + 1;
    let signature_offset = eth_addr_offset + eth_addr.len();
    let message_data_offset = signature_offset + sig_r.len() + sig_s.len() + 1;

    // Copy data from slice into sized arrays.
    let mut addr = [0u8; HASHED_PUBKEY_SERIALIZED_SIZE];
    addr.copy_from_slice(&eth_addr[..]);
    let mut msg = [0u8; RENVM_MINT_MESSAGE_SIZE];
    msg.copy_from_slice(message_arr);

    // Write Secp256k1InstructionData data.
    let secp256k1_instruction_data = Secp256k1InstructionData::MintSignature {
        eth_addr: addr,
        sig_r: *sig_r,
        sig_s: *sig_s,
        sig_v: sig_v - 27,
        msg: msg,
    };
    let packed_data = secp256k1_instruction_data.pack();
    instruction_data[data_start..total_size].copy_from_slice(packed_data.as_slice());

    // Write offsets data.
    let offsets = SecpSignatureOffsets {
        signature_offset: signature_offset as u16,
        signature_instruction_index: 1,
        eth_address_offset: eth_addr_offset as u16,
        eth_address_instruction_index: 1,
        message_data_offset: message_data_offset as u16,
        message_data_size: message_arr.len() as u16,
        message_instruction_index: 1,
    };
    let writer = std::io::Cursor::new(&mut instruction_data[1..data_start]);
    bincode::serialize_into(writer, &offsets).unwrap();

    Instruction {
        program_id: solana_sdk::secp256k1_program::id(),
        accounts: vec![],
        data: instruction_data,
    }
}

/// ABI-encode the values for creating the signature hash.
pub fn encode_msg(
    amount: u64,
    shash: &[u8; 32],
    to: &[u8; 32],
    p_hash: &[u8; 32],
    n_hash: &[u8; 32],
) -> Vec<u8> {
    let mut encoded_msg = vec![0u8; RENVM_MINT_MESSAGE_SIZE];
    let msg = RenVmMintMessage {
        p_hash: *p_hash,
        amount: amount,
        selector_hash: *shash,
        to: *to,
        n_hash: *n_hash,
    };
    msg.pack_into_slice(encoded_msg.as_mut_slice());
    encoded_msg
}
