//+build windows
package tui

import (
	"syscall"
	"unsafe"
)

// SetConsoleTitle sets the console title
func SetConsoleTitle(title string) error {
	handle, err := syscall.LoadLibrary("Kernel32.dll")
	if err != nil {
		return err
	}
	defer syscall.FreeLibrary(handle)
	proc, err := syscall.GetProcAddress(handle, "SetConsoleTitleW")
	if err != nil {
		return err
	}
	_, _, err = syscall.Syscall(proc, 1, uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(title))), 0, 0)
	return err
}
