pub enum Error {
    WrongSelectorLength,
    UnknownSelector,
    ParseDataFailed,
    ProofVerificationFailed,
    ProofDecodeFailed,
    FailedToLoadVerifierParams,
    FoldingDecodeFailed,
    WrongLengthOfArguments,
    UnsupportInputsOutputs,
    InputOutOfBound,
    SerializeError,
    DeserializeError,
    DecodeError,
    ExecError,
}

impl Error {
    pub fn code(&self) -> u8 {
        match self {
            Error::WrongSelectorLength => 1,
            Error::UnknownSelector => 2,
            Error::ParseDataFailed => 3,
            Error::ProofVerificationFailed => 4,
            Error::ProofDecodeFailed => 5,
            Error::FailedToLoadVerifierParams => 6,
            Error::FoldingDecodeFailed => 7,
            Error::WrongLengthOfArguments => 8,
            Error::UnsupportInputsOutputs => 9,
            Error::InputOutOfBound => 10,
            Error::SerializeError => 11,
            Error::DeserializeError => 12,
            Error::DecodeError => 13,
            Error::ExecError => 14,
        }
    }
}

pub type Result<T> = std::result::Result<T, Error>;
