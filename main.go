package main

import (
	"fmt"
	"os"

	"github.com/gdamore/tcell"
)

const (
	BTEX_VERSION = "0.0.1"
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
	s.HideCursor()
	s.Clear()
	editorDrawRows(s)
	s.ShowCursor(1, 0)
}

// drawString sets the content at the starting location given by x and y
func drawString(s tcell.Screen, x int, y int, stringToDraw string) {
	bs := []rune(stringToDraw)
	for i, v := range bs {
		s.SetContent(x+i, y, v, nil, tcell.StyleDefault)
	}
	s.Show()
}

func editorDrawRows(s tcell.Screen) {
	w, h := s.Size()
	for y := 0; y < h; y++ {
		s.SetContent(0, y, '~', nil, tcell.StyleDefault)
	}
	// Draw Welcome Screen
	func(s tcell.Screen) {
		textToDraw := fmt.Sprintf("btex editor -- version %s", BTEX_VERSION)
		drawString(s, w/3, h/4, textToDraw)
		drawString(s, (w/3)-1, (h/4)+1, "Press Ctrl+C or Ctrl+Q to Quit")
	}(s)

}

func initEditor() tcell.Screen {
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
	return s
}

func main() {
	s := initEditor()

	for {
		editorRefreshScreen(s)

		s.Show()
		c := editorReadKey(s)
		print(string(c))

	}

}
