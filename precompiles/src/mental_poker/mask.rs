use barnett_smart_card_protocol::BarnettSmartProtocol;
use ethabi::ParamType;

use super::{
    utils::{deserialize, serialize},
    CAggregatePublicKey, CCard, CCardParameters, CCardProtocol, CMaskedCard, CScalar,
};
use crate::{utils, Error, Result};
use ark_std::One;

pub const MASK_PER_GAS: u64 = 50000;
pub struct Mask {
    params: Vec<u8>,
    shared_key: Vec<u8>,
    encoded: Vec<u8>,
}

impl Mask {
    fn params_type() -> [ParamType; 3] {
        let bytes = ParamType::Bytes;
        /*
            bytes params
            bytes shared_key
            bytes encoded
        */
        [bytes.clone(), bytes.clone(), bytes]
    }

    pub fn new(data: &[u8]) -> Result<Self> {
        let res = ethabi::decode(&Self::params_type(), data).map_err(|_| Error::ParseDataFailed)?;

        let params = utils::into_bytes(res.get(0).cloned())?;
        let shared_key = utils::into_bytes(res.get(1).cloned())?;
        let encoded = utils::into_bytes(res.get(2).cloned())?;

        Ok(Self {
            params,
            shared_key,
            encoded,
        })
    }

    pub fn check(self) -> Result<Vec<u8>> {
        let params: CCardParameters = deserialize(&self.params)?;
        let shared_key: CAggregatePublicKey = deserialize(&self.shared_key)?;
        let encoded: CCard = deserialize(&self.encoded)?;

        let masked: CMaskedCard =
            CCardProtocol::mask_only(&params, &shared_key, &encoded, &CScalar::one())
                .map_err(|_| Error::ProofVerificationFailed)?;
        serialize(&masked)
    }

    pub fn gas(self) -> u64 {
        MASK_PER_GAS
    }
}
