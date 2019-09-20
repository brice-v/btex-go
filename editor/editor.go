package editor

import (
	"io/ioutil"
	"log"
	"os"

	"btex-go/peace"
	"btex-go/screen"

	"github.com/gdamore/tcell"
)

const (
	//BTEX_VERSION is the current version of the text editor
	BTEX_VERSION = "0.0.2"
	// TODO: Possibly remove
	// NEWLINE_CHAR = '\n'
)

//
// EDITOR FUNCTIONS
//

//Editor holds all the information related to the editor
// TODO: export fields and make this a more suitable
// embeddable project
type Editor struct {
	s   tcell.Screen
	cur Cursor

	displayWelcome bool

	pt        *peace.PieceTable
	rowoffset int
}

func (E *Editor) displayCursor() {
	// TODO: restrict how the cursor can be displayed in an editing mode
	w, h := E.s.Size()

	if E.cur.x < 1 {
		E.cur.x = 1
	}
	if E.cur.y < 0 {
		E.cur.y = 0
		E.rowoffset--
	}
	if E.cur.x > w {
		E.cur.x = w - 1
	}
	if E.cur.y > h {
		E.cur.y = h - 1
		E.rowoffset++
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
			// E.drawRune(NEWLINE_CHAR)
			// E.cur.Move(DOWN)
			// E.cur.x = 1
			// E.s.Show()
		case tcell.KeyCtrlL:
			E.s.Clear()
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

//
// FILE / IO
//

//OpenFile will open the file and set the buffers accordingly
func (E *Editor) OpenFile(f string) {
	data, err := ioutil.ReadFile(f)
	if err != nil {
		//TODO Better handle this failure
		return
	}

	E.pt = peace.NewPT(([]rune(string(data))))

}

//NewEditor returns the editor object
func NewEditor() *Editor {
	E := new(Editor)
	E.cur.x, E.cur.y = 1, 0

	style := tcell.StyleDefault.
		Foreground(tcell.ColorWhite).
		Background(tcell.ColorBlack)

	E.s = screen.InitScreen(style)

	E.pt = peace.NewPT(nil)

	E.rowoffset = 2

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
			log.Fatal(err)
		}
	}
	E.displayWelcome = true
	return E
}
