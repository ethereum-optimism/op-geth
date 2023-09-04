use crate::{
    utils::{self, bytes_asset, check_address_format_from_folding},
    Error, Result,
};
use ethabi::ParamType;
use hypr_algebra::{bn254::BN254Scalar, serialization::FromToBytes};
use hypr_api::{
    anon_xfr::{
        ownership::{verify_ownership_note, OwnershipBody, OwnershipNote},
        AXfrAddressFoldingInstance, AXfrPlonkPf,
    },
    parameters::VerifierParams,
};
use sha3::{Digest, Sha3_512};

pub struct OwnerShip {
    root_version: u64,
    amount: u128,
    asset: Vec<u8>,
    nullifier: Vec<u8>,
    proof: Vec<u8>,
    // folding_instance: Vec<u8>,
    merkle_root: [u8; 32],
    hash: Vec<u8>,
}

impl OwnerShip {
    fn params_type() -> [ParamType; 7] {
        let bytes32 = ParamType::FixedBytes(32);
        let uint256 = ParamType::Uint(256);
        let bytes = ParamType::Bytes;

        /*
            uint256 root_version,
            uint256 amount,
            bytes asset
            bytes nullifier
            bytes proof
            bytes folding_instance
            bytes32 merkle_root
            bytes hash
        */
        [
            uint256.clone(),
            uint256,
            bytes.clone(),
            bytes.clone(),
            bytes.clone(),
            bytes32,
            bytes,
        ]
    }

    pub fn new(data: &[u8]) -> Result<Self> {
        let res = ethabi::decode(&Self::params_type(), data).map_err(|_| Error::ParseDataFailed)?;

        let root_version = utils::into_uint(res.get(0).cloned())? as u64;
        let amount = utils::into_uint(res.get(1).cloned())?;
        let asset = utils::into_bytes(res.get(2).cloned())?;
        let nullifier = utils::into_bytes(res.get(3).cloned())?;
        let proof = utils::into_bytes(res.get(4).cloned())?;
        // let folding_instance = utils::into_bytes(res.get(5).cloned())?;
        let merkle_root = utils::into_bytes32(res.get(6).cloned())?;
        let hash = utils::into_bytes(res.get(7).cloned())?;

        let r = Self {
            root_version,
            amount,
            asset,
            nullifier,
            proof,
            // folding_instance,
            merkle_root,
            hash,
        };

        Ok(r)
    }

    pub fn check(self) -> Result<()> {
        verify_ownership(
            self.root_version,
            self.amount,
            self.asset,
            self.nullifier,
            self.proof,
            self.merkle_root,
            self.hash,
        )
    }
    pub fn gas(self) -> u64 {
        OWNERSHIP_PER_FEE
    }
}

pub const OWNERSHIP_PER_FEE: u64 = 75000;

#[allow(clippy::too_many_arguments)]
fn verify_ownership(
    root_version: u64,
    amount: u128,
    asset: Vec<u8>,
    nullifier: Vec<u8>,
    proof: Vec<u8>,
    merkle_root: [u8; 32],
    hash: Vec<u8>,
) -> Result<()> {
    let (proof, folding_instance): (AXfrPlonkPf, AXfrAddressFoldingInstance) = bincode::deserialize(&proof).map_err(|_| Error::ProofDecodeFailed)?;
    let address_format = check_address_format_from_folding(&folding_instance);
    let asset = bytes_asset(&asset)?;
    let input = BN254Scalar::from_bytes(&nullifier).map_err(|_| Error::ParseDataFailed)?;
    let merkle_root = BN254Scalar::from_bytes(&merkle_root).map_err(|_| Error::ParseDataFailed)?;

    let note = OwnershipNote {
        body: OwnershipBody {
            input,
            asset,
            amount,
            merkle_root,
            merkle_root_version: root_version,
        },
        proof,
        folding_instance,
    };

    let mut hasher = Sha3_512::new();
    hasher.update(&hash);
    hasher.update(&bincode::serialize(&note.body).map_err(|_| Error::ParseDataFailed)?);

    let params = VerifierParams::get_ownership(address_format)
        .map_err(|_| Error::FailedToLoadVerifierParams)?;

    verify_ownership_note(&params, &note, &merkle_root, hasher)
        .map_err(|_| Error::ProofVerificationFailed)
}
