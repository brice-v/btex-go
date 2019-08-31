package main

import (
	"os"

	"github.com/gdamore/tcell"
)

func editorReadKey(s tcell.Screen) rune {
	var k rune

	ev := s.PollEvent()
	switch ev := ev.(type) {
	case *tcell.EventKey:
		switch ev.Key() {
		case tcell.KeyCtrlC, tcell.KeyCtrlQ:
			s.Fini()
			os.Exit(0)
		}
		k = ev.Rune()
	}
	return k
}

func editorRefreshScreen(s tcell.Screen) {
	s.Clear()
}

func main() {

	tcell.SetEncodingFallback(tcell.EncodingFallbackASCII)
	s, e := tcell.NewScreen()
	if e != nil {
		// TODO Remove panic for better method
		panic(e)
	}

	s.SetStyle(tcell.StyleDefault.
		Foreground(tcell.ColorBlack).
		Background(tcell.ColorWhite))
	s.Clear()
	s.Init()

	// quit := make(chan struct{})
	// go func() {
	// 	for {
	// 		ev := s.PollEvent()
	// 		switch ev := ev.(type) {
	// 		case *tcell.EventKey:
	// 			switch ev.Key() {
	// 			case tcell.KeyEscape, tcell.KeyCtrlQ, tcell.KeyCtrlC:
	// 				close(quit)
	// 				return
	// 			default:
	// 				fmt.Println(ev.Key())
	// 			}
	// 		case *tcell.EventResize:
	// 			s.Sync()
	// 		}
	// 	}
	// }()
	for {
		editorRefreshScreen()
		c := editorReadKey(s)
		print(string(c))
	}

}
