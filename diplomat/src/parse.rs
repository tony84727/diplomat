use std::{cmp::PartialEq, collections::HashMap};

use serde::de::DeserializeOwned;
use serde::{Deserialize, Serialize};

use crate::preprocessor::{ChinesePreprocessorDesc, CopyPreprocessorDesc};

#[derive(Serialize, Deserialize, Debug, PartialEq)]
struct TemplateOption {
    filename: String,
}

#[derive(Serialize, Deserialize, Debug, PartialEq)]
struct TemplateConfig {
    #[serde(rename = "type")]
    output_type: String,
    options: TemplateOption,
}

#[derive(Serialize, Deserialize, Debug, PartialEq)]
struct Output {
    selectors: Vec<String>,
    templates: Vec<TemplateConfig>,
}

#[derive(Serialize, Deserialize, Debug, PartialEq)]
#[serde(tag = "type", content = "options")]
enum PreprocessorConfig {
    #[serde(rename = "chinese")]
    Chinese(ChinesePreprocessorDesc),
    #[serde(rename = "copy")]
    Copy(CopyPreprocessorDesc),
}

#[derive(Serialize, Deserialize, Debug, PartialEq)]
pub struct Config {
    preprocessors: Vec<PreprocessorConfig>,
    outputs: Vec<Output>,
}

#[derive(Debug)]
pub enum YAMLParsingError {
    IO(std::io::Error),
    Parse(serde_yaml::Error),
}

pub trait YAMLConfig: Sized {
    fn from_yaml(path: &str) -> Result<Self, YAMLParsingError>;
}

impl<T> YAMLConfig for T
where
    T: DeserializeOwned,
{
    fn from_yaml(path: &str) -> Result<T, YAMLParsingError> {
        let source = std::fs::File::open(path).map_err(|x| YAMLParsingError::IO(x))?;
        serde_yaml::from_reader(&source).map_err(|x| YAMLParsingError::Parse(x))
    }
}

#[derive(Debug, PartialEq, Serialize, Deserialize)]
#[serde(untagged)]
enum YAMLTranslation {
    Nested(HashMap<String, YAMLTranslation>),
    Leaf(String),
}

#[cfg(test)]
mod tests {
    use crate::preprocessor::ChinesePreprocessorMode;

    use super::*;

    #[test]
    fn test_parse_config() {
        let expected = Config {
            preprocessors: vec![
                PreprocessorConfig::Chinese(ChinesePreprocessorDesc {
                    from: String::from("zh-TW"),
                    mode: ChinesePreprocessorMode::T2S,
                    to: String::from("zh-CN"),
                }),
                PreprocessorConfig::Copy(CopyPreprocessorDesc {
                    from: String::from("en"),
                    to: String::from("fr"),
                }),
            ],
            outputs: vec![Output {
                selectors: vec!["admin".to_string(), "manage".to_string()],
                templates: vec![TemplateConfig {
                    output_type: "js-object".to_string(),
                    options: TemplateOption {
                        filename: "control-panel.js".to_string(),
                    },
                }],
            }],
        };
        let config = Config::from_yaml("../testdata/diplomat.yaml").unwrap();
        assert_eq!(expected, config);
    }

    #[test]
    fn test_parse_yaml_translation() {
        let expected = YAMLTranslation::Nested(hashmap! {
            "admin".to_string() => YAMLTranslation::Nested(hashmap! {
                "admin".to_string() => YAMLTranslation::Nested(hashmap! {
                    "zh-TW".to_string() => YAMLTranslation::Leaf("管理員".to_string()),
                    "en".to_string() => YAMLTranslation::Leaf("Admin".to_string()),
                }),
                "message".to_string() => YAMLTranslation::Nested(
                    hashmap! {
                        "hello".to_string() => YAMLTranslation::Nested(
                            hashmap! {
                                "zh-TW".to_string() => YAMLTranslation::Leaf("您好".to_string()),
                                "en".to_string() => YAMLTranslation::Leaf("Hello!".to_string()),
                            }
                        )
                    }
                ),
            }),
        });
        let actual = YAMLTranslation::from_yaml("../testdata/admin.yaml").unwrap();
        assert_eq!(expected, actual);
    }
}
