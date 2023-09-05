mod compute_aggregate_key;
mod mask;
mod reveal;
mod utils;
mod verify_key_ownership;
mod verify_reveal;
mod verify_shuffle;

use {
    self::{
        compute_aggregate_key::ComputeAggregateKey, mask::Mask, reveal::Reveal,
        verify_key_ownership::VerifyKeyOwnership, verify_reveal::VerifyReveal,
        verify_shuffle::VerifyShuffle,
    },
    crate::{Error, Result},
    ark_bn254::{g1::Config, Fr, G1Affine, G1Projective},
    ark_ec::models::short_weierstrass::{Affine, Projective},
    barnett_smart_card_protocol::discrete_log_cards::{
        DLCards, MaskedCard as InMaskedCard, Parameters, RevealToken as InRevealToken,
    },
    proof_essentials::{
        homomorphic_encryption::el_gamal::{ElGamal, Plaintext},
        vector_commitment::pedersen::PedersenCommitment,
        zkp::{
            arguments::shuffle::proof::Proof as InShuffleProof,
            proofs::{
                chaum_pedersen_dl_equality::proof::Proof as InRevealProof,
                schnorr_identification::proof::Proof as InKeyownershipProof,
            },
        },
    },
    std::slice,
};

type CConfig = Config;
type CProjective<T> = Projective<T>;
type CCurve = G1Projective;
type CCardProtocol<'a> = DLCards<'a, CCurve>;
type CCardParameters = Parameters<CProjective<CConfig>>;
type CPublicKey = Affine<CConfig>;
type CMaskedCard = InMaskedCard<CCurve>;
type CRevealToken = InRevealToken<CCurve>;
type CAggregatePublicKey = G1Affine;
type CRevealProof = InRevealProof<CCurve>;
type CShuffleProof = InShuffleProof<Fr, ElGamal<CCurve>, PedersenCommitment<CCurve>>;
type CKeyownershipProof = InKeyownershipProof<CProjective<CConfig>>;
type CCard = Plaintext<CProjective<CConfig>>;
type CScalar = Fr;

#[no_mangle]
pub extern "C" fn __precompile_mental_pokey_verify(data_ptr: *const u8, data_len: usize) -> u8 {
    let data = unsafe { slice::from_raw_parts(data_ptr, data_len) };

    if let Err(e) = check(data) {
        e.code()
    } else {
        0
    }
}

#[no_mangle]
pub extern "C" fn __precompile_mental_pokey_exec(
    data_ptr: *const u8,
    data_len: usize,
    ret_val: *mut u8,
) -> u8 {
    let data = unsafe { slice::from_raw_parts(data_ptr, data_len) };
    let ret = unsafe { slice::from_raw_parts_mut(ret_val, 32) };

    match exec(data) {
        Ok(v) => {
            ret.copy_from_slice(&v);
            0
        }
        Err(e) => e.code(),
    }
}

#[no_mangle]
pub extern "C" fn __precompile_mental_pokey_verify_gas(
    data_ptr: *const u8,
    data_len: usize,
) -> u64 {
    let data = unsafe { slice::from_raw_parts(data_ptr, data_len) };

    verify_gas(data).unwrap_or_default()
}

#[no_mangle]
pub extern "C" fn __precompile_mental_pokey_exec_gas(data_ptr: *const u8, data_len: usize) -> u64 {
    let data = unsafe { slice::from_raw_parts(data_ptr, data_len) };

    exec_gas(data).unwrap_or_default()
}

pub fn verify_gas(data: &[u8]) -> Result<u64> {
    Ok(ArgumentVerifys::new(&data[4..])?.gas())
}

pub fn exec_gas(data: &[u8]) -> Result<u64> {
    Ok(Arguments::new(data)?.gas())
}

pub fn check(data: &[u8]) -> Result<()> {
    let args = ArgumentVerifys::new(data)?;
    args.check()?;

    Ok(())
}

pub fn exec(data: &[u8]) -> Result<Vec<u8>> {
    let args = Arguments::new(data)?;
    args.exec()
}

// verifyKeyOwnership(bytes,bytes,bytes,bytes) 0x3931f649
pub const VERIFY_KEY_OWNERSHIP: [u8; 4] = [0x39, 0x31, 0xf6, 0x49];
// verifyReveal(bytes,bytes,bytes,bytes,bytes) 0x9ca80d77
pub const VERIFY_REVEAL: [u8; 4] = [0x9c, 0xa8, 0x0d, 0x77];
// computeAggregateKey(bytes[]) 0x5b2bfec7
pub const COMPUTE_AGGREGATE_KEY: [u8; 4] = [0x5b, 0x2b, 0xfe, 0xc7];
// verifyShuffle(bytes,bytes,bytes[],bytes[],bytes) 0x2a379865
pub const VERIFY_SHUFFLE: [u8; 4] = [0x2a, 0x37, 0x98, 0x65];
// reveal(bytes[],bytes) 0x6a33d652
pub const REVEAL: [u8; 4] = [0x6a, 0x33, 0xd6, 0x52];
// mask(bytes,bytes,bytes) 0x5a8890bc
pub const MASK: [u8; 4] = [0x5a, 0x88, 0x90, 0xbc];

pub enum ArgumentVerifys {
    VerifyKeyOwnership(VerifyKeyOwnership),
    VerifyReveal(VerifyReveal),
    VerifyShuffle(VerifyShuffle),
}

impl ArgumentVerifys {
    pub fn new(data: &[u8]) -> Result<Self> {
        if data.len() < 4 {
            return Err(Error::WrongSelectorLength);
        }
        match [data[0], data[1], data[2], data[3]] {
            VERIFY_KEY_OWNERSHIP => Ok(Self::VerifyKeyOwnership(VerifyKeyOwnership::new(
                &data[4..],
            )?)),
            VERIFY_REVEAL => Ok(Self::VerifyReveal(VerifyReveal::new(&data[4..])?)),
            VERIFY_SHUFFLE => Ok(Self::VerifyShuffle(VerifyShuffle::new(&data[4..])?)),
            _ => Err(Error::UnknownSelector),
        }
    }

    pub fn check(self) -> Result<()> {
        match self {
            Self::VerifyKeyOwnership(v) => v.check(),
            Self::VerifyReveal(v) => v.check(),
            Self::VerifyShuffle(v) => v.check(),
        }
    }

    pub fn gas(self) -> u64 {
        match self {
            Self::VerifyKeyOwnership(v) => v.gas(),
            Self::VerifyReveal(v) => v.gas(),
            Self::VerifyShuffle(v) => v.gas(),
        }
    }
}

pub enum Arguments {
    ComputeAggregateKey(ComputeAggregateKey),
    Reveal(Reveal),
    Mask(Mask),
}

impl Arguments {
    pub fn new(data: &[u8]) -> Result<Self> {
        if data.len() < 4 {
            return Err(Error::WrongSelectorLength);
        }
        match [data[0], data[1], data[2], data[3]] {
            COMPUTE_AGGREGATE_KEY => Ok(Self::ComputeAggregateKey(ComputeAggregateKey::new(
                &data[4..],
            )?)),
            REVEAL => Ok(Self::Reveal(Reveal::new(&data[4..])?)),
            MASK => Ok(Self::Mask(Mask::new(&data[4..])?)),
            _ => Err(Error::UnknownSelector),
        }
    }

    pub fn exec(self) -> Result<Vec<u8>> {
        match self {
            Self::ComputeAggregateKey(v) => v.check(),
            Self::Reveal(v) => v.check(),
            Self::Mask(v) => v.check(),
        }
    }

    pub fn gas(self) -> u64 {
        match self {
            Self::ComputeAggregateKey(v) => v.gas(),
            Self::Reveal(v) => v.gas(),
            Self::Mask(v) => v.gas(),
        }
    }
}
