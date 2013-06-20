// +build ignore

package main

import (
	"fmt"
	"image/color"
	"os"
	"strings"
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

var (
	prog *gl.Program
	buf  *gl.Buffer
)

func main() {
	if err := sdl.Init(20); err != nil {
		panic(err)
	}
	go mainFunc()
	thread0.Hijack()
}

func mainFunc() {
	win, err := sdl.NewWindow("Test", width, height)
	if err != nil {
		panic(err)
	}

	if err := setup(); err != nil {
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

func setup() error {
	var err error
	thread0.Do(func() {
		v := strings.NewReader(vertShader)
		f := strings.NewReader(fragShader)
		prog, err = gl.NewProgram(v, f)
		if err != nil {
			return
		}

		buf, err = gl.NewArrayBuffer()
		if err != nil {
			return
		}

		err = buf.SetData(
			gl.StaticDraw,
			0.75, 0.75, 0.0, 1.0,
			0.75, -0.75, 0.0, 1.0,
			-0.75, -0.75, 0.0, 1.0,
		)
	})
	return err
}

func draw() error {
	var err error
	thread0.Do(func() {
		gl.ClearColor(color.Black)
		if err = gl.Clear(gl.ColorBufferBit); err != nil {
			return
		}

		if err = buf.Bind(); err != nil {
			return
		}

		if err = prog.SetVertexAttributeData("position", 4, 0, 0); err != nil {
			return
		}

		err = prog.DrawArrays(gl.Triangles, 0, 3)
	})
	return err
}

var (
	vertShader = `
	        #version 120
	        
	        attribute vec4 position;
	        
	        void main(){
	               gl_Position = position;
	        }`

	fragShader = `
	        #version 120
	        
	        void main(){
	               gl_FragColor = vec4(1, 1, 1, 1);
	        }`
)
