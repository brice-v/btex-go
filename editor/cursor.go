package editor

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

type Cursor struct {
	x int
	y int
}

func (c *Cursor) Move(d direction) {
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
