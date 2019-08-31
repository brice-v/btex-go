package main

import (
	"fmt"
	"os"

	"github.com/gdamore/tcell"
)

const (
	BTEX_VERSION = "0.0.1"
)

// drawString sets the content at the starting location given by x and y
func drawString(s tcell.Screen, x int, y int, stringToDraw string) {
	bs := []rune(stringToDraw)
	for i, v := range bs {
		s.SetContent(x+i, y, v, nil, tcell.StyleDefault)
	}
	s.Show()
}

//
// CURSOR FUNCTIONS
//

type direction int

const (
	UP direction = iota
	DOWN
	LEFT
	RIGHT
)

type cursor struct {
	x int
	y int
}

func (c *cursor) move(d direction) {
	switch d {
	case UP:
		c.y--
	case DOWN:
		c.y++
	case LEFT:
		c.x--
	case RIGHT:
		c.x++
	}
}

//
// EDITOR FUNCTIONS
//

type editor struct {
	s   tcell.Screen
	cur cursor

	displayWelcome bool
}

func (E *editor) ReadKey() rune {
	var k rune

	ev := E.s.PollEvent()
	switch ev := ev.(type) {
	case *tcell.EventKey:
		switch ev.Key() {
		case tcell.KeyCtrlC, tcell.KeyCtrlQ:
			E.s.Fini()
			os.Exit(0)
		}
		k = ev.Rune()
		switch k {
		case 'w':
			E.cur.move(UP)
		case 'a':
			E.cur.move(LEFT)
		case 'd':
			E.cur.move(RIGHT)
		case 's':
			E.cur.move(DOWN)
		}
		// as soon as typing begins, get rid of the welcome screen
		E.displayWelcome = false
	default:
		return k
	}
	return k
}

func (E *editor) RefreshScreen() {
	E.s.HideCursor()
	E.s.Clear()
	E.initRows()
	E.s.ShowCursor(E.cur.x, E.cur.y)
	E.s.Show()
}

func (E *editor) initRows() {
	w, h := E.s.Size()
	for y := 0; y < h; y++ {
		E.s.SetContent(0, y, '~', nil, tcell.StyleDefault)
	}
	// Draw Welcome Screen
	if E.displayWelcome {
		textToDraw := fmt.Sprintf("btex editor -- version %s", BTEX_VERSION)
		drawString(E.s, w/3, h/4, textToDraw)
		drawString(E.s, (w/3)-1, (h/4)+1, "Press Ctrl+C or Ctrl+Q to Quit")
	}

}

func initScreen() tcell.Screen {
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

func newEditor() *editor {
	E := new(editor)
	E.cur.x, E.cur.y = 1, 0
	E.s = initScreen()
	E.displayWelcome = true
	return E
}

func main() {
	e := newEditor()

	for {
		e.RefreshScreen()
		c := e.ReadKey()
		print(string(c))
	}

}
