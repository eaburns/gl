package ui

import (
	"strconv"
)

// A Win is an open window.
type Win interface {
	// Events returns a channel of the window's events.
	Events() <-chan interface{}

	// Close closes the window.
	Close()
}

// A Key is the identifier of a key on the keyboard.
// The letter and number keys are represented by their ASCII encoding.  All other
// keys are represented by constants.
type Key int

const (
	// KeyEnter is the Enter or return key.
	KeyEnter = iota + 256
	// KeySpace is the space bar key.
	KeySpace
	// KeyUpArrow is the up key on the directional pad.
	KeyUpArrow
	// KeyDownArrow is the down key on the directional pad.
	KeyDownArrow
	// KeyLeftArrow is the left key on the directional pad.
	KeyLeftArrow
	// KeyRightArrow is the right key on the directional pad.
	KeyRightArrow
	// KeyRightShift is the right shift key.
	KeyRightShift
	// KeyLeftShift is the left shift key.
	KeyLeftShift
	// KeyBackSpace is the backspace key.
	KeyBackSpace
	// KeyDelete is the delete key.
	KeyDelete
)

var keyNames = map[Key]string{
	KeyEnter:      "KeyEnter",
	KeySpace:      "KeySpace",
	KeyUpArrow:    "KeyUpArrow",
	KeyDownArrow:  "KeyDownArrow",
	KeyLeftArrow:  "KeyLeftArrow",
	KeyRightArrow: "KeyRightArrow",
	KeyRightShift: "KeyRightShift",
	KeyLeftShift:  "KeyLeftShift",
	KeyBackSpace:  "KeyBackSpace",
	KeyDelete:     "KeyDelete",
}

// String returns the human-readable representation of the Key.
func (k Key) String() string {
	if (k >= 'a' && k <= 'z') || (k >= 'A' && k <= 'Z') {
		return string([]rune{rune(k)})
	}
	if n, ok := keyNames[k]; ok {
		return n
	}
	return "Key(" + strconv.Itoa(int(k)) + ")"
}

// A KeyEvent signals a change in the state of a key on the keyboard.
type KeyEvent struct {
	// Down is true if they key is was pressed down, otherwise it is false.
	Down bool

	// Key is the key that was either pressed or depressed.
	Key Key
}

// A WinEventType identifies the type of event that happened on a window.
type WinEventType int

const (
	// WinClose is the WinEvent type when the "x" is clicked to close the window.
	WinClose WinEventType = iota
	// WinResize is the WinEvent type when the window size changes.
	WinResize
	// WinEnter is the WinEvent type when the window has gained mouse focus.
	WinEnter
	// WinLeave is the WinEvent type when the window has lost mouse focus.
	WinLeave
	// WinFocus is the WinEvent type when the window has gained keyboard focus.
	WinFocus
	// WinUnFocus is the WinEvent type when the window has lost keyboard focus.
	WinUnFocus
)

var winEventTypeNames = map[WinEventType]string{
	WinClose:   "WinClose",
	WinResize:  "WinResize",
	WinEnter:   "WinEnter",
	WinLeave:   "WinLeave",
	WinFocus:   "WinFocus",
	WinUnFocus: "WinUnFocus",
}

// String returns a human-readable representation of a WinEventType.
func (t WinEventType) String() string {
	if s, ok := winEventTypeNames[t]; ok {
		return s
	}
	return "WinEventType(" + strconv.Itoa(int(t)) + ")"
}

// A WinEvent signals a change to a window.
type WinEvent struct {
	// Type is the type of event.
	Type WinEventType

	// Width and Height give the size of the window.
	Width, Height int
}
