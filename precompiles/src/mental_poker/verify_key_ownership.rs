use {
    super::{utils::deserialize, CCardParameters, CCardProtocol, CKeyownershipProof, CPublicKey},
    crate::{utils, Error, Result},
    barnett_smart_card_protocol::BarnettSmartProtocol,
    ethabi::ParamType,
};

pub struct VerifyKeyOwnership {
    params: Vec<u8>,
    pub_key: Vec<u8>,
    memo: Vec<u8>,
    key_proof: Vec<u8>,
}

const VERIFY_KEY_OWNERSHIP_PER_GAS: u64 = 50000;

impl VerifyKeyOwnership {
    fn params_type() -> [ParamType; 4] {
        let bytes = ParamType::Bytes;
        /*
            bytes params
            bytes pub_key
            bytes memo
            bytes key_proof
        */
        [bytes.clone(), bytes.clone(), bytes.clone(), bytes]
    }

    pub fn new(data: &[u8]) -> Result<Self> {
        let res = ethabi::decode(&Self::params_type(), data).map_err(|_| Error::ParseDataFailed)?;

        let params = utils::into_bytes(res.get(4).cloned())?;
        let pub_key = utils::into_bytes(res.get(4).cloned())?;
        let memo = utils::into_bytes(res.get(4).cloned())?;
        let key_proof = utils::into_bytes(res.get(4).cloned())?;

        Ok(Self {
            params,
            pub_key,
            memo,
            key_proof,
        })
    }

    pub fn check(self) -> Result<()> {
        let params: CCardParameters = deserialize(&self.params)?;
        let pub_key: CPublicKey = deserialize(&self.pub_key)?;
        let memo: Vec<u8> = deserialize(&self.memo)?;
        let key_proof: CKeyownershipProof = deserialize(&self.key_proof)?;

        CCardProtocol::verify_key_ownership(&params, &pub_key, &memo.to_vec(), &key_proof)
            .map_err(|_| Error::ProofVerificationFailed)
    }

    pub fn gas(self) -> u64 {
        VERIFY_KEY_OWNERSHIP_PER_GAS
    }
}
