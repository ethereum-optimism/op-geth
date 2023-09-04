use crate::{utils, Error, Result};
use ethabi::ParamType;
use hypr_algebra::{bn254::BN254Scalar, serialization::FromToBytes};
use hypr_crypto::anemoi_jive::{AnemoiJive, AnemoiJive254, ANEMOI_JIVE_BN254_SALTS};
use num_bigint::BigUint;
use std::slice;

#[no_mangle]
pub extern "C" fn __precompile_anemoi(
    data_ptr: *const u8,
    data_len: usize,
    ret_val: *mut u8,
) -> u8 {
    let data = unsafe { slice::from_raw_parts(data_ptr, data_len) };
    let ret = unsafe { slice::from_raw_parts_mut(ret_val, 32) };

    match compute(data, ret) {
        Ok(()) => 0,
        Err(e) => e.code(),
    }
}

#[no_mangle]
pub extern "C" fn __precompile_anemoi_gas(data_ptr: *const u8, data_len: usize) -> u64 {
    let data = unsafe { slice::from_raw_parts(data_ptr, data_len) };

    gas(data).unwrap_or_default()
}

fn eval_jive4(data: &[u8], ret: &mut [u8]) -> Result<()> {
    let param = ParamType::FixedBytes(32);

    let r = ethabi::decode(&[param.clone(), param.clone(), param.clone(), param], data)
        .map_err(|_| Error::ParseDataFailed)?;

    let x0 = utils::into_bytes32(r.get(0).cloned())?;
    let x1 = utils::into_bytes32(r.get(1).cloned())?;
    let y0 = utils::into_bytes32(r.get(2).cloned())?;
    let y1 = utils::into_bytes32(r.get(3).cloned())?;

    let res = AnemoiJive254::eval_jive(
        &[
            BN254Scalar::from_bytes(&x0).map_err(|_| Error::ParseDataFailed)?,
            BN254Scalar::from_bytes(&x1).map_err(|_| Error::ParseDataFailed)?,
        ],
        &[
            BN254Scalar::from_bytes(&y0).map_err(|_| Error::ParseDataFailed)?,
            BN254Scalar::from_bytes(&y1).map_err(|_| Error::ParseDataFailed)?,
        ],
    );

    let r = BigUint::from(res).to_bytes_le();

    ret.copy_from_slice(&r);

    Ok(())
}

fn jive_254_salts(data: &[u8], ret: &mut [u8]) -> Result<()> {
    let r = ethabi::decode(&[ParamType::Uint(32)], data).map_err(|_| Error::ParseDataFailed)?;

    let index = utils::into_uint(r.get(0).cloned())?;

    if index < 64 {
        let point = BigUint::from(ANEMOI_JIVE_BN254_SALTS[index as usize]).to_bytes_le();

        ret.copy_from_slice(&point);
        Ok(())
    } else {
        Err(Error::InputOutOfBound)
    }
}

// anemoi_jive_4(bytes32,bytes32,bytes32,bytes32) 0x73808263
pub const ANEMOI_JIVE_4_SELECTOR: [u8; 4] = [0x73, 0x80, 0x82, 0x63];
// anemoi_jive_254_salts(uint128) 0xbc4f54ce
pub const ANEMOI_JIVE_254_SALTS_SELECTOR: [u8; 4] = [0xbc, 0x4f, 0x54, 0xce];

fn compute(data: &[u8], ret: &mut [u8]) -> Result<()> {
    if data.len() < 4 {
        return Err(Error::WrongSelectorLength);
    }

    match [data[0], data[1], data[2], data[3]] {
        ANEMOI_JIVE_4_SELECTOR => eval_jive4(&data[4..], ret)?,
        ANEMOI_JIVE_254_SALTS_SELECTOR => jive_254_salts(&data[4..], ret)?,
        _ => return Err(Error::UnknownSelector),
    }

    Ok(())
}

pub const ANEMOI_SALT_GAS: u64 = 10;
pub const ANEMOI_EVAL_4: u64 = 400;

fn gas(data: &[u8]) -> Result<u64> {
    if data.len() < 4 {
        return Err(Error::WrongSelectorLength);
    }

    match [data[0], data[1], data[2], data[3]] {
        ANEMOI_JIVE_4_SELECTOR => Ok(ANEMOI_SALT_GAS),
        ANEMOI_JIVE_254_SALTS_SELECTOR => Ok(ANEMOI_EVAL_4),
        _ => Err(Error::UnknownSelector),
    }
}
