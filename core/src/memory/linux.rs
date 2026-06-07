use std::fs;
use std::num::ParseIntError;

#[derive(Debug)]
pub enum MemoryError {
    IoError(std::io::Error),
    ParseError(String),
}

impl From<std::io::Error> for MemoryError {
    fn from(e: std::io::Error) -> Self {
        MemoryError::IoError(e)
    }
}

pub fn get_base_address(pid: u32) -> Result<usize, MemoryError> {
    let maps_path = format!("/proc/{}/maps", pid);
    let contents = fs::read_to_string(&maps_path)?;

    // Each line of /proc/{}/maps looks like:
    // 7f4a3b000000-7f4a3b001000 r--p 00000000 fd:01 123456  /path/to/binary
    // We want the start address of the very first line
    let first_line = contents
        .lines()
        .next()
        .ok_or_else(|| MemoryError::ParseError("maps file is empty".to_string()))?;

    let address_range = first_line
        .split_whitespace()
        .next()
        .ok_or_else(|| MemoryError::ParseError("could not parse first field".to_string()))?;

    let start_address = address_range
        .split('-')
        .next()
        .ok_or_else(|| MemoryError::ParseError("could not parse address range".to_string()))?;

    usize::from_str_radix(start_address, 16)
        .map_err(|e| MemoryError::ParseError(format!("could not parse hex address: {}", e)))
}