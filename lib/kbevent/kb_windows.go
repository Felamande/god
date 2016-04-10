package kbevent

import (
	"errors"
	"fmt"
	"runtime"
	"strings"

	"github.com/Felamande/god/lib/sys"
)

func New() Interface {
	return &Keyboard{
		handlers: make(map[uintptr]func()),
		started:  false,
	}
}

type Keyboard struct {
	handlers map[uintptr]func()
	started  bool
}

func (k *Keyboard) Init() error {

	return nil
}

func (k *Keyboard) KeyCodeOf(key string) (uint8, bool) {
	code, exist := str2Key[key]
	return uint8(code), exist
}

func (k *Keyboard) ModifierCodeOf(m string) (uint8, bool) {
	code, exist := str2Mod[m]
	return uint8(code), exist
}

func (k *Keyboard) Parse(seq string) (mods []uint8, key uint8, err error) {

	seq = strings.ToLower(seq)
	keyStrs := strings.Split(seq, "+")
	unsupported := []string{}

	for _, ks := range keyStrs {
		ks = strings.TrimSpace(ks)

		if kcode, exist := k.KeyCodeOf(ks); exist {
			key = kcode
		} else if mcode, exist := k.ModifierCodeOf(ks); exist {
			mods = append(mods, mcode)
		} else {
			unsupported = append(unsupported, ks)
		}
	}
	if len(unsupported) != 0 {
		err = fmt.Errorf("unsupported keys or modifiers: %v", unsupported)
		return
	}
	return
}

func (k *Keyboard) GetSeq(mods []uint8, keyCode uint8) (string, error) {
	var seq string

	for _, code := range mods {
		if mod, exist := mod2Str[sys.Modifier(code)]; exist {
			seq += mod + "+"
		} else {
			return "", fmt.Errorf("unsupported modifier code: %v", code)
		}
	}

	if key, exist := key2Str[sys.Key(keyCode)]; exist {
		seq += key
	} else {
		return "", fmt.Errorf("unsupported key code: %v", keyCode)
	}

	return strings.TrimRight(seq, "+"), nil

}

func (k *Keyboard) Call(seq string) error {
	if f := k.HandlerOf(seq); f != nil {
		f()
	}
	return fmt.Errorf("no handlers for %s", seq)

}

func (k *Keyboard) HandlerOf(seq string) func() {
	mods, key, err := k.Parse(seq)
	if err != nil {
		return nil
	}

	lparamCode := uintptr(0)
	for _, mcode := range mods {
		lparamCode += uintptr(mcode)
	}
	lparamCode += uintptr(key) << 16

	return k.handlers[lparamCode]

}

func (k *Keyboard) Bind(seq string, f func()) error {
	mods, key, err := k.Parse(seq)

	if err != nil {
		return fmt.Errorf("bind %s: %v", seq, err)
	}
	if key == 0 {
		return fmt.Errorf("bind %s: %s", seq, "empty key")
	}

	lparamCode := uintptr(0)
	for _, mcode := range mods {
		lparamCode += uintptr(mcode)
	}

	if !sys.RegisterHotKey(sys.HWND(0), 1, sys.Modifier(lparamCode), sys.Key(key)) {
		return fmt.Errorf("bind %s: failed", seq)
	}

	lparamCode += uintptr(key) << 16

	k.handlers[lparamCode] = f
	// k.Start()
	return nil

}

func (k *Keyboard) ReadEvents(seqChan chan string, errChan chan error) error {
	if k.started {
		return errors.New("started twice")
	}

	// go func() {
	msg := new(sys.MSG)
	for sys.GetMessage(msg, sys.HWND(0), 0, 0) {
		var modCodes []uint8
		for i := uint(0); i < 4; i++ {
			modCode := uint8(msg.LParam & (1 << i))
			if modCode == 0 {
				continue
			}
			modCodes = append(modCodes, modCode)
		}
		keyCode := uint8(msg.LParam >> 16)
		seq, err := k.GetSeq(modCodes, keyCode)
		if err != nil {
			errChan <- err
		}
		seqChan <- seq

	}
	errChan <- errors.New("GetMessage terminated")
	// }()
	k.started = true
	return nil
}

func (k *Keyboard) Start() error {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	if k.started {
		return errors.New("started twice")
	}

	msg := new(sys.MSG)
	k.started = true

	for sys.GetMessage(msg, sys.HWND(0), 0, 0) {
		if msg.Message != sys.WM_HOTKEY {
			continue
		}
		// fmt.Printf("%v\n", key2Str[sys.Key(msg.LParam>>16)])
		if f, ok := k.handlers[msg.LParam]; ok {
			f()
		}
	}
	return nil

}

var str2Mod = map[string]sys.Modifier{
	"shift": sys.MOD_SHIFT,
	"alt":   sys.MOD_ALT,
	"ctrl":  sys.MOD_CONTROL,
	"win":   sys.MOD_WIN,
}
var str2Key = map[string]sys.Key{
	"0":   sys.VK_0,
	"1":   sys.VK_1,
	"2":   sys.VK_2,
	"3":   sys.VK_3,
	"4":   sys.VK_4,
	"5":   sys.VK_5,
	"6":   sys.VK_6,
	"7":   sys.VK_7,
	"8":   sys.VK_8,
	"9":   sys.VK_9,
	"a":   sys.VK_A,
	"b":   sys.VK_B,
	"c":   sys.VK_C,
	"d":   sys.VK_D,
	"e":   sys.VK_E,
	"f":   sys.VK_F,
	"g":   sys.VK_G,
	"h":   sys.VK_H,
	"i":   sys.VK_I,
	"j":   sys.VK_J,
	"k":   sys.VK_K,
	"l":   sys.VK_L,
	"m":   sys.VK_M,
	"n":   sys.VK_N,
	"o":   sys.VK_O,
	"p":   sys.VK_P,
	"q":   sys.VK_Q,
	"r":   sys.VK_R,
	"s":   sys.VK_S,
	"t":   sys.VK_T,
	"u":   sys.VK_U,
	"v":   sys.VK_V,
	"w":   sys.VK_W,
	"x":   sys.VK_X,
	"y":   sys.VK_Y,
	"z":   sys.VK_Z,
	"f1":  sys.VK_F1,
	"f2":  sys.VK_F2,
	"f3":  sys.VK_F3,
	"f4":  sys.VK_F4,
	"f5":  sys.VK_F5,
	"f6":  sys.VK_F6,
	"f7":  sys.VK_F7,
	"f8":  sys.VK_F8,
	"f9":  sys.VK_F9,
	"f10": sys.VK_F10,
	"f11": sys.VK_F11,
	"f12": sys.VK_F12,
}

var mod2Str = map[sys.Modifier]string{
	sys.MOD_SHIFT:   "shift",
	sys.MOD_ALT:     "alt",
	sys.MOD_CONTROL: "ctrl",
	sys.MOD_WIN:     "win",
}

var key2Str = map[sys.Key]string{
	sys.VK_0:   "0",
	sys.VK_1:   "1",
	sys.VK_2:   "2",
	sys.VK_3:   "3",
	sys.VK_4:   "4",
	sys.VK_5:   "5",
	sys.VK_6:   "6",
	sys.VK_7:   "7",
	sys.VK_8:   "8",
	sys.VK_9:   "9",
	sys.VK_A:   "a",
	sys.VK_B:   "b",
	sys.VK_C:   "c",
	sys.VK_D:   "d",
	sys.VK_E:   "e",
	sys.VK_F:   "f",
	sys.VK_G:   "g",
	sys.VK_H:   "h",
	sys.VK_I:   "i",
	sys.VK_J:   "j",
	sys.VK_K:   "k",
	sys.VK_L:   "l",
	sys.VK_M:   "m",
	sys.VK_N:   "n",
	sys.VK_O:   "o",
	sys.VK_P:   "p",
	sys.VK_Q:   "q",
	sys.VK_R:   "r",
	sys.VK_S:   "s",
	sys.VK_T:   "t",
	sys.VK_U:   "u",
	sys.VK_V:   "v",
	sys.VK_W:   "w",
	sys.VK_X:   "x",
	sys.VK_Y:   "y",
	sys.VK_Z:   "z",
	sys.VK_F1:  "f1",
	sys.VK_F2:  "f2",
	sys.VK_F3:  "f3",
	sys.VK_F4:  "f4",
	sys.VK_F5:  "f5",
	sys.VK_F6:  "f6",
	sys.VK_F7:  "f7",
	sys.VK_F8:  "f8",
	sys.VK_F9:  "f9",
	sys.VK_F10: "f10",
	sys.VK_F11: "f11",
	sys.VK_F12: "f12",
}
