// +build ignore

package main

import (
	"errors"
	"fmt"
	"image/color"
	"os"
	"time"

	"github.com/eaburns/gui/gl"
	"github.com/eaburns/gui/sdl"
	"github.com/eaburns/gui/thread0"
	"github.com/eaburns/gui/ui"
)

const (
	width  = 640
	height = 480
)

func main() {
	if err := sdl.Init(20); err != nil {
		panic(err)
	}
	gl.Init(0, width-1, 0, height-1)
	go mainFunc()
	thread0.Hijack()
}

func mainFunc() {
	win, err := sdl.NewWindow("Test", width, height)
	if err != nil {
		panic(err)
	}
	tick := time.NewTicker(20 * time.Millisecond)
	for {
		select {
		case ev, ok := <-win.Events():
			if !ok {
				os.Exit(0)
			}
			fmt.Printf("%T%v\n", ev, ev)
			if _, ok := ev.(ui.WinClose); ok {
				win.Close()
			}
		case <-tick.C:
			if err := draw(); err != nil {
				panic(err)
			}
			win.Present()
		}
	}
	panic("Unreachable")
}

func draw() error {
	var err error
	thread0.Do(func() {
		gl.Color(color.White)
		gl.ClearColorBuffer()
		gl.BeginQuads()
		gl.Vertex2(100, 100)
		gl.Vertex2(200, 100)
		gl.Vertex2(200, 200)
		gl.Vertex2(100, 200)
		gl.End()
		if s, ok := gl.ErrorString[gl.GetError()]; ok && s != "GL_NO_ERROR" {
			err = errors.New(s)
		}
	})
	return err
}
