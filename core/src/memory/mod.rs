#[cfg(target_os = "linux")]
mod linux;

#[cfg(target_os = "windows")]
mod windows;

#[cfg(target_os = "linux")]
pub use linux::MemoryError;

#[cfg(target_os = "windows")]
pub use windows::MemoryError;

#[cfg(target_os = "linux")]
pub fn get_base_address(pid: u32) -> Result<usize, MemoryError> {
    linux::get_base_address(pid)
}

#[cfg(target_os = "windows")]
pub fn get_base_address(pid: u32) -> Result<usize, MemoryError> {
    windows::get_base_address(pid)
}