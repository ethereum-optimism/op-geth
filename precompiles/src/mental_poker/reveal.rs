use {
    super::{
        utils::{deserialize, serialize},
        CMaskedCard, CRevealToken,
    },
    crate::{utils, Error, Result},
    ark_std::Zero,
    barnett_smart_card_protocol::Reveal as BReveal,
    ethabi::ParamType,
};

pub const REVEAL_PER_GAS: u64 = 50000;

pub struct Reveal {
    reveal_tokens: Vec<Vec<u8>>,
    masked: Vec<u8>,
}

impl Reveal {
    fn params_type() -> [ParamType; 2] {
        let bytes_array = ParamType::Array(Box::new(ParamType::Bytes));
        let bytes = ParamType::Bytes;
        [bytes_array, bytes]
    }

    pub fn new(data: &[u8]) -> Result<Self> {
        let res = ethabi::decode(&Self::params_type(), data).map_err(|_| Error::ParseDataFailed)?;

        let reveal_tokens = utils::into_bytes_array(res.get(0).cloned())?;
        let masked = utils::into_bytes(res.get(1).cloned())?;

        Ok(Self {
            reveal_tokens,
            masked,
        })
    }

    pub fn check(self) -> Result<Vec<u8>> {
        let mut aggregate_reveal_token = CRevealToken::zero();
        for reveal_token in self.reveal_tokens {
            let reveal_token: CRevealToken = deserialize(reveal_token.as_slice())?;
            aggregate_reveal_token = aggregate_reveal_token + reveal_token;
        }
        let masked: CMaskedCard = deserialize(&self.masked)?;

        let decrypted = aggregate_reveal_token
            .reveal(&masked)
            .map_err(|_| Error::ExecError)?;

        serialize(&decrypted)
    }

    pub fn gas(self) -> u64 {
        REVEAL_PER_GAS
    }
}
