use clap::{App, SubCommand};
use serde::{Deserialize, Serialize};

#[derive(Serialize, Deserialize)]
enum ChinesePreprocessorMode {
    T2S,
    S2T,
}

#[derive(Serialize, Deserialize)]
struct ChinesePreprocessorOption {
    from: String,
    mode: ChinesePreprocessorMode,
    to: String,
}

#[derive(Serialize, Deserialize)]
struct ChinesePreprocessorConfig {
    #[serde(rename = "type")]
    preprocessor_type: String,
    options: ChinesePreprocessorMode,
}

#[derive(Serialize, Deserialize)]
struct TemplateOption {
    filename: String,
}

#[derive(Serialize, Deserialize)]
struct TemplateConfig {
    #[serde(rename = "type")]
    output_type: String,
    options: TemplateOption,
}

#[derive(Serialize, Deserialize)]
struct Output {
    selectors: Vec<String>,
    templates: Vec<TemplateConfig>,
}

#[derive(Serialize, Deserialize)]
struct Config {
    outputs: Vec<Output>,
}

enum ConfigParsingError {
    IO(std::io::Error),
    Parse(serde_yaml::Error),
}

impl Config {
    fn from_file(path: &str) -> Result<Self, ConfigParsingError> {
        let source = std::fs::File::open(path).map_err(|x| ConfigParsingError::IO(x))?;
        serde_yaml::from_reader(&source).map_err(|x| ConfigParsingError::Parse(x))
    }
}

fn main() {
    let matches = App::new("diplomat")
        .version("1.0")
        .author("Tony Duan <tony84727@gmail.com>")
        .subcommand(
            SubCommand::with_name("configtest")
                .usage("test the format of the config file")
                .arg_from_usage("-c, --config=[FILE] 'path to config file diplomat.yml'"),
        )
        .get_matches();

    if let Some(matches) = matches.subcommand_matches("configtest") {
        let config = matches.value_of("config").unwrap_or("diplomat.yml");
        match Config::from_file(config) {
            Err(_err) => println!("error when reading/parsing config"),
            Ok(config) => {
                println!("ok!");
                match serde_json::to_string_pretty(&config) {
                    Ok(pretty_print) => println!("Config: {}", pretty_print),
                    Err(err) => {
                        println!("Warning: error when trying to pretty print json: {}", err)
                    }
                };
            }
        }
    }
}
