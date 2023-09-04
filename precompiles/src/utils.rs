use ethabi::{Address, Token};
use ethereum_types::U256;
use hypr_api::{
    anon_xfr::AXfrAddressFoldingInstance,
    parameters::AddressFormat,
    structs::{AssetType, ASSET_TYPE_LENGTH},
};

use crate::{Error, Result};

pub fn into_bytes32(tk: Option<Token>) -> Result<[u8; 32]> {
    tk.ok_or(Error::ParseDataFailed)?
        .into_fixed_bytes()
        .ok_or(Error::ParseDataFailed)?
        .try_into()
        .map_err(|_| Error::ParseDataFailed)
}

pub fn into_bytes32_array(tk: Option<Token>) -> Result<Vec<[u8; 32]>> {
    let token = tk.ok_or(Error::ParseDataFailed)?;
    let arr = token.into_array().ok_or(Error::ParseDataFailed)?;

    let mut res = Vec::with_capacity(arr.len());

    for item in arr {
        let bytes32 = into_bytes32(Some(item))?;

        res.push(bytes32)
    }

    Ok(res)
}

pub fn into_address(tk: Option<Token>) -> Result<Address> {
    Ok(tk
        .ok_or(Error::ParseDataFailed)?
        .into_address()
        .ok_or(Error::ParseDataFailed)?)
}

pub fn into_uint(tk: Option<Token>) -> Result<u128> {
    Ok(tk
        .ok_or(Error::ParseDataFailed)?
        .into_uint()
        .ok_or(Error::ParseDataFailed)?
        .as_u128())
}

pub fn into_uint256(tk: Option<Token>) -> Result<U256> {
    Ok(tk
        .ok_or(Error::ParseDataFailed)?
        .into_uint()
        .ok_or(Error::ParseDataFailed)?)
}

pub fn into_uint_array(tk: Option<Token>) -> Result<Vec<u128>> {
    let token = tk.ok_or(Error::ParseDataFailed)?;
    let arr = token.into_array().ok_or(Error::ParseDataFailed)?;

    let mut res = Vec::with_capacity(arr.len());

    for item in arr {
        let uint = into_uint(Some(item))?;

        res.push(uint)
    }

    Ok(res)
}

pub fn into_uint256_array(tk: Option<Token>) -> Result<Vec<U256>> {
    let token = tk.ok_or(Error::ParseDataFailed)?;
    let arr = token.into_array().ok_or(Error::ParseDataFailed)?;

    let mut res = Vec::with_capacity(arr.len());

    for item in arr {
        let uint = into_uint256(Some(item))?;

        res.push(uint)
    }

    Ok(res)
}

pub fn into_bytes(tk: Option<Token>) -> Result<Vec<u8>> {
    tk.ok_or(Error::ParseDataFailed)?
        .into_bytes()
        .ok_or(Error::ParseDataFailed)
}

pub fn into_bytes_array(tk: Option<Token>) -> Result<Vec<Vec<u8>>> {
    let token = tk.ok_or(Error::ParseDataFailed)?;
    let arr = token.into_array().ok_or(Error::ParseDataFailed)?;

    let mut res = Vec::with_capacity(arr.len());

    for item in arr {
        let bytes = into_bytes(Some(item))?;

        res.push(bytes)
    }

    Ok(res)
}

pub fn join_bytes32(byte32s: &[[u8; 32]]) -> Vec<u8> {
    let mut v = Vec::with_capacity(byte32s.len() * 32);

    for b in byte32s {
        v.extend_from_slice(b);
    }

    v
}

pub fn split_bytes32(bytes32: &Vec<u8>) -> Result<Vec<&[u8; 32]>> {
    let num = bytes32.len() / 32;

    let mut res = Vec::with_capacity(num);

    for i in 0..num {
        let begin = i * 32;
        let end = (i + 1) * 32;

        let b32 = bytes32.get(begin..end).ok_or(Error::ParseDataFailed)?;

        res.push(b32.try_into().map_err(|_| Error::ParseDataFailed)?)
    }

    Ok(res)
}
pub fn bytes_asset(bytes: &[u8]) -> Result<AssetType> {
    if bytes.len() != ASSET_TYPE_LENGTH {
        return Err(Error::ParseDataFailed);
    }
    let mut asset_bytes = [0u8; ASSET_TYPE_LENGTH];
    asset_bytes.copy_from_slice(bytes);
    Ok(AssetType(asset_bytes))
}
pub fn check_address_format_from_folding(folding: &AXfrAddressFoldingInstance) -> AddressFormat {
    match folding {
        AXfrAddressFoldingInstance::Ed25519(_) => AddressFormat::ED25519,
        AXfrAddressFoldingInstance::Secp256k1(_) => AddressFormat::SECP256K1,
    }
}
