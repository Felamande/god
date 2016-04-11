package sys

import (
	// _ "runtime/cgo"

	"syscall"
	"unsafe"
)

// libkernel32 *syscall.DLL
var (
	libuser32 *syscall.LazyDLL
)

var (
	registerHotKey *syscall.LazyProc
	getMessage     *syscall.LazyProc
	sendMessage    *syscall.LazyProc
)

// var getVolumeInformation  *syscall.Proc

func init() {
	// libkernel32 = syscall.MustLoadDLL("kernel32")
	libuser32 = syscall.NewLazyDLL("user32")

	// setConsoleCtrlHandler = libkernel32.MustFindProc("SetConsoleCtrlHandler")
	registerHotKey = libuser32.NewProc("RegisterHotKey")
	getMessage = libuser32.NewProc("GetMessageW")
	sendMessage = libuser32.NewProc("SendMessageW")

}

const (
	MaxPath = uint32(261)
)

func RegisterHotKey(hwnd IHwnd, id int, fsModifiers Modifier, vk Key) (bool, error) {
	ret, _, lerr := registerHotKey.Call(
		hwnd.Hwnd(),
		uintptr(id),
		uintptr(fsModifiers),
		uintptr(vk),
	)
	return ret != 0, lerr
}

func SendMessage(hwnd IHwnd, msg uint32, wParam, lParam uintptr) uintptr {
	ret, _, _ := sendMessage.Call(
		hwnd.Hwnd(),
		uintptr(msg),
		wParam,
		lParam,
	)

	return ret
}

func GetMessage(msg *MSG, hWnd IHwnd, msgFilterMin, msgFilterMax uint32) bool {
	ret, _, _ := getMessage.Call(
		uintptr(unsafe.Pointer(msg)),
		hWnd.Hwnd(),
		uintptr(msgFilterMin),
		uintptr(msgFilterMax),
	)

	return ret != 0
}
