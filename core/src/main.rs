use clap::{Parser, ValueEnum};
use log::{LevelFilter, debug, error, info, trace, warn};

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
    #[arg(
        short,
        long,
        required_if_eq("action", "write"),
    )]
    value: Option<u32>,

    /// Substring of process to attach to
    #[arg(
        short,
        long,
    )]
    process: Option<String>,

    /// Turn on debug strings
    #[arg(
        short,
        long,
        default_value_t = false,
    )]
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
        Action::GetBaseAddress => info!(
            "Getting base address of process: {}",
            args.process.unwrap(),
        ),
        Action::Read => info!(
            "Reading value at address: {}",
            args.address.unwrap(),
        ),
        Action::Write => info!(
            "Writing value {} to address {}",
            args.value.unwrap(),
            args.address.unwrap(),
        ),
    }
}


fn get_base_address(process_name: &str) {

}