package ui

import (
	"strconv"
)

// A Win is an open window.
type Win interface {
	// Present swaps the front and back buffers, displaying that which
	// was drawn since the last call to Present.
	Present()

	// Events returns a channel of the window's events.
	Events() <-chan interface{}

	// Close closes the window.
	Close()
}

// A Key identifies a key on the keyboard.
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
	if (k >= 'a' && k <= 'z') || (k >= 'A' && k <= 'Z') || (k >= '0' && k <= '9') {
		return string([]rune{rune(k)})
	}
	if n, ok := keyNames[k]; ok {
		return n
	}
	return "Key(" + strconv.Itoa(int(k)) + ")"
}

// A Button identifies a button on a mouse.
type Button int

const (
	// ButtonLeft is the left mouse button.
	ButtonLeft Button = iota + 1
	// ButtonRight is the right mouse button.
	ButtonRight
	//ButtonCenter is the center mouse button.
	ButtonCenter
)

var buttonNames = map[Button]string{
	ButtonLeft:   "ButtonLeft",
	ButtonRight:  "ButtonRight",
	ButtonCenter: "ButtonCenter",
}

// String returns the human-readable representation of a Button.
func (m Button) String() string {
	if n, ok := buttonNames[m]; ok {
		return n
	}
	return "Button(" + strconv.Itoa(int(m)) + ")"
}

// KeyDown is an event signaling that a key on the keyboard was pressed.
type KeyDown struct {
	Key Key
}

// KeyUp is an event signaling that a key on the keyboard was released.
type KeyUp struct {
	Key Key
}

// MouseDown is an event signaling that a mouse button was pressed.
type MouseDown struct {
	Button Button
	// X and Y give the location of the pointer in the window.
	// 0, 0 is the upper left corner of the window, and
	// positive Y is downward.
	X, Y int
}

// MouseUp is an event signaling that a mouse button was released.
type MouseUp struct {
	Button Button
	// X and Y give the location of the pointer in the window.
	// 0, 0 is the upper left corner of the window, and
	// positive Y is downward.
	X, Y int
}

// MouseMove is an event signaling that the mouse has moved.
type MouseMove struct {
	// X and Y give the location of the pointer in the window.
	// 0, 0 is the upper left corner of the window, and
	// positive Y is downward.
	X, Y int
}

// WinClose is an event signaling that the "x" was clicked to close a window.
type WinClose struct{}
