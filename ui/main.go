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

	setup()

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
			draw()
			win.Present()
		}
	}
	panic("Unreachable")
}

func setup() {
	thread0.Do(func() {
		v := strings.NewReader(vertShader)
		f := strings.NewReader(fragShader)
		var err error
		if prog, err = gl.NewProgram(v, f); err != nil {
			panic(err)
		}

		buf = gl.NewArrayBuffer()
		buf.SetData(
			gl.StaticDraw,
			0.75, 0.75, 0.0, 1.0,
			0.75, -0.75, 0.0, 1.0,
			-0.75, -0.75, 0.0, 1.0,
		)
	})
}

func draw() {
	thread0.Do(func() {
		gl.ClearColor(color.Black)
		gl.Clear(gl.ColorBufferBit)
		buf.Bind()
		prog.SetVertexAttributeData("position", 4, 0, 0)
		prog.DrawArrays(gl.Triangles, 0, 3)
	})
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
