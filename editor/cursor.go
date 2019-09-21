package editor

//
// CURSOR FUNCTIONS
//

type direction int

const (
	// UP direction
	UP direction = iota
	// DOWN direction
	DOWN
	// LEFT direction
	LEFT
	// RIGHT direction
	RIGHT
)

// Cursor is the current position of the cursor on the screen
// TODO: Implement multicursor
type Cursor struct {
	x int
	y int
}

// Move increments or decrements the cursor's position based on
// the direction
func (c *Cursor) Move(d direction) {
	switch d {
	case UP:
		if c.y == 0 {
			return
		}
		c.y--
	case DOWN:
		c.y++
	case LEFT:
		if c.x == 0 {
			return
		}
		c.x--
	case RIGHT:
		c.x++
	}
}
