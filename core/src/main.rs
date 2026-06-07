use clap::{Parser, ValueEnum};
use log::{LevelFilter, error};

mod memory;

#[derive(Parser, Debug)]
#[command(name = "cheat-codex-core")]
#[command(version = "1.0")]
#[command(about = "Read/Write memory to ROM Emulators", long_about = None)]
struct Args {
    #[arg(short, long, value_enum, required = true)]
    action: Action,

    #[arg(long, required_if_eq("action", "find-ram-base"),
                required_if_eq("action", "read"),
                required_if_eq("action", "write"))]
    pid: Option<u32>,

    /// Absolute address to read/write (ram_base + offset from yaml)
    #[arg(long, required_if_eq("action", "read"),
               required_if_eq("action", "write"))]
    address: Option<String>,

    #[arg(short, long, required_if_eq("action", "write"))]
    value: Option<u32>,

    #[arg(short, long, default_value_t = false)]
    verbose: bool,
}

#[derive(ValueEnum, Clone, Debug)]
enum Action {
    FindRamBase,
    Read,
    Write,
}

fn main() {
    let args: Args = Args::parse();

    let log_level = if args.verbose { LevelFilter::Debug } else { LevelFilter::Info };
    colog::default_builder()
        .filter_level(log_level)
        .target(env_logger::Target::Stderr)
        .init();

    match args.action {
        Action::FindRamBase => {
            let pid = args.pid.unwrap();
            match memory::find_gba_ram_base(pid) {
                Ok(base) => println!("0x{:x}", base),
                Err(e) => { error!("{:?}", e); std::process::exit(1); }
            }
        },

        Action::Read => {
            let pid = args.pid.unwrap();
            let address = parse_hex(&args.address.unwrap());
            match memory::read_u32(pid, address) {
                Ok(val) => println!("{}", val),
                Err(e) => { error!("{:?}", e); std::process::exit(1); }
            }
        },

        Action::Write => {
            let pid = args.pid.unwrap();
            let address = parse_hex(&args.address.unwrap());
            let value = args.value.unwrap();
            match memory::write_u32(pid, address, value) {
                Ok(_) => println!("ok"),
                Err(e) => { error!("{:?}", e); std::process::exit(1); }
            }
        },
    }
}

fn parse_hex(s: &str) -> usize {
    let hex = s.trim_start_matches("0x").trim_start_matches("0X");
    usize::from_str_radix(hex, 16).unwrap_or_else(|_| {
        eprintln!("Invalid hex address: {}", s);
        std::process::exit(1);
    })
}