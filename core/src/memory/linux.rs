use std::fs;

#[derive(Debug)]
pub enum MemoryError {
    IoError(std::io::Error),
    ParseError(String),
    RegionNotFound,
    ReadError(String),
    WriteError(String),
}

impl From<std::io::Error> for MemoryError {
    fn from(e: std::io::Error) -> Self {
        MemoryError::IoError(e)
    }
}

const EWRAM_SIZE: usize = 0x40000; // 256KB - GBA external RAM

pub fn find_gba_ram_base(pid: u32) -> Result<usize, MemoryError> {
    let contents = fs::read_to_string(format!("/proc/{}/maps", pid))?;

    for line in contents.lines() {
        let mut fields = line.split_whitespace();

        let range = fields.next()
            .ok_or_else(|| MemoryError::ParseError("bad line".to_string()))?;

        let perms = fields.next()
            .ok_or_else(|| MemoryError::ParseError("no perms".to_string()))?;

        let mut range_parts = range.split('-');
        let start = parse_hex(range_parts.next().unwrap_or(""))?;
        let end   = parse_hex(range_parts.next().unwrap_or(""))?;
        let size  = end - start;

        // GBA EWRAM is exactly 256KB, readable and writable, anonymous (no file backing)
        // Anonymous regions have no filename at the end of the line
        let is_anonymous = line.split_whitespace().count() == 5;
        let is_rw = perms.starts_with("rw");

        if size == EWRAM_SIZE && is_rw && is_anonymous {
            return Ok(start);
        }
    }

    Err(MemoryError::RegionNotFound)
}

pub fn read_u32(pid: u32, address: usize) -> Result<u32, MemoryError> {
    use std::io::{Read, Seek, SeekFrom};

    let mut file = std::fs::OpenOptions::new()
        .read(true)
        .open(format!("/proc/{}/mem", pid))?;

    file.seek(SeekFrom::Start(address as u64))?;

    let mut buf = [0u8; 4];
    file.read_exact(&mut buf)
        .map_err(|e| MemoryError::ReadError(e.to_string()))?;

    Ok(u32::from_le_bytes(buf))
}

pub fn write_u32(pid: u32, address: usize, value: u32) -> Result<(), MemoryError> {
    use std::io::{Seek, SeekFrom, Write};

    let mut file = std::fs::OpenOptions::new()
        .write(true)
        .open(format!("/proc/{}/mem", pid))?;

    file.seek(SeekFrom::Start(address as u64))?;
    file.write_all(&value.to_le_bytes())
        .map_err(|e| MemoryError::WriteError(e.to_string()))?;

    Ok(())
}

fn parse_hex(s: &str) -> Result<usize, MemoryError> {
    usize::from_str_radix(s, 16)
        .map_err(|e| MemoryError::ParseError(e.to_string()))
}