#[macro_use]
extern crate maplit;

use clap::{App, SubCommand};

use parse::YAMLConfig;

mod parse;
mod preprocessor;

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
        match parse::Config::from_yaml(config) {
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
