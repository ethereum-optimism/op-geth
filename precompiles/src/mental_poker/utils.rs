use {
    crate::{Error, Result},
    ark_serialize::{CanonicalDeserialize, CanonicalSerialize},
};

pub(crate) fn deserialize<T: CanonicalDeserialize>(data: &[u8]) -> Result<T> {
    CanonicalDeserialize::deserialize_compressed(data).map_err(|_e| Error::DeserializeError)
}

pub(crate) fn serialize<T: CanonicalSerialize>(data: &T) -> Result<Vec<u8>> {
    let mut res = Vec::with_capacity(data.compressed_size());
    data.serialize_compressed(&mut res)
        .map(|_v| res)
        .map_err(|_e| Error::SerializeError)
}
