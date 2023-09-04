use {
    super::{
        utils::deserialize, CCardParameters, CCardProtocol, CMaskedCard, CPublicKey, CShuffleProof,
    },
    crate::{utils, Error, Result},
    barnett_smart_card_protocol::BarnettSmartProtocol,
    ethabi::ParamType,
};

pub const VERIFY_SHUFFLE_PER_GAS: u64 = 50000;

pub struct VerifyShuffle {
    params: Vec<u8>,
    shared_key: Vec<u8>,
    cur_decks: Vec<Vec<u8>>,
    new_decks: Vec<Vec<u8>>,
    shuffle_proof: Vec<u8>,
}

impl VerifyShuffle {
    fn params_type() -> [ParamType; 5] {
        let bytes = ParamType::Bytes;
        let bytes_array = ParamType::Array(Box::new(ParamType::Bytes));
        [
            bytes.clone(),
            bytes.clone(),
            bytes_array.clone(),
            bytes_array,
            bytes,
        ]
    }

    pub fn new(data: &[u8]) -> Result<Self> {
        let res = ethabi::decode(&Self::params_type(), data).map_err(|_| Error::ParseDataFailed)?;

        let params = utils::into_bytes(res.get(0).cloned())?;
        let shared_key = utils::into_bytes(res.get(1).cloned())?;
        let cur_decks = utils::into_bytes_array(res.get(3).cloned())?;
        let new_decks = utils::into_bytes_array(res.get(4).cloned())?;
        let shuffle_proof = utils::into_bytes(res.get(5).cloned())?;

        Ok(Self {
            params,
            shared_key,
            cur_decks,
            new_decks,
            shuffle_proof,
        })
    }

    pub fn check(self) -> Result<()> {
        let params: CCardParameters = deserialize(&self.params)?;
        let shared_key: CPublicKey = deserialize(&self.shared_key)?;
        let mut cur_decks3: Vec<CMaskedCard> = Vec::new();
        for v_cur_deck in &self.cur_decks {
            let v_cur_deck: CMaskedCard = deserialize(v_cur_deck)?;
            cur_decks3.push(v_cur_deck);
        }
        let mut new_decks3: Vec<CMaskedCard> = Vec::new();
        for v_new_deck in &self.new_decks {
            let v_new_deck: CMaskedCard = deserialize(v_new_deck)?;
            new_decks3.push(v_new_deck);
        }
        let shuffle_proof: CShuffleProof = deserialize(&self.shuffle_proof)?;

        CCardProtocol::verify_shuffle(
            &params,
            &shared_key,
            &cur_decks3,
            &new_decks3,
            &shuffle_proof,
        )
        .map_err(|_| Error::ProofVerificationFailed)
    }

    pub fn gas(self) -> u64 {
        VERIFY_SHUFFLE_PER_GAS
    }
}
