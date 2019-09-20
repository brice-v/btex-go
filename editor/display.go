package editor

import (
	"fmt"

	"github.com/gdamore/tcell"
	"github.com/mattn/go-runewidth"
)

const (
	// LEFTSIDE_CHAR is the character that appears on the leftmost
	// portion of the screen
	// TODO: Make this a settable property in the editor
	LEFTSIDE_CHAR = '~'
	// TAB_SIZE is the amount of spaces to display for tab
	// TODO: Make this a settable property in the editor
	TAB_SIZE = 4
)

// DrawRows draws all the rows onto the screen from the E.row.chars
// this is going to change soon
func (E *Editor) DrawRows() {
	style := tcell.StyleDefault.
		Foreground(tcell.ColorWhite).
		Background(tcell.ColorBlack)

	w, h := E.s.Size()
	totLen := E.pt.Length()
	// println("E.rowoffset", E.rowoffset)
	for i := 0; i < h; i++ {
		E.s.SetContent(0, i, '~', nil, style)
		if totLen != 0 {
			line, err := E.pt.GetLineStr(uint(i + E.rowoffset))
			if err != nil {
				continue
			}
			if E.rowoffset > 1 {
				E.s.Clear()
				E.Puts(style, 1, i+E.rowoffset, line)
			}
			E.Puts(style, 1, i, line)

		}

	}

	// Draw Welcome Screen
	if E.displayWelcome && totLen < 1 {
		textToDraw := fmt.Sprintf("btex editor -- version %s", BTEX_VERSION)
		E.Puts(style, w/3, h/4, textToDraw)
		E.Puts(style, (w/3)-1, (h/4)+1, "Press Ctrl+C or Ctrl+Q to Quit")
	}
}

// Puts paints a unicode string on to the display if it is supported
// This function is from: https://github.com/gdamore/tcell/blob/master/_demos/unicode.go
func (E *Editor) Puts(style tcell.Style, x, y int, str string) {
	i := 0
	var deferred []rune
	dwidth := 0
	zwj := false
	for _, r := range str {
		if r == '\n' || r == '\r' {
			// dont display newlines
			continue
		}
		if r == '\t' {
			for spaceCount := 0; spaceCount < TAB_SIZE; spaceCount++ {
				deferred = append(deferred, ' ')

			}
			dwidth = TAB_SIZE
			continue
		}
		if r == '\u200d' {
			if len(deferred) == 0 {
				deferred = append(deferred, ' ')
				dwidth = 1
			}
			deferred = append(deferred, r)
			zwj = true
			continue
		}
		if zwj {
			deferred = append(deferred, r)
			zwj = false
			continue
		}
		switch runewidth.RuneWidth(r) {
		case 0:
			if len(deferred) == 0 {
				deferred = append(deferred, ' ')
				dwidth = 1
			}
		case 1:
			if len(deferred) != 0 {
				E.s.SetContent(x+i, y, deferred[0], deferred[1:], style)
				i += dwidth
			}
			deferred = nil
			dwidth = 1
		case 2:
			if len(deferred) != 0 {
				E.s.SetContent(x+i, y, deferred[0], deferred[1:], style)
				i += dwidth
			}
			deferred = nil
			dwidth = 2
		}
		deferred = append(deferred, r)
	}
	if len(deferred) != 0 {
		E.s.SetContent(x+i, y, deferred[0], deferred[1:], style)
		i += dwidth
	}
}
