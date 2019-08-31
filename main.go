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
	default:
		return k
	}
	return k
}

func editorRefreshScreen(s tcell.Screen) {
	s.Clear()
}

func editorDrawRows(s tcell.Screen) {
	for y := 0; y < 24; y++ {
		s.SetContent(0, y, '~', nil, tcell.StyleDefault)
	}

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

	for i := 0; i < 100; i++ {
		editorRefreshScreen(s)
		editorDrawRows(s)
		s.ShowCursor(1, 0)
		s.Show()
		c := editorReadKey(s)
		print(string(c))

	}

}
