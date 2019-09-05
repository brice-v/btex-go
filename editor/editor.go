package editor

import (
	"fmt"
	"io/ioutil"
	"os"

	"btex-go/screen"

	"github.com/gdamore/tcell"
)

const (
	BTEX_VERSION = "0.0.1"
	NEWLINE_CHAR = '\n'
)

//
// EDITOR FUNCTIONS
//

type editorRow struct {
	length int
	chars  []rune
}

type Editor struct {
	s   tcell.Screen
	cur Cursor

	displayWelcome bool

	rows    []editorRow
	numrows int
}

func (E *Editor) displayCursor() {
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

func (E *Editor) deleteUnder() {
	E.s.SetContent(E.cur.x, E.cur.y, ' ', nil, tcell.StyleDefault)
	E.s.Show()
}

func (E *Editor) drawRune(r rune) {
	E.s.SetContent(E.cur.x, E.cur.y, r, nil, tcell.StyleDefault)
	E.cur.Move(RIGHT)
}

// ProcessKey polls the key pressed and responds with the correct event
// it will return the key if it is not a command key
func (E *Editor) ProcessKey() rune {
	var k rune

	ev := E.s.PollEvent()
	switch ev := ev.(type) {
	case *tcell.EventKey:
		switch ev.Key() {
		case tcell.KeyCtrlC, tcell.KeyCtrlQ:
			E.s.Fini()
			os.Exit(0)
		case tcell.KeyUp:
			E.cur.Move(UP)
		case tcell.KeyLeft:
			E.cur.Move(LEFT)
		case tcell.KeyRight:
			E.cur.Move(RIGHT)
		case tcell.KeyDown:
			// to handle directional issues with text i might just allow the user to Move anywhere and when they
			// start typing bring them back to the bottom of the text or something like that
			// still need to look in using ropes for string storage
			E.cur.Move(DOWN)
		case tcell.KeyBackspace2, tcell.KeyBackspace:
			E.cur.Move(LEFT)
			E.deleteUnder()
		case tcell.KeyDelete:
			// TODO need to Move all text left when something is deleted
			E.deleteUnder()
		case tcell.KeyEnter:
			// TODO come back to this because i dont know what this will need to do
			E.drawRune(NEWLINE_CHAR)
			E.cur.Move(DOWN)
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
func (E *Editor) RefreshScreen() {
	E.s.HideCursor()
	E.s.Clear()
	E.DrawRows()
	E.displayCursor()
	E.s.Show()
}

// DrawRows draws all the rows onto the screen from the E.row.chars
// this is going to change soon
func (E *Editor) DrawRows() {
	w, h := E.s.Size()
	for i := 0; i < h; i++ {
		E.s.SetContent(0, i, '~', nil, tcell.StyleDefault)
	}
	// Draw Welcome Screen
	if E.displayWelcome && len(E.rows) < 1 {
		textToDraw := fmt.Sprintf("btex editor -- version %s", BTEX_VERSION)
		DrawString(E.s, w/3, h/4, textToDraw)
		DrawString(E.s, (w/3)-1, (h/4)+1, "Press Ctrl+C or Ctrl+Q to Quit")
	}

}

//
// FILE / IO
//

func getRows(data []byte) []editorRow {
	ers := []editorRow{}
	buf := []byte{}

	length := 0
	indx := 0
	for _, char := range data {
		buf = append(buf, char)
		if char == '\n' {
			er := editorRow{length: length, chars: []rune(string(buf))}
			ers = append(ers, er)
			//save the newline and the byte slices lenth here
			buf = []byte{}
			length = 0
			indx++
		}
		length++
	}
	return ers
}

//OpenFile will open the file and set the buffers accordingly
func (E *Editor) OpenFile(f string) {
	data, err := ioutil.ReadFile(f)
	if err != nil {
		//TODO Better handle this failure
		return
	}
	E.rows = getRows(data)
	// fmt.Println(E.rows)
	// os.Exit(0)

}

//NewEditor returns the editor object
func NewEditor() *Editor {
	E := new(Editor)
	E.cur.x, E.cur.y = 1, 0
	E.s = screen.InitScreen()

	// for now only opening file when exactly the 1st argument on the command line
	if len(os.Args) == 2 {
		if _, err := os.Stat(os.Args[1]); err == nil {
			E.OpenFile(os.Args[1])
		} else if os.IsNotExist(err) {
			// // if it doesnt exist go ahead and create it
			// newFile, err := os.Create(os.Args[1])
			// if err != nil {
			// 	//TODO handle
			// 	panic(err)
			// }
			// E.OpenFile(newFile)
		} else {
			// something crazier happened?
			panic(err)
		}
	}
	E.numrows = len(E.rows)

	E.displayWelcome = true
	return E
}
