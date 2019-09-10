package main

// an attempt at making a piecetable/piecemap in go but right now im using list(doubly) from container/list
// and honestly a bunch of other random stuff but im ready to start recording my eventual implementation

import (
	"container/list"
	"fmt"
	"io/ioutil"
	"log"
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
	//Remove is a NodeType used to cleanup the garbage nodes
	Remove
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
	} else if n.typ == Remove {
		result = fmt.Sprintf("{NodeType: Remove, start: %d, length: %d, lineOffsets: %v}}",
			n.start, n.length, n.lineOffsets)
	} else if n.typ == Sentinel {
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
			return
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

//CreateNode generates a new node
func (PT *PieceTable) CreateNode(typ NodeType, start, length int) *Node {
	buf := []rune(PT.buffer[typ][start : start+length])
	los := getLineOffsets(buf)

	return &Node{
		typ:         typ,
		start:       start,
		length:      length,
		lineOffsets: los,
	}
}

// ------------------------------------------------------------------------
//
// REMOVE FUNCTIONS
//

func (PT *PieceTable) cleanupRemoveNodes() {
	// this is probably really slow but i couldnt make a better way
	// cleanup all the removes
	for PT.anyRemoveLeft() {
		for e := PT.nodes.Front(); e != nil; e = e.Next() {
			n, ok := e.Value.(*Node)
			if !ok {
				log.Fatal("Found non Node when trying to unwrap in deletestringat")
			}
			if n.typ == Remove {
				PT.deleteNode(e)
				break
			}
		}
	}
}

// DeleteStringAt will delete the string in the nodes from
// start is the char the delete starts
// length is the length of the delete. to delete only 1 char it will be a 0 length
// EXAMPLE ---------------------------------------------------------------------------
// DeleteStringAt(start=4,length=4)
//
// Node: Start=0, Length=10
// ]<->[ 1,2,3,4,5,6,7,8,9,10]<->[..
//
// NodeLeft: Start=currentNodeStart, length=start-currentNodeStart
// NodeRight: Start=currentNodeStart+length, length=currentNodeLength-length
// ]<->[ 1,2,3]<->[8,9,10]<->[..
// ------------------------------------------------------------------------------------
func (PT *PieceTable) DeleteStringAt(offset, length int) error {
	// return error when trying to use negative values
	if offset < 0 || length < 0 {
		return fmt.Errorf("Offset or Length is less than 0. offset=%d, length=%d",
			offset, length)
	}

	//totLen records the total length of the visible buffer as we continue through
	totLen := 0

	// cleanup when were done with the initial loop
	defer PT.cleanupRemoveNodes()

	//first have to find which node this starts in
	for e := PT.nodes.Front(); e != nil; e = e.Next() {
		n, ok := e.Value.(*Node)
		if !ok {
			log.Fatal("Found non Node when trying to unwrap in deletestringat")
		}
		if n.typ == Sentinel || n.typ == Remove {
			continue
		}

		// add each nodes length to get the current place in the visible buffer
		totLen += n.length
		endLen := offset + length
		// this is the nodes start point in the documents visible buffer
		nodeStartPoint := totLen - n.length

		distanceToRightNodeInChars := totLen - offset
		distanceFromLeftToOffset := n.length - distanceToRightNodeInChars

		// still need to keep going if we arent at the offset yet
		if offset > totLen || nodeStartPoint > endLen {
			continue
		} else if totLen > offset && endLen >= totLen {
			// delete the node entirely if we are beyond the start of the offset and the
			// offset + length is still more than the current total length
			n.typ = Remove
			continue
		} else if totLen > offset && endLen < totLen && nodeStartPoint <= offset {
			//in this case we remove the node were in, and make sure to add a new node if necessary
			// for the remainder of end offset to the totlen
			nodeLeft := PT.CreateNode(n.typ, n.start, distanceFromLeftToOffset)

			// only insert the left node if it has a length
			// ignore negatives just in case?:
			if nodeLeft.length > 0 {
				PT.nodes.InsertBefore(nodeLeft, e)
			}

			nodeRightStart := distanceFromLeftToOffset + length
			nodeRightLength := n.length - nodeRightStart
			nodeRight := PT.CreateNode(n.typ, nodeRightStart, nodeRightLength)
			PT.nodes.InsertBefore(nodeRight, e)
			n.typ = Remove
			break
		} else if totLen > offset && endLen > totLen && nodeStartPoint < offset {
			// this is only node left
			nodeLeft := PT.CreateNode(n.typ, n.start, distanceFromLeftToOffset)
			// only insert the left node if it has a length
			if nodeLeft.length != 0 {
				PT.nodes.InsertBefore(nodeLeft, e)
				n.typ = Remove
			}
		} else if totLen > offset && endLen < totLen && nodeStartPoint > offset {
			// this is only the right node
			nodeRightStart := (n.length - (totLen - endLen)) + n.start
			nodeRightLength := n.length - (nodeRightStart - n.start)

			nodeRight := PT.CreateNode(n.typ, nodeRightStart, nodeRightLength)
			PT.nodes.InsertBefore(nodeRight, e)
			n.typ = Remove
			break
		} else if totLen == offset && length == 1 {
			nodeLeftLength := totLen - 1
			nodeLeft := PT.CreateNode(n.typ, n.start, nodeLeftLength)
			PT.nodes.InsertBefore(nodeLeft, e)
			n.typ = Remove
			break
		} else {
			return fmt.Errorf("Case not handled totlen=%v, offset=%v", totLen, offset)
		}
	}
	return nil
}

func (PT *PieceTable) anyRemoveLeft() bool {
	for e := PT.nodes.Front(); e != nil; e = e.Next() {
		n, ok := e.Value.(*Node)
		if !ok {
			log.Fatal("Found non Node when trying to unwrap in deletestringat")
		}
		if n.typ == Remove {
			return true
		}

	}
	return false
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
	pt.InsertStringAt(6, "AAABBB")
	pt.InsertStringAt(10, "CCC")
	// need to get this working
	pt.DeleteStringAt(0, 12)
	// need to get this working
	// pt.DeleteStringAt(3, 10)
	// pt.DeleteStringAt(3, 8)
	// need to get this working
	// pt.DeleteStringAt(7, 1)

	for e := pt.nodes.Front(); e != nil; e = e.Next() {
		n := e.Value.(*Node)
		fmt.Println(n)
	}
	cat(pt)

}
