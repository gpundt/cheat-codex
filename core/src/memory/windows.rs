use winapi::shared::minwindef::DWORD;
use winapi::um::handleapi::CloseHandle;
use winapi::um::memoryapi::{ReadProcessMemory, VirtualQueryEx, WriteProcessMemory};
use winapi::um::processthreadsapi::OpenProcess;
use winapi::um::winnt::{
    HANDLE, MEMORY_BASIC_INFORMATION, MEM_COMMIT, PAGE_READWRITE,
    PROCESS_QUERY_INFORMATION, PROCESS_VM_OPERATION, PROCESS_VM_READ, PROCESS_VM_WRITE,
};

#[derive(Debug)]
pub enum MemoryError {
    WindowsError(u32),
    RegionNotFound,
    ReadError(String),
    WriteError(String),
}

const EWRAM_SIZE: usize = 0x40000; // 256KB - GBA external RAM

pub fn find_gba_ram_base(pid: u32) -> Result<usize, MemoryError> {
    unsafe {
        let handle = open_process(pid, PROCESS_QUERY_INFORMATION | PROCESS_VM_READ)?;

        let mut address: usize = 0;
        let mut mbi: MEMORY_BASIC_INFORMATION = std::mem::zeroed();

        loop {
            let result = VirtualQueryEx(
                handle,
                address as *const _,
                &mut mbi,
                std::mem::size_of::<MEMORY_BASIC_INFORMATION>(),
            );

            if result == 0 {
                break;
            }

            let size = mbi.RegionSize;
            let is_committed = mbi.State == MEM_COMMIT;
            let is_rw = mbi.Protect == PAGE_READWRITE;
            // Private memory (not mapped from a file) has Type == MEM_PRIVATE
            let is_private = mbi.Type == winapi::um::winnt::MEM_PRIVATE;

            if size == EWRAM_SIZE && is_committed && is_rw && is_private {
                CloseHandle(handle);
                return Ok(mbi.BaseAddress as usize);
            }

            // Advance to next region
            address = (mbi.BaseAddress as usize) + mbi.RegionSize;
        }

        CloseHandle(handle);
        Err(MemoryError::RegionNotFound)
    }
}

pub fn read_u32(pid: u32, address: usize) -> Result<u32, MemoryError> {
    unsafe {
        let handle = open_process(pid, PROCESS_QUERY_INFORMATION | PROCESS_VM_READ)?;
        let mut buffer: u32 = 0;
        let mut bytes_read = 0usize;

        let result = ReadProcessMemory(
            handle,
            address as *const _,
            &mut buffer as *mut u32 as *mut _,
            std::mem::size_of::<u32>(),
            &mut bytes_read,
        );
        CloseHandle(handle);

        if result == 0 {
            return Err(MemoryError::ReadError(format!(
                "ReadProcessMemory failed: {}",
                winapi::um::errhandlingapi::GetLastError()
            )));
        }
        Ok(buffer)
    }
}

pub fn write_u32(pid: u32, address: usize, value: u32) -> Result<(), MemoryError> {
    unsafe {
        let handle = open_process(
            pid,
            PROCESS_QUERY_INFORMATION | PROCESS_VM_READ | PROCESS_VM_WRITE | PROCESS_VM_OPERATION,
        )?;
        let mut bytes_written = 0usize;

        let result = WriteProcessMemory(
            handle,
            address as *mut _,
            &value as *const u32 as *const _,
            std::mem::size_of::<u32>(),
            &mut bytes_written,
        );
        CloseHandle(handle);

        if result == 0 {
            return Err(MemoryError::WriteError(format!(
                "WriteProcessMemory failed: {}",
                winapi::um::errhandlingapi::GetLastError()
            )));
        }
        Ok(())
    }
}

fn open_process(pid: u32, access: DWORD) -> Result<HANDLE, MemoryError> {
    unsafe {
        let handle = OpenProcess(access, 0, pid);
        if handle.is_null() {
            return Err(MemoryError::WindowsError(
                winapi::um::errhandlingapi::GetLastError()
            ));
        }
        Ok(handle)
    }
}