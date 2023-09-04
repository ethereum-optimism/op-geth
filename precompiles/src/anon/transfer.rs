use ethabi::{Address, ParamType, Token};
use ethereum_types::U256;
use hypr_algebra::{bn254::BN254Scalar, serialization::FromToBytes};
use hypr_api::{
    anon_xfr::{
        abar_to_abar::{verify_anon_xfr_note, AXfrBody, AXfrNote},
        AXfrAddressFoldingInstance, AXfrPlonkPf,
    },
    parameters::VerifierParams,
    structs::{AnonAssetRecord, AxfrOwnerMemo},
};
use sha3::{Digest, Sha3_512};

use crate::{
    utils::{self, bytes_asset, check_address_format_from_folding},
    Error, Result,
};

pub struct TransferEntity {
    root_version: u64,
    root: [u8; 32],
    asset: [u8; 32],
    fee_amount: U256,
    transparent_account: Address,
    transparent_asset: [u8; 32],
    transparent_amount: U256,
    hash: [u8; 32],
    nullifiers: Vec<[u8; 32]>,
    commitments: Vec<[u8; 32]>,
    memos: Vec<Vec<u8>>,
    proof: Vec<u8>,
}

pub struct Transfer {
    params: Vec<TransferEntity>,
}

impl Transfer {
    // abi "tuple(uint64,bytes32,bytes32,uint256,address,bytes32,uint256,bytes32,bytes32[],bytes32[],bytes[],bytes)[]"
    fn params_type() -> Vec<ParamType> {
        vec![
            ParamType::Uint(64),
            ParamType::FixedBytes(32),
            ParamType::FixedBytes(32),
            ParamType::Uint(256),
            ParamType::Address,
            ParamType::FixedBytes(32),
            ParamType::Uint(256),
            ParamType::FixedBytes(32),
            ParamType::Array(Box::new(ParamType::FixedBytes(32))),
            ParamType::Array(Box::new(ParamType::FixedBytes(32))),
            ParamType::Array(Box::new(ParamType::Bytes)),
            ParamType::Bytes,
        ]
    }

    pub fn new(data: &[u8]) -> Result<Self> {
        let res = ethabi::decode(
            &[ParamType::Array(Box::new(ParamType::Tuple(
                Self::params_type(),
            )))],
            data,
        )
        .map_err(|_| Error::ParseDataFailed)?;

        assert_eq!(1, res.len());

        let mut result = Vec::with_capacity(res.len());

        match res.get(0) {
            Some(base) => match base {
                Token::Array(array_base) => {
                    for tk in array_base {
                        match tk {
                            Token::Tuple(inner) => {
                                let root_version = utils::into_uint(inner.get(0).cloned())? as u64;
                                let root = utils::into_bytes32(inner.get(1).cloned())?;
                                let asset = utils::into_bytes32(inner.get(2).cloned())?;
                                let fee_amount = utils::into_uint256(inner.get(3).cloned())?;
                                let transparent_account =
                                    utils::into_address(inner.get(4).cloned())?;
                                let transparent_asset = utils::into_bytes32(inner.get(5).cloned())?;
                                let transparent_amount =
                                    utils::into_uint256(inner.get(6).cloned())?;
                                let hash = utils::into_bytes32(inner.get(7).cloned())?;
                                let nullifiers = utils::into_bytes32_array(inner.get(8).cloned())?;
                                let commitments = utils::into_bytes32_array(inner.get(9).cloned())?;
                                let memos = utils::into_bytes_array(inner.get(10).cloned())?;
                                let proof = utils::into_bytes(inner.get(11).cloned())?;
                                result.push(TransferEntity {
                                    root_version,
                                    root,
                                    asset,
                                    fee_amount,
                                    transparent_account,
                                    transparent_asset,
                                    transparent_amount,
                                    hash,
                                    nullifiers,
                                    commitments,
                                    memos,
                                    proof,
                                });
                            }
                            _ => {
                                return Err(Error::ParseDataFailed);
                            }
                        }
                    }
                }
                _ => {
                    return Err(Error::ParseDataFailed);
                }
            },
            _ => return Err(Error::ParseDataFailed),
        }

        Ok(Self { params: result })
    }

    pub fn check(self) -> Result<()> {
        let mut res = Vec::new();
        self.params.into_iter().for_each(|x| {
            res.push(verify_atoa(
                x.nullifiers,
                x.commitments,
                x.root,
                x.proof,
                x.fee_amount.as_u128(),
                x.asset,
                x.transparent_asset,
                x.root_version,
                x.transparent_amount.as_u128(),
                x.memos,
                x.hash,
            ));
        });
        for r in res {
            r?
        }
        Ok(())
    }

    pub fn gas(self) -> u64 {
        let mut gas: u64 = 0;
        self.params.into_iter().for_each(|x| {
            gas += TRANSFER_PER_INPUT * x.nullifiers.len() as u64
                + TRANSFER_PER_OUTPUT * x.commitments.len() as u64
        });
        gas
    }
}

pub const TRANSFER_PER_INPUT: u64 = 4000;
pub const TRANSFER_PER_OUTPUT: u64 = 30000;

fn verify_atoa(
    nullifiers: Vec<[u8; 32]>,
    commitments: Vec<[u8; 32]>,
    merkle_root: [u8; 32],
    proof: Vec<u8>,
    fee: u128,
    fee_asset: [u8; 32],
    transparent_asset: [u8; 32],
    root_version: u64,
    transparent: u128,
    memos: Vec<Vec<u8>>,
    hash: [u8; 32],
) -> Result<()> {
    let (proof, folding_instance): (AXfrPlonkPf, AXfrAddressFoldingInstance) =
        bincode::deserialize(&proof).map_err(|_| Error::ProofDecodeFailed)?;

    let address_format = check_address_format_from_folding(&folding_instance);
    let fee_type = bytes_asset(&fee_asset)?;
    let transparent_type = bytes_asset(&transparent_asset)?;
    let merkle_root = BN254Scalar::from_bytes(&merkle_root).map_err(|_| Error::ParseDataFailed)?;

    let mut inputs = vec![];
    for bytes in nullifiers.iter() {
        inputs.push(BN254Scalar::from_bytes(bytes).map_err(|_| Error::ParseDataFailed)?);
    }

    let mut outputs = vec![];
    for bytes in commitments.iter() {
        outputs.push(AnonAssetRecord {
            commitment: BN254Scalar::from_bytes(bytes).map_err(|_| Error::ParseDataFailed)?,
        });
    }
    let (inputs_len, outputs_len) = (inputs.len(), outputs.len());
    let note = AXfrNote {
        body: AXfrBody {
            inputs,
            outputs,
            merkle_root,
            merkle_root_version: root_version,
            fee,
            fee_type,
            transparent,
            transparent_type,
            owner_memos: memos
                .iter()
                .map(|bytes| AxfrOwnerMemo::from_bytes(bytes))
                .collect(),
        },
        proof,
        folding_instance,
    };

    let mut hasher = Sha3_512::new();
    hasher.update(hash);
    hasher.update(&bincode::serialize(&note.body).map_err(|_| Error::ParseDataFailed)?);

    let params = VerifierParams::get_abar_to_abar(inputs_len, outputs_len, address_format)
        .map_err(|_| Error::FailedToLoadVerifierParams)?;

    verify_anon_xfr_note(&params, &note, &merkle_root, hasher)
        .map_err(|_| Error::ProofVerificationFailed)
}
