package sdl

/*
#include <SDL.h>

#cgo darwin CFLAGS: -I/Library/Frameworks/SDL2.framework/Headers
#cgo darwin LDFLAGS: -framework SDL2

#cgo linux pkg-config: sdl2

static Uint32 sdlEventType(SDL_Event *e){
	return e->type;
}
*/
import "C"

import (
	"errors"
	"time"
	"unsafe"

	"github.com/eaburns/gui/thread0"
	"github.com/eaburns/gui/ui"
)

// Wins is the map of all open windows.  It can only be accessed by the Init go routine.
var wins = make(map[uint32]*window, 1)

// Init initializes SDL and starts the event polling loop.  It must be called from
// the main go routine, and it only returns if there is an error.  Hz specifies
// the frequency (in milliseconds) that SDL polls events.
func Init(hz int) error {
	if C.SDL_Init(C.SDL_INIT_EVERYTHING) != 0 {
		return sdlError()
	}
	go func() {
		tick := time.NewTicker(time.Duration(hz) * time.Millisecond)
		for _ = range tick.C {
			thread0.Do(poll)
		}
	}()
	return nil
}

func poll() {
	for {
		var e C.SDL_Event
		if C.SDL_PollEvent(&e) == 0 {
			return
		}

		switch C.sdlEventType(&e) {
		case C.SDL_KEYDOWN:
			keyEvent(&e, true)
		case C.SDL_KEYUP:
			keyEvent(&e, false)
		case C.SDL_MOUSEBUTTONDOWN:
			mouseButtonEvent(&e, true)
		case C.SDL_MOUSEBUTTONUP:
			mouseButtonEvent(&e, false)
		case C.SDL_MOUSEMOTION:
			mouseMoveEvent(&e)
		case C.SDL_WINDOWEVENT:
			winEvent(&e)
		}
	}
}

var keys = map[int]ui.Key{
	C.SDLK_RETURN:    ui.KeyEnter,
	C.SDLK_SPACE:     ui.KeySpace,
	C.SDLK_UP:        ui.KeyUpArrow,
	C.SDLK_DOWN:      ui.KeyDownArrow,
	C.SDLK_LEFT:      ui.KeyLeftArrow,
	C.SDLK_RIGHT:     ui.KeyRightArrow,
	C.SDLK_RSHIFT:    ui.KeyRightShift,
	C.SDLK_LSHIFT:    ui.KeyLeftShift,
	C.SDLK_BACKSPACE: ui.KeyBackSpace,
	C.SDLK_DELETE:    ui.KeyDelete,
}

func mapKeySym(kcode C.SDL_Keycode) (ui.Key, bool) {
	if (kcode >= 'a' && kcode <= 'z') || (kcode >= 'A' && kcode <= 'Z') || (kcode >= '0' && kcode <= '9') {
		return ui.Key(kcode), true
	}
	key, ok := keys[int(kcode)]
	return ui.Key(key), ok
}

func keyEvent(e *C.SDL_Event, down bool) {
	k := (*C.SDL_KeyboardEvent)(unsafe.Pointer(e))
	if k.repeat != 0 {
		return
	}
	win, ok := wins[uint32(k.windowID)]
	if !ok {
		return
	}

	key, ok := mapKeySym(k.keysym.sym)
	if !ok {
		return
	}

	if down {
		win.events <- ui.KeyDown{Key: key}
	} else {
		win.events <- ui.KeyUp{Key: key}
	}
}

var mouseButtons = map[C.Uint8]ui.Button{
	C.SDL_BUTTON_LEFT:   ui.ButtonLeft,
	C.SDL_BUTTON_RIGHT:  ui.ButtonRight,
	C.SDL_BUTTON_MIDDLE: ui.ButtonCenter,
}

func mouseButtonEvent(e *C.SDL_Event, down bool) {
	m := (*C.SDL_MouseButtonEvent)(unsafe.Pointer(e))
	win, ok := wins[uint32(m.windowID)]
	if !ok {
		return
	}

	b, ok := mouseButtons[m.button]
	if !ok {
		return
	}

	x, y := int(m.x), int(m.y)
	if down {
		win.events <- ui.MouseDown{Button: b, X: x, Y: y}
	} else {
		win.events <- ui.MouseUp{Button: b, X: x, Y: y}
	}
}

func mouseMoveEvent(e *C.SDL_Event) {
	m := (*C.SDL_MouseMotionEvent)(unsafe.Pointer(e))
	win, ok := wins[uint32(m.windowID)]
	if !ok {
		return
	}
	win.events <- ui.MouseMove{X: int(m.x), Y: int(m.y)}
}

func winEvent(e *C.SDL_Event) {
	w := (*C.SDL_WindowEvent)(unsafe.Pointer(e))
	win, ok := wins[uint32(w.windowID)]
	if !ok || w.event != C.SDL_WINDOWEVENT_CLOSE {
		return
	}
	win.events <- ui.WinClose{}
}

type window struct {
	win    *C.SDL_Window
	rend   *C.SDL_Renderer
	w, h   int
	events chan interface{}
}

// NewWindow returns a new ui.Window.
func NewWindow(title string, w, h int) (ui.Win, error) {
	var err error
	win := &window{w: w, h: h}

	thread0.Do(func() {
		win.win = C.SDL_CreateWindow(
			C.CString(title),
			C.SDL_WINDOWPOS_UNDEFINED,
			C.SDL_WINDOWPOS_UNDEFINED,
			C.int(w),
			C.int(h),
			C.SDL_WINDOW_SHOWN|C.SDL_WINDOW_OPENGL)
		if win.win == nil {
			err = sdlError()
			win = nil
			return
		}
		win.rend = C.SDL_CreateRenderer(win.win, -1, C.SDL_RENDERER_ACCELERATED)
		if win.rend == nil {
			err = sdlError()
			C.SDL_DestroyWindow(win.win)
			win = nil
			return
		}

		win.events = make(chan interface{}, 10)
		wins[uint32(C.SDL_GetWindowID(win.win))] = win
	})

	return win, err
}

func (win *window) Present() {
	thread0.Do(func() {
		C.SDL_RenderPresent(win.rend)
	})
}

// Events returns the window's event channel.
func (win *window) Events() <-chan interface{} {
	return win.events
}

// Close closes the window and cleans up its resources.
func (win *window) Close() {
	thread0.Do(func() {
		id := uint32(C.SDL_GetWindowID(win.win))
		C.SDL_DestroyWindow(win.win)
		delete(wins, id)
		close(win.events)
	})
}

func sdlError() error {
	return errors.New(C.GoString(C.SDL_GetError()))
}
