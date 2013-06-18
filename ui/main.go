// +build ignore

package main

import (
	"fmt"
	"gui/ui"
	"gui/sdl"
)

const (
	width  = 640
	height = 480
)

func main() {
	if err := sdl.Init(20); err != nil {
		panic(err)
	}
	go mainFunc()
	ui.Start()
}

func mainFunc() {
	win, err := sdl.NewWindow("Test", width, height)
	if err != nil {
		panic(err)
	}
	for ev := range win.Events() {
		fmt.Println(ev)
		if wEv, ok := ev.(ui.WinEvent); ok && wEv.Type == ui.WinClose {
			win.Close()
		}
	}
	ui.Stop()
}
