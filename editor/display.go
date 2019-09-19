package editor

import (
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

// drawString sets the content at the starting location given by x and y
func DrawString(s tcell.Screen, x int, y int, stringToDraw string) {
	bs := []rune(stringToDraw)
	for i, v := range bs {
		s.SetContent(x+i, y, v, nil, tcell.StyleDefault)
	}
	s.Show()
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
