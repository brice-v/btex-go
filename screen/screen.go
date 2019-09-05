package screen

import (
	"log"

	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/encoding"
)

// InitScreen initializes and returns the screen object
// style is set to bg='black' and fg='white'
// Unicode is registered here
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
