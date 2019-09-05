package screen

import (
	"log"

	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/encoding"
)

func InitScreen() tcell.Screen {
	//unicode support
	encoding.Register()
	s, e := tcell.NewScreen()
	if e != nil {
		// log error instead of panic
		log.Fatal(e)
	}
	if e := s.Init(); e != nil {
		log.Fatal(e)
	}

	s.SetStyle(tcell.StyleDefault.
		Foreground(tcell.ColorWhite).
		Background(tcell.ColorBlack))

	return s
}
