package sys

type MSG struct {
	HWnd    HWND
	Message uint32
	WParam  uintptr
	LParam  uintptr
	Time    uint32
	Pt      struct { //POINT
		x, y int32
	}
}

type IHwnd interface {
	Hwnd() uintptr
}

type HWND uintptr
type BOOL uint32


func (h HWND) Hwnd() uintptr {
	return uintptr(h)
}

const (
	NULL  uintptr = 0
	FALSE BOOL    = iota
	TRUE
)

func (b BOOL) Bool() bool {
	return b > 0
}

func (b BOOL) Hwnd() uintptr {
	return uintptr(b)
}

const (
	MOD_ALT     Modifier = 1
	MOD_CONTROL Modifier = 2
	MOD_SHIFT   Modifier = 4
	MOD_WIN     Modifier = 8
)

const (
	VK_LBUTTON             Key = 0x01
	VK_RBUTTON             Key = 0x02
	VK_CANCEL              Key = 0x03
	VK_MBUTTON             Key = 0x04
	VK_XBUTTON1            Key = 0x05
	VK_XBUTTON2            Key = 0x06
	VK_BACK                Key = 0x08
	VK_TAB                 Key = 0x09
	VK_CLEAR               Key = 0x0C
	VK_RETURN              Key = 0x0D
	VK_SHIFT               Key = 0x10
	VK_CONTROL             Key = 0x11
	VK_MENU                Key = 0x12
	VK_PAUSE               Key = 0x13
	VK_CAPITAL             Key = 0x14
	VK_KANA                Key = 0x15
	VK_HANGUEL             Key = 0x15
	VK_HANGUL              Key = 0x15
	VK_JUNJA               Key = 0x17
	VK_FINAL               Key = 0x18
	VK_HANJA               Key = 0x19
	VK_KANJI               Key = 0x19
	VK_ESCAPE              Key = 0x1B
	VK_CONVERT             Key = 0x1C
	VK_NONCONVERT          Key = 0x1D
	VK_ACCEPT              Key = 0x1E
	VK_MODECHANGE          Key = 0x1F
	VK_SPACE               Key = 0x20
	VK_PRIOR               Key = 0x21
	VK_NEXT                Key = 0x22
	VK_END                 Key = 0x23
	VK_HOME                Key = 0x24
	VK_LEFT                Key = 0x25
	VK_UP                  Key = 0x26
	VK_RIGHT               Key = 0x27
	VK_DOWN                Key = 0x28
	VK_SELECT              Key = 0x29
	VK_PRINT               Key = 0x2A
	VK_EXECUTE             Key = 0x2B
	VK_SNAPSHOT            Key = 0x2C
	VK_INSERT              Key = 0x2D
	VK_DELETE              Key = 0x2E
	VK_HELP                Key = 0x2F
	VK_0                   Key = 0x30
	VK_1                   Key = 0x31
	VK_2                   Key = 0x32
	VK_3                   Key = 0x33
	VK_4                   Key = 0x34
	VK_5                   Key = 0x35
	VK_6                   Key = 0x36
	VK_7                   Key = 0x37
	VK_8                   Key = 0x38
	VK_9                   Key = 0x39
	VK_A                   Key = 0x41
	VK_B                   Key = 0x42
	VK_C                   Key = 0x43
	VK_D                   Key = 0x44
	VK_E                   Key = 0x45
	VK_F                   Key = 0x46
	VK_G                   Key = 0x47
	VK_H                   Key = 0x48
	VK_I                   Key = 0x49
	VK_J                   Key = 0x4A
	VK_K                   Key = 0x4B
	VK_L                   Key = 0x4C
	VK_M                   Key = 0x4D
	VK_N                   Key = 0x4E
	VK_O                   Key = 0x4F
	VK_P                   Key = 0x50
	VK_Q                   Key = 0x51
	VK_R                   Key = 0x52
	VK_S                   Key = 0x53
	VK_T                   Key = 0x54
	VK_U                   Key = 0x55
	VK_V                   Key = 0x56
	VK_W                   Key = 0x57
	VK_X                   Key = 0x58
	VK_Y                   Key = 0x59
	VK_Z                   Key = 0x5A
	VK_LWIN                Key = 0x5B
	VK_RWIN                Key = 0x5C
	VK_APPS                Key = 0x5D
	VK_SLEEP               Key = 0x5F
	VK_NUMPAD0             Key = 0x60
	VK_NUMPAD1             Key = 0x61
	VK_NUMPAD2             Key = 0x62
	VK_NUMPAD3             Key = 0x63
	VK_NUMPAD4             Key = 0x64
	VK_NUMPAD5             Key = 0x65
	VK_NUMPAD6             Key = 0x66
	VK_NUMPAD7             Key = 0x67
	VK_NUMPAD8             Key = 0x68
	VK_NUMPAD9             Key = 0x69
	VK_MULTIPLY            Key = 0x6A
	VK_ADD                 Key = 0x6B
	VK_SEPARATOR           Key = 0x6C
	VK_SUBTRACT            Key = 0x6D
	VK_DECIMAL             Key = 0x6E
	VK_DIVIDE              Key = 0x6F
	VK_F1                  Key = 0x70
	VK_F2                  Key = 0x71
	VK_F3                  Key = 0x72
	VK_F4                  Key = 0x73
	VK_F5                  Key = 0x74
	VK_F6                  Key = 0x75
	VK_F7                  Key = 0x76
	VK_F8                  Key = 0x77
	VK_F9                  Key = 0x78
	VK_F10                 Key = 0x79
	VK_F11                 Key = 0x7A
	VK_F12                 Key = 0x7B
	VK_F13                 Key = 0x7C
	VK_F14                 Key = 0x7D
	VK_F15                 Key = 0x7E
	VK_F16                 Key = 0x7F
	VK_F17                 Key = 0x80
	VK_F18                 Key = 0x81
	VK_F19                 Key = 0x82
	VK_F20                 Key = 0x83
	VK_F21                 Key = 0x84
	VK_F22                 Key = 0x85
	VK_F23                 Key = 0x86
	VK_F24                 Key = 0x87
	VK_NUMLOCK             Key = 0x90
	VK_SCROLL              Key = 0x91
	VK_LSHIFT              Key = 0xA0
	VK_RSHIFT              Key = 0xA1
	VK_LCONTROL            Key = 0xA2
	VK_RCONTROL            Key = 0xA3
	VK_LMENU               Key = 0xA4
	VK_RMENU               Key = 0xA5
	VK_BROWSER_BACK        Key = 0xA6
	VK_BROWSER_FORWARD     Key = 0xA7
	VK_BROWSER_REFRESH     Key = 0xA8
	VK_BROWSER_STOP        Key = 0xA9
	VK_BROWSER_SEARCH      Key = 0xAA
	VK_BROWSER_FAVORITES   Key = 0xAB
	VK_BROWSER_HOME        Key = 0xAC
	VK_VOLUME_MUTE         Key = 0xAD
	VK_VOLUME_DOWN         Key = 0xAE
	VK_VOLUME_UP           Key = 0xAF
	VK_MEDIA_NEXT_TRACK    Key = 0xB0
	VK_MEDIA_PREV_TRACK    Key = 0xB1
	VK_MEDIA_STOP          Key = 0xB2
	VK_MEDIA_PLAY_PAUSE    Key = 0xB3
	VK_LAUNCH_MAIL         Key = 0xB4
	VK_LAUNCH_MEDIA_SELECT Key = 0xB5
	VK_LAUNCH_APP1         Key = 0xB6
	VK_LAUNCH_APP2         Key = 0xB7
	VK_OEM_1               Key = 0xBA
	VK_OEM_PLUS            Key = 0xBB
	VK_OEM_COMMA           Key = 0xBC
	VK_OEM_MINUS           Key = 0xBD
	VK_OEM_PERIOD          Key = 0xBE
	VK_OEM_2               Key = 0xBF
	VK_OEM_3               Key = 0xC0
	VK_OEM_4               Key = 0xDB
	VK_OEM_5               Key = 0xDC
	VK_OEM_6               Key = 0xDD
	VK_OEM_7               Key = 0xDE
	VK_OEM_8               Key = 0xDF
	VK_OEM_102             Key = 0xE2
	VK_PROCESSKEY          Key = 0xE5
	VK_PACKET              Key = 0xE7
	VK_ATTN                Key = 0xF6
	VK_CRSEL               Key = 0xF7
	VK_EXSEL               Key = 0xF8
	VK_EREOF               Key = 0xF9
	VK_PLAY                Key = 0xFA
	VK_ZOOM                Key = 0xFB
	VK_NONAME              Key = 0xFC
	VK_PA1                 Key = 0xFD
	VK_OEM_CLEAR           Key = 0xFE
)

func (k Key) Raw() uintptr {
	return uintptr(k)
}

func (m Modifier) Raw() uintptr {
	return uintptr(m)
}

const (
	WM_HOTKEY uint32 = 0x0312
)
