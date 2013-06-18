package sdl

/*
#include <SDL.h>

#cgo darwin CFLAGS: -I/Library/Frameworks/SDL2.framework/Headers
#cgo darwin LDFLAGS: -framework SDL2

#cgo linux CFLAGS: -I/usr/local/include/SDL2
#cgo linux LDFLAGS: -L/usr/local/lib -lSDL2

static Uint32 sdlEventType(SDL_Event *e){
	return e->type;
}
*/
import "C"

import (
	"errors"
	"time"
	"unsafe"

	"gui/ui"
)

// Wins is the map of all open windows.  It can only be accessed by the Init go routine.
var wins = make(map[uint32]*window, 1)

// Init initializes SDL.  It must be called from the main go routine, and it only
// returns if there is an error.  Hz specifies the frequency (in milliseconds)
// that SDL polls events.
func Init(hz int) error {
	if C.SDL_Init(C.SDL_INIT_EVERYTHING) != 0 {
		return sdlError()
	}
	go func() {
		tick := time.NewTicker(time.Duration(hz) * time.Millisecond)
		for _ = range tick.C {
			ui.Do(poll)
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

func keyEvent(e *C.SDL_Event, down bool) {
	k := (*C.SDL_KeyboardEvent)(unsafe.Pointer(e))
	if k.repeat != 0 {
		return
	}
	win, ok := wins[uint32(k.windowID)]
	if !ok {
		return
	}
	s := k.keysym.sym
	key, ok := keys[int(s)]
	if (s >= 'a' && s <= 'z') || (s >= 'A' && s <= 'Z') {
		key = ui.Key(s)
	} else if !ok {
		return
	}
	win.events <- ui.KeyEvent{
		Down: down,
		Key:  key,
	}
}

var winEvents = map[int]ui.WinEventType{
	C.SDL_WINDOWEVENT_CLOSE:        ui.WinClose,
	C.SDL_WINDOWEVENT_RESIZED:      ui.WinResize,
	C.SDL_WINDOWEVENT_SIZE_CHANGED: ui.WinResize,
	C.SDL_WINDOWEVENT_ENTER:        ui.WinEnter,
	C.SDL_WINDOWEVENT_LEAVE:        ui.WinLeave,
	C.SDL_WINDOWEVENT_FOCUS_GAINED: ui.WinFocus,
	C.SDL_WINDOWEVENT_FOCUS_LOST:   ui.WinUnFocus,
}

func winEvent(e *C.SDL_Event) {
	w := (*C.SDL_WindowEvent)(unsafe.Pointer(e))
	win, ok := wins[uint32(w.windowID)]
	if !ok {
		return
	}

	if w.event == C.SDL_WINDOWEVENT_RESIZED || w.event == C.SDL_WINDOWEVENT_SIZE_CHANGED {
		win.w = int(w.data1)
		win.h = int(w.data2)
	}

	if t, ok := winEvents[int(w.event)]; ok {
		win.events <- ui.WinEvent{
			Type:   t,
			Width:  win.w,
			Height: win.h,
		}
	}
}

type window struct {
	win    *C.SDL_Window
	w, h   int
	events chan interface{}
}

// NewWindow returns a new ui.Window.
func NewWindow(title string, w, h int) (ui.Win, error) {
	var err error
	win := &window{w: w, h: h}

	ui.Do(func() {
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

		win.events = make(chan interface{}, 10)
		wins[uint32(C.SDL_GetWindowID(win.win))] = win
	})

	return win, err
}

// Events returns the window's event channel.
func (win *window) Events() <-chan interface{} {
	return win.events
}

// Close closes the window and cleans up its resources.
func (win *window) Close() {
	ui.Do(func() {
		id := uint32(C.SDL_GetWindowID(win.win))
		C.SDL_DestroyWindow(win.win)
		delete(wins, id)
		close(win.events)
	})
}

func sdlError() error {
	return errors.New(C.GoString(C.SDL_GetError()))
}
