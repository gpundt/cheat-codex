#[cfg(target_os = "linux")]
mod linux;
#[cfg(target_os = "windows")]
mod windows;

#[cfg(target_os = "linux")]
pub use linux::MemoryError;
#[cfg(target_os = "windows")]
pub use windows::MemoryError;

#[cfg(target_os = "linux")]
pub fn find_gba_ram_base(pid: u32) -> Result<usize, MemoryError> {
    linux::find_gba_ram_base(pid)
}
#[cfg(target_os = "windows")]
pub fn find_gba_ram_base(pid: u32) -> Result<usize, MemoryError> {
    windows::find_gba_ram_base(pid)
}

#[cfg(target_os = "linux")]
pub fn read_u32(pid: u32, address: usize) -> Result<u32, MemoryError> {
    linux::read_u32(pid, address)
}
#[cfg(target_os = "windows")]
pub fn read_u32(pid: u32, address: usize) -> Result<u32, MemoryError> {
    windows::read_u32(pid, address)
}

#[cfg(target_os = "linux")]
pub fn write_u32(pid: u32, address: usize, value: u32) -> Result<(), MemoryError> {
    linux::write_u32(pid, address, value)
}
#[cfg(target_os = "windows")]
pub fn write_u32(pid: u32, address: usize, value: u32) -> Result<(), MemoryError> {
    windows::write_u32(pid, address, value)
}