package main

// an attempt at making a piecetable/piecemap in go but right now im using list(doubly) from container/list
// and honestly a bunch of other random stuff but im ready to start recording my eventual implementation

import (
	"container/list"
	"fmt"
	"io/ioutil"
	"log"
	"math"
)

var println = fmt.Println

// NodeType Enum they are both readonly/appendonly buffers
type NodeType int

const (
	//Added Buffer NodeType descriptor
	Added NodeType = iota
	//Original Buffer NodeType descriptor
	Original
	//Sentinel is the NodeType that refers to the empty head and tail nodes
	Sentinel
)

// PieceTable is currently 2 buffers but will be modified in the future
type PieceTable struct {
	buffer map[NodeType][]rune
	nodes  *list.List
}

func (n *Node) String() (result string) {
	if n.typ == Original {
		result = fmt.Sprintf("{NodeType: Original, start: %d, length: %d, lineOffsets: %v}",
			n.start, n.length, n.lineOffsets)
	} else if n.typ == Added {
		result = fmt.Sprintf("{NodeType: Added, start: %d, length: %d, lineOffsets: %v}}",
			n.start, n.length, n.lineOffsets)
	} else {
		result = fmt.Sprintf("{NodeType: Sentinel, start: %d, length: %d, lineOffsets: %v}}",
			n.start, n.length, n.lineOffsets)
	}
	return
}

// Node is the element in the list that contains some metadata for the contents and the operation
type Node struct {
	typ    NodeType
	start  int
	length int

	lineOffsets []int
}

func (PT *PieceTable) deleteNode(e *list.Element) {
	// remove the node we are currently on top of and
	// break the loop because we are done (this could be a return)
	PT.nodes.MoveToBack(e)
	if this, ok := PT.nodes.Back().Value.(*Node); ok {
		if this.typ != Sentinel {
			PT.nodes.Remove(PT.nodes.Back())
		}
	}
}

func getLineOffsets(buf []rune) []int {
	var bucket []int
	for i := 0; i < len(buf); i++ {
		if buf[i] == '\n' {
			bucket = append(bucket, i)
		}
	}
	return bucket
}

// ------------------------------------------------------------------------
//
// REMOVE FUNCTIONS
//

// DeleteStringAt will delete the string in the nodes from
// start is the char the delete starts
// length is the length of the delete. to delete only 1 char it will be a 0 length
func (PT *PieceTable) DeleteStringAt(start, length int) error {
	if start < 0 {
		return fmt.Errorf("Start must be positive")
	}
	// var endOfDeleteStringInDoc int
	// if length < 0 {
	// 	endOfDeleteStringInDoc = start
	// } else {
	// 	endOfDeleteStringInDoc = start + length
	// }

	//totLen records the total length of the visible buffer as we continue through
	totLen := 0

	//first have to find which node this starts in
	for e := PT.nodes.Front(); e != nil; e = e.Next() {
		n, ok := e.Value.(*Node)
		if !ok {
			log.Fatal("Found non Node when trying to unwrap in deletestringat")
		}
		if n.typ == Sentinel {
			continue
		}

		// add each nodes length to get the current place in the visible buffer
		totLen += n.length
		offset := start

		// if the length is negative we calculate a different starting point
		// otherwise it stays as the offset
		if length < 0 {
			offset = start - int(math.Abs(float64(length)))
			if offset < 0 {
				return fmt.Errorf("Offset was calculated to be less than 0 this is not allowed")
			}
			length = int(math.Abs(float64(length)))
		}

		// still need to keep going if we arent at the offset yet
		if offset > totLen {
			continue
		} else if totLen > offset && (offset+length < totLen) {
			// so the code below just covers the one case where the node we are inside of is the correct node to delete
			// we now need to cover a delete spanning multiple nodes
			// this is just the node that we are starting on

			//in this case we remove the node were in, and make sure to add a new node if necessary
			// for the remainder of end offset to the totlen

			// EXAMPLE ---------------------------------------------------------------------------
			// DeleteStringAt(start=3,length=4)
			//
			// Node: Start=0, Length=10
			// ]<->[ 1,2,3,4,5,6,7,8,9,10]<->[..
			//
			// NodeLeft: Start=currentNodeStart, length=start-currentNodeStart
			// NodeRight: Start=currentNodeStart+length, length=currentNodeLength-length
			// ]<->[ 1,2,3]<->[8,9,10]<->[..
			// ------------------------------------------------------------------------------------

			nodeLeftStart := n.start
			nodeLeftLength := offset + n.start
			nodeLeftBuf := []rune(PT.buffer[n.typ][nodeLeftStart : nodeLeftStart+nodeLeftLength])
			nodeLeftLos := getLineOffsets(nodeLeftBuf)

			nodeRightStart := n.start + length
			nodeRightLength := n.length - nodeRightStart
			nodeRightBuf := []rune(PT.buffer[n.typ][nodeRightStart : nodeRightStart+nodeRightLength])
			nodeRightLos := getLineOffsets(nodeRightBuf)

			nodeLeft := &Node{
				typ:         n.typ,
				start:       nodeLeftStart,
				length:      nodeLeftLength,
				lineOffsets: nodeLeftLos,
			}
			nodeRight := &Node{
				typ:         n.typ,
				start:       nodeRightStart,
				length:      nodeRightLength,
				lineOffsets: nodeRightLos,
			}
			if offset != n.start {
				PT.nodes.InsertBefore(nodeLeft, e)
			}
			PT.nodes.InsertBefore(nodeRight, e)
			PT.deleteNode(e)

			return nil

		} else if totLen == offset && length == 1 {
			nodeLeftStart := n.start
			nodeLeftLength := totLen - 1
			nodeLeftBuf := []rune(PT.buffer[n.typ][nodeLeftStart : nodeLeftStart+nodeLeftLength])
			nodeLeftLos := getLineOffsets(nodeLeftBuf)

			nodeLeft := &Node{
				typ:         n.typ,
				start:       nodeLeftStart,
				length:      nodeLeftLength,
				lineOffsets: nodeLeftLos,
			}
			PT.nodes.InsertBefore(nodeLeft, e)
			PT.deleteNode(e)
			return nil
		} else {
			return fmt.Errorf("Case not handled totlen=%v, offset=%v, start=%v", totLen, offset, start)
		}

	}
	return fmt.Errorf("Should not make it out of the node for loop")
}

// ------------------------------------------------------------------------
//
// ADD FUNCTIONS
//

//AppendString allows a new string to be added to the add buffer
// this is strictly for append
// just syntactic
func (PT *PieceTable) AppendString(data string) {
	PT.InsertStringAt(len(PT.buffer[Added])+len(PT.buffer[Original]), data)
}

//InsertStringAt will insert a string into the piece table at an offset
// this offset can be considered the byte location from the beginning of
// the visible buffers.
// data is the string to (append to the add buffer) be added to the
// PieceTable
func (PT *PieceTable) InsertStringAt(offset int, data string) bool {
	//record where we are in the document pretty much
	totLen := 0
	// -------------------------------------------------------------
	// This is the new node that we are adding to the `Added` buffer
	// because this new Node in the dll is getting inserted at an offset
	// that is what actually determines its place in the dll and not
	// the content that refers to the view on the buffer
	// -------------------------------------------------------------
	nodeMiddleTyp := Added
	// calulate line offsets for the newly inserted data
	nodeMiddleLos := getLineOffsets([]rune(data))
	//the start of the new node is the current length of the `Added`
	// buffer because that is where this new data will be visible from.
	nodeMiddleStart := len(PT.buffer[Added])
	// this is the length of the data we are passing in
	nodeMiddleLength := len(data)
	// append the rest of the string to the add buffer
	PT.buffer[Added] = append(PT.buffer[Added], []rune(data)...)
	newNodeMiddle := &Node{
		typ:         nodeMiddleTyp,
		start:       nodeMiddleStart,
		length:      nodeMiddleLength,
		lineOffsets: nodeMiddleLos,
	}

	// currentTotLen := 0
	// looop through the nodes and find out where the offset is gonna be, use the length += next length to
	for e := PT.nodes.Front(); e != nil; e = e.Next() {
		n, ok := e.Value.(*Node)
		if !ok {
			panic("Not unrwapping a node")
		}
		if n.typ == Sentinel {
			continue
		}

		//
		// SOME NOTES
		//
		// if the offset is in the middle of a nodes start -> start+length
		// 		=> then this is the one we need to "Split Up"
		// mainly meaning that this node will have to get removed (from the dll)
		// and 2 new nodes will be made (to fill in the gaps left and right with
		// the newNode in the middle {previously where the original node was})
		// with the proper start and length that will make up where it got "split"
		// from the inserted node
		//
		//

		totLen += n.length

		if offset > totLen {
			continue
		} else if offset < totLen {
			// insert 3 new nodes, left | middle (new data) | right

			//lets create our new nodes for the 2 new views
			// start and type are the same
			nodeLeftTyp := n.typ
			nodeLeftStart := n.start
			// i think this is right
			nodeLeftLength := n.length - (totLen - offset)
			// hopefully this works
			nodeLeftLos := getLineOffsets(PT.buffer[n.typ][n.start:offset])
			newNodeLeft := &Node{
				typ:         nodeLeftTyp,
				start:       nodeLeftStart,
				length:      nodeLeftLength,
				lineOffsets: nodeLeftLos,
			}

			//new node for the right
			nodeRightTyp := n.typ
			nodeRightStart := nodeLeftLength
			nodeRightLenth := totLen - offset
			nodeRightLos := getLineOffsets(PT.buffer[n.typ][nodeRightStart : nodeRightLenth+nodeRightStart])
			newNodeRight := &Node{typ: nodeRightTyp, start: nodeRightStart, length: nodeRightLenth, lineOffsets: nodeRightLos}

			//now that we have our nodes, load them into the dll
			// and then remove the node we are currently on
			PT.nodes.InsertBefore(newNodeLeft, e)
			PT.nodes.InsertBefore(newNodeMiddle, e)
			PT.nodes.InsertBefore(newNodeRight, e)

			PT.deleteNode(e)
			return true
		} else if offset == totLen {
			//insert between 2 dll nodes
			PT.nodes.InsertAfter(newNodeMiddle, e)
			return true
		}

	}
	return false
}

// ------------------------------------------------------------------------

func newEmptyList() *list.List {
	hn := &Node{typ: Sentinel, start: 0, length: 0, lineOffsets: []int{}}
	tn := &Node{typ: Sentinel, start: 0, length: 0, lineOffsets: []int{}}
	l := list.New()
	l.PushBack(tn)
	l.PushBack(hn)
	return l
}

// NewPT will eventually return a piecetable/map and will probably have a separate
// new function for the optional buffer (this would be starting a new buffer for instance)
func NewPT(optBuf []rune) *PieceTable {
	if optBuf != nil {
		optBufLen := len(optBuf)
		bufs := map[NodeType][]rune{Original: optBuf, Added: []rune("")}
		pt := &PieceTable{buffer: bufs, nodes: newEmptyList()}
		//calculate lineoffsets
		los := getLineOffsets(optBuf)
		newAppendNode := &Node{
			typ:         Original,
			start:       0,
			length:      optBufLen,
			lineOffsets: los,
		}
		pt.nodes.InsertBefore(newAppendNode, pt.nodes.Back())
		return pt
	}
	bufs := map[NodeType][]rune{Original: []rune(""), Added: []rune("")}
	return &PieceTable{buffer: bufs, nodes: newEmptyList()}
}

func openAndReadFile(f string) []rune {
	data, err := ioutil.ReadFile(f)
	if err != nil {
		log.Fatal(err)
	}
	return []rune(string(data))
}

func cat(pt *PieceTable) {
	for e := pt.nodes.Front(); e != nil; e = e.Next() {
		n, ok := e.Value.(*Node)
		if !ok {
			panic("Not unrwapping a node")
		}
		if n.typ == Original {
			fmt.Print(string(pt.buffer[Original][n.start : n.start+n.length]))
		} else if n.typ == Added {
			fmt.Print(string(pt.buffer[Added][n.start : n.start+n.length]))
		} else {
			// e = e.Next()
			continue
		}
	}
}

func main() {

	// data := openAndReadFile("unicode.txt")

	data := []rune(`Thequi Σckbrown`)
	// println("len(input)=", len(input))
	pt := NewPT(data)
	pt.InsertStringAt(6, "AAA")
	// need to get this working
	// pt.DeleteStringAt(7, 1)
	// need to get this working
	pt.DeleteStringAt(0, 8)
	// need to get this working
	// pt.DeleteStringAt(7, 1)

	// fmt.Println(result)

	// pt.AppendString(`//EXTRA
	// 	asfjk

	// 	// data to have at the bottom test`)

	// 	pt.InsertStringAt(6, `Here is the new
	// data`)

	// 	pt.AppendString("\n||||||||||||||||||\n")

	// pt.InsertStringAt(28, `Here is the new afjslkjasflkjasflk
	// afskjfaskasfljfa
	// asfjasfkjfkasj
	// data`)
	for e := pt.nodes.Front(); e != nil; e = e.Next() {
		n := e.Value.(*Node)
		fmt.Println(n)
	}
	cat(pt)

}
