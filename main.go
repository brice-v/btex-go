package main

import "btex-go/editor"

func main() {
	e := editor.NewEditor()

	for {
		e.RefreshScreen()
		// this returns the rune but we may not need it
		_ = e.ProcessKey()
	}

}
