use {
    super::{
        utils::deserialize, CCardParameters, CCardProtocol, CMaskedCard, CPublicKey, CRevealProof,
        CRevealToken,
    },
    crate::{utils, Error, Result},
    barnett_smart_card_protocol::BarnettSmartProtocol,
    ethabi::ParamType,
};

pub const VERIFY_REVEAL_PER_GAS: u64 = 50000;
pub struct VerifyReveal {
    params: Vec<u8>,
    pub_key: Vec<u8>,
    reveal_token: Vec<u8>,
    masked: Vec<u8>,
    reveal_proof: Vec<u8>,
}
impl VerifyReveal {
    fn params_type() -> [ParamType; 5] {
        let bytes = ParamType::Bytes;
        /*
            bytes params
            bytes pub_key
            bytes reveal_token
            bytes masked
            bytes reveal_proof
        */
        [
            bytes.clone(),
            bytes.clone(),
            bytes.clone(),
            bytes.clone(),
            bytes,
        ]
    }

    pub fn new(data: &[u8]) -> Result<Self> {
        let res = ethabi::decode(&Self::params_type(), data).map_err(|_| Error::ParseDataFailed)?;

        let params = utils::into_bytes(res.get(0).cloned())?;
        let pub_key = utils::into_bytes(res.get(1).cloned())?;
        let reveal_token = utils::into_bytes(res.get(2).cloned())?;
        let masked = utils::into_bytes(res.get(3).cloned())?;
        let reveal_proof = utils::into_bytes(res.get(4).cloned())?;

        Ok(Self {
            params,
            pub_key,
            reveal_token,
            masked,
            reveal_proof,
        })
    }

    pub fn check(self) -> Result<()> {
        let params: CCardParameters = deserialize(&self.params)?;
        let pub_key: CPublicKey = deserialize(&self.pub_key)?;
        let reveal_token: CRevealToken = deserialize(&self.reveal_token)?;
        let masked: CMaskedCard = deserialize(&self.masked)?;
        let reveal_proof: CRevealProof = deserialize(&self.reveal_proof)?;

        CCardProtocol::verify_reveal(&params, &pub_key, &reveal_token, &masked, &reveal_proof)
            .map_err(|_| Error::ProofVerificationFailed)
    }

    pub fn gas(self) -> u64 {
        VERIFY_REVEAL_PER_GAS
    }
}
