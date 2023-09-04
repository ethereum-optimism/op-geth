use {
    super::{
        utils::{deserialize, serialize},
        CAggregatePublicKey, CPublicKey,
    },
    crate::{utils, Error, Result},
    ark_ec::AffineRepr,
    ethabi::ParamType,
};

pub const COMPUTE_AGGREGATE_KEY_PER_GAS: u64 = 50000;

pub struct ComputeAggregateKey {
    pub_keys: Vec<Vec<u8>>,
}

impl ComputeAggregateKey {
    fn params_type() -> [ParamType; 1] {
        let bytes_array = ParamType::Array(Box::new(ParamType::Bytes));
        [bytes_array]
    }

    pub fn new(data: &[u8]) -> Result<Self> {
        let res = ethabi::decode(&Self::params_type(), data).map_err(|_| Error::ParseDataFailed)?;

        let pub_keys = utils::into_bytes_array(res.get(0).cloned())?;

        Ok(Self { pub_keys })
    }

    pub fn check(self) -> Result<Vec<u8>> {
        let mut aggregate_pub_key = CAggregatePublicKey::zero();
        for v_pub_key in self.pub_keys {
            let v_pub_key: CPublicKey = deserialize(v_pub_key.as_slice())?;
            aggregate_pub_key = (aggregate_pub_key + v_pub_key).into();
        }

        Ok(serialize(&aggregate_pub_key).map_err(|_| Error::ProofVerificationFailed)?)
    }

    pub fn gas(self) -> u64 {
        COMPUTE_AGGREGATE_KEY_PER_GAS
    }
}
