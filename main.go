package main

import (
	"fmt"
	"os"

	"github.com/gdamore/tcell"
)

const (
	BTEX_VERSION  = "0.0.1"
	NEWLINE_CHAR  = '\n'
	LEFTSIDE_CHAR = '~'
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
// FILE / IO
//

func (E *editor) fileOpen() {
	open_file := `
"HelloWorld"!
Sometimes there is more words


more stuff down here`
	E.row.size = len(open_file)
	E.row.chars = []rune(open_file)

	E.numrows = 6
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

type editorRow struct {
	size  int
	chars []rune
}

type editor struct {
	s   tcell.Screen
	cur cursor

	displayWelcome bool

	row     editorRow
	numrows int
}

func (E *editor) displayCursor() {
	w, h := E.s.Size()

	if E.cur.x < 1 {
		E.cur.x = 1
	}
	if E.cur.y < 0 {
		E.cur.y = 0
	}
	if E.cur.x > w {
		E.cur.x = w - 1
	}
	if E.cur.y > h {
		E.cur.y = h - 1
	}
	E.s.ShowCursor(E.cur.x, E.cur.y)
}

func (E *editor) deleteUnder() {
	E.s.SetContent(E.cur.x, E.cur.y, 'A', nil, tcell.StyleDefault)
	E.s.Show()
}

func (E *editor) drawRune(r rune) {
	E.cur.move(RIGHT)
	E.s.SetContent(E.cur.x, E.cur.y, r, nil, tcell.StyleDefault)
}

// ProcessKey polls the key pressed and responds with the correct event
// it will return the key if it is not a command key
func (E *editor) ProcessKey() rune {
	var k rune

	ev := E.s.PollEvent()
	switch ev := ev.(type) {
	case *tcell.EventKey:
		switch ev.Key() {
		case tcell.KeyCtrlC, tcell.KeyCtrlQ:
			E.s.Fini()
			os.Exit(0)
		case tcell.KeyUp:
			E.cur.move(UP)
		case tcell.KeyLeft:
			E.cur.move(LEFT)
		case tcell.KeyRight:
			E.cur.move(RIGHT)
		case tcell.KeyDown:
			// to handle directional issues with text i might just allow the user to move anywhere and when they
			// start typing bring them back to the bottom of the text or something like that
			// still need to look in using ropes for string storage
			E.cur.move(DOWN)
		case tcell.KeyBackspace2, tcell.KeyBackspace:
			E.cur.move(LEFT)
			E.deleteUnder()
		case tcell.KeyDelete:
			// TODO need to move all text left when something is deleted
			E.deleteUnder()
		case tcell.KeyEnter:
			E.drawRune(NEWLINE_CHAR)
			E.cur.move(DOWN)
			E.cur.x = 1
			E.s.Show()
		default:
			k = ev.Rune()
			E.drawRune(k)
		}
		// as soon as typing begins, get rid of the welcome screen
		E.displayWelcome = false
	default:
		return k
	}
	return k
}

// RefreshScreen calls all the necessary functions between terminal screen refreshes
func (E *editor) RefreshScreen() {
	E.s.HideCursor()
	E.s.Clear()
	E.initRows()
	E.displayCursor()
	E.s.Show()
}

func (E *editor) initRows() {
	w, h := E.s.Size()
	for y := 0; y < h; y++ {
		E.s.SetContent(0, y, LEFTSIDE_CHAR, nil, tcell.StyleDefault)
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

	E.fileOpen()

	E.displayWelcome = true
	return E
}

func main() {
	e := newEditor()

	for {
		e.RefreshScreen()
		c := e.ProcessKey()
		print(string(c))
	}

}
