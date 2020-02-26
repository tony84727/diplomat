use serde::{Deserialize, Serialize};

#[derive(Serialize, Deserialize, Debug, PartialEq)]
pub enum ChinesePreprocessorMode {
    #[serde(alias = "t2s")]
    T2S,
    #[serde(alias = "s2t")]
    S2T,
}

#[derive(Serialize, Deserialize, Debug, PartialEq)]
pub struct ChinesePreprocessorDesc {
    pub from: String,
    pub mode: ChinesePreprocessorMode,
    pub to: String,
}

#[derive(Serialize, Deserialize, Debug, PartialEq)]
pub struct CopyPreprocessorDesc {
    pub from: String,
    pub to: String,
}
