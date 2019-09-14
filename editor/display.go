package editor

import (
	"github.com/gdamore/tcell"
	"github.com/mattn/go-runewidth"
)

const (
	LEFTSIDE_CHAR = '~'
)

// drawString sets the content at the starting location given by x and y
func DrawString(s tcell.Screen, x int, y int, stringToDraw string) {
	bs := []rune(stringToDraw)
	for i, v := range bs {
		s.SetContent(x+i, y, v, nil, tcell.StyleDefault)
	}
	s.Show()
}

func (E *Editor) DrawEditorChars(xPos int, yPos int, leftChar rune) {
	w, h := E.s.Size()
	for rowindx := 0; rowindx < E.numrows && rowindx < h; rowindx++ {
		for charindx := uint(0); charindx < E.rows[rowindx].length && charindx < uint(w); charindx++ {
			E.s.SetContent(int(charindx+2), rowindx, E.rows[rowindx].chars[charindx], nil, tcell.StyleDefault)
		}
		// lineNo := fmt.Sprintf("%d", rowindx)
		E.s.SetContent(0, rowindx, LEFTSIDE_CHAR, nil, tcell.StyleDefault)
		E.s.SetContent(1, rowindx, ' ', nil, tcell.StyleDefault)
	}

	E.s.Show()
}

// puts paints a unicode string on to the display if it is supported
// This function is from: https://github.com/gdamore/tcell/blob/master/_demos/unicode.go
func puts(s tcell.Screen, style tcell.Style, x, y int, str string) {
	i := 0
	var deferred []rune
	dwidth := 0
	zwj := false
	for _, r := range str {
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
				s.SetContent(x+i, y, deferred[0], deferred[1:], style)
				i += dwidth
			}
			deferred = nil
			dwidth = 1
		case 2:
			if len(deferred) != 0 {
				s.SetContent(x+i, y, deferred[0], deferred[1:], style)
				i += dwidth
			}
			deferred = nil
			dwidth = 2
		}
		deferred = append(deferred, r)
	}
	if len(deferred) != 0 {
		s.SetContent(x+i, y, deferred[0], deferred[1:], style)
		i += dwidth
	}
}
