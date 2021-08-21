//+build windows
package tui

import (
	"syscall"
	"unsafe"
)

// SetConsoleTitle sets the console title
func SetConsoleTitle(title string) error {
	kernel32 := syscall.MustLoadDLL("Kernel32.dll")
	defer kernel32.Release()

	setConsoleTitle := kernel32.MustFindProc("SetConsoleTitleW")
	_, _, err := setConsoleTitle.Call(uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(title))))
	return err
}
