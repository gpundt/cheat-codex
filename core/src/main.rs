use clap::{Parser, ValueEnum};
use log::{LevelFilter, debug, error, info, trace, warn};

mod memory;

#[cfg(target_os = "windows")]
use proc_mem::{Process, ProcMemError};

#[derive(Parser, Debug)]
#[command(name = "cheat-codex-core")]
#[command(version = "1.0")]
#[command(about = "Read/Write memory to ROM Emulators", long_about = None)]
struct Args {
    /// Action to take (read/write/get-base-address)
    #[arg(short, long, value_enum, required = true)]
    action: Action,

    /// Address to target
    #[arg(
        long,
        required_if_eq("action", "read"),
        required_if_eq("action", "write"),
    )]
    address: Option<String>,

    /// Value to overwrite address with
    #[arg(short, long, required_if_eq("action", "write"))]
    value: Option<u32>,

    /// PID of process to attach to
    #[arg(long, required_if_eq("action", "get-base-address"))]
    pid: Option<u32>,

    /// Turn on debug strings
    #[arg(short, long, default_value_t = false)]
    verbose: bool,
}

#[derive(ValueEnum, Clone, Debug)]
enum Action {
    Read,
    Write,
    GetBaseAddress,
}

fn main() {
    let args: Args = Args::parse();

    let log_level: LevelFilter = if args.verbose {
    LevelFilter::Debug
    } else {
    LevelFilter::Info
    };
    colog::default_builder()
        .filter_level(log_level)
        .init();

    match args.action {
        Action::GetBaseAddress => {
            let pid = args.pid.expect("--pid is required for --action get-base-address");
            match memory::get_base_address(pid) {
                Ok(base) => println!("{}", base),
                Err(e) => {
                    error!("Error getting base address: {:?}", e);
                    std::process::exit(1)
                }
            }
        },
        Action::Read => {
            match args.address {
                Some(address) => {
                    match read_value_at_address(&address) {
                        Ok(value) => println!("{}", value),
                        Err(e) => {
                            error!("Error reading value at {}: {:?}", address, e);
                            std::process::exit(1)
                        }
                    }
                }
                None => {
                    error!("Must provide --address for --action read");
                    std::process::exit(1);
                }
            }
        }
        Action::Write => {
            match (args.value, args.address) {
                (Some(value), Some(address)) => {
                    match write_value_to_address(value, &address) {
                        Ok(_) => println!("ok"),
                        Err(e) => {
                            error!("Error writing {} to {}: {:?}", value, address, e);
                            std::process::exit(1)
                        }
                    }
                }
                _ => {
                    error!("Must provide --address and --value for --action write");
                    std::process::exit(1);
                }
            }
        }
    }
}


fn read_value_at_address(address: &String) -> Result<String, String> {
    Ok(format!("0x0"))
}

fn write_value_to_address(value: u32, address: &String) -> Result<(), String> {
    Ok(())
}