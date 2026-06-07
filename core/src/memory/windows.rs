use std::ffi::OsStr;
use std::os::windows::ffi::OsStrExt;

use winapi::shared::minwindef::{DWORD, HMODULE};
use winapi::um::handleapi::CloseHandle;
use winapi::um::processthreadsapi::OpenProcess;
use winapi::um::psapi::EnumProcessModules;
use winapi::um::psapi::GetModuleFileNameExW;
use winapi::um::winnt::{PROCESS_QUERY_INFORMATION, PROCESS_VM_READ};

#[derive(Debug)]
pub enum MemoryError {
    IoError(std::io::Error),
    ParseError(String),
    WindowsError(u32),
}

pub fn get_base_address(pid: u32) -> Result<usize, MemoryError> {
    unsafe {
        let handle = OpenProcess(
            PROCESS_QUERY_INFORMATION | PROCESS_VM_READ,
            0,
            pid,
        );

        if handle.is_null() {
            return Err(MemoryError::WindowsError(
                winapi::um::errhandlingapi::GetLastError()
            ));
        }

        // Get the first module (the executable itself)
        let mut module: HMODULE = std::ptr::null_mut();
        let mut cb_needed: DWORD = 0;

        let result =  EnumProcessModules(
            handle,
            &mut module,
            std::mem::size_of::<HMODULE>() as DWORD,
            &mut cb_needed,
        );

        CloseHandle(handle);

        if result == 0 {
            return Err(MemoryError::WindowsError(
                winapi::um::errhandlingapi::GetLastError()
            ));
        }

        // The HMODULE of the first module is the base address
        Ok(module as usize)
    }
}