//+build windows
package tui

import (
	"os"
	"syscall"
	"unsafe"
)

// SetConsoleTitle sets the console title
func SetConsoleTitle(title string) error {
	kernel32 := syscall.MustLoadDLL("Kernel32.dll")
	defer kernel32.Release()

	const ATTACH_PARENT_PROCESS = uintptr(^uint32(0))
	kernel32.MustFindProc("AttachConsole").Call(ATTACH_PARENT_PROCESS)
	hout, _ := syscall.GetStdHandle(syscall.STD_OUTPUT_HANDLE)
	hin, _ := syscall.GetStdHandle(syscall.STD_INPUT_HANDLE)
	os.Stdout = os.NewFile(uintptr(hout), "/dev/stdout")
	os.Stdin = os.NewFile(uintptr(hin), "/dev/stdin")

	setConsoleTitle := kernel32.MustFindProc("SetConsoleTitleW")
	_, _, err := setConsoleTitle.Call(uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(title))))
	return err
}
