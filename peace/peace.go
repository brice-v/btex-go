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

func (PT *PieceTable) newNodeAppendOnly(typ NodeType, start, length int, lineOffsets []int) {
	PT.nodes.InsertBefore(&Node{typ: typ, start: start, length: length, lineOffsets: lineOffsets}, PT.nodes.Back())
}

func (PT *PieceTable) newNodeBefore(typ NodeType, start, length int, lineOffsets []int, currentNode *list.Element) {
	abc := currentNode.Next()
	if abc != nil {
		PT.nodes.InsertBefore(&Node{typ: typ, start: start, length: length, lineOffsets: lineOffsets}, currentNode.Next())
	}

}

func (PT *PieceTable) newNodeAfter(typ NodeType, start, length int, lineOffsets []int, currentNode *list.Element) {
	abc := currentNode.Next()
	if abc != nil {
		PT.nodes.InsertAfter(&Node{typ: typ, start: start, length: length, lineOffsets: lineOffsets}, currentNode.Prev())
	}
	// PT.nodes.InsertAfter(&Node{typ: typ, start: start, length: length, lineOffsets: lineOffsets}, currentNode.Next())
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

//AppendString allows a new string to be added to the add buffer
// this is strictly for append
func (PT *PieceTable) AppendString(data string) {
	addBufBeforeLen := len(PT.buffer[Added])
	d := []rune(data)
	dLen := len(d)
	los := getLineOffsets(d)
	PT.buffer[Added] = append(PT.buffer[Added], d...)
	PT.newNodeAppendOnly(Added, addBufBeforeLen, dLen, los)
}

//InsertStringAt will insert a string into the piece table at an offset
// this offset can be considered the byte location from the beginning of
// the visible buffers.
// data is the string to (append to the add buffer) be added to the
// PieceTable
func (PT *PieceTable) InsertStringAt(offset int, data string) {
	//record where we are in the document pretty much
	totLen := 0
	// -------------------------------------------------------------
	// This is the new node that we are adding to the `Added` buffer
	// because this new Node in the dll is getting inserted at an offset
	// that is what actually determines its place in the dll and not
	// the content that refers to the view on the buffer
	// -------------------------------------------------------------
	newNodeTyp := Added
	// calulate line offsets for the newly inserted data
	newNodeLos := getLineOffsets([]rune(data))
	//the start of the new node is the current length of the `Added`
	// buffer because that is where this new data will be visible from.
	newNodeStart := len(PT.buffer[Added])
	// this is the length of the data we are passing in
	newNodeLength := len(data)
	// append the rest of the string to the add buffer
	PT.buffer[Added] = append(PT.buffer[Added], []rune(data)...)
	newNodeMiddle := &Node{typ: newNodeTyp, start: newNodeStart, length: newNodeLength, lineOffsets: newNodeLos}

	// currentTotLen := 0
	// looop through the nodes and find out where the offset is gonna be, use the length += next length to
	for e := PT.nodes.Front(); e != nil; e = e.Next() {
		n, ok := e.Value.(*Node)
		if !ok {
			panic("Not unrwapping a node")
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
			newNodeLeft := &Node{typ: nodeLeftTyp, start: nodeLeftStart, length: nodeLeftLength, lineOffsets: nodeLeftLos}


			nodeMiddleTyp := Added
			nodeMiddleStart := len(PT.buffer[Added])
			nodeMiddleLength := len(data)
			nodeMiddleLos := getLineOffsets([]rune(data))
			newNodeMiddle := &Node{typ: Added, start: nodeMiddleStart, length: nodeMiddleLength, lineOffsets: nodeMiddleLos}


			//new node for the right
			nodeRightTyp := n.typ
			nodeRightStart:= nodeLeftLength
			nodeRightLenth := totLen - offset
			nodeRightLos := getLineOffsets(PT.buffer[n.typ][nodeRightStart:nodeRightLenth+nodeRightStart])
			newNodeRight := &Node{typ: nodeRightTyp, start: nodeRightStart, length: nodeRightLenth, lineOffsets: nodeRightLos}
		} else if offset == totLen {
			//insert between 2 dll nodes
		} 

		if offset >= n.start && offset <= n.start+n.length {
			
			// now create the right node which is trickier
			nodeRightTyp := n.typ
			nodeRightStart := offset
			nodeRightLength := 
		}

	}
}

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
		pt.newNodeAppendOnly(Original, 0, optBufLen, los)
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
			println("n.start=", n.start)
			println("n.length=", n.length)
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

	// data := openAndReadFile("peace.go")

	input := []rune(`The quick 
	brown`)
	println("len(input)=", len(input))
	pt := NewPT(input)
	pt.AppendString(`//EXTRA
	asfjk

	// data to have at the bottom test`)

	pt.InsertStringAt(6, `Here is the new 
data`)

	// pt.InsertStringAt(28, `Here is the new afjslkjasflkjasflk
	// afskjfaskasfljfa
	// asfjasfkjfkasj
	// data`)
	for e := pt.nodes.Front(); e != nil; e = e.Next() {
		n := e.Value.(*Node)
		fmt.Println(n)
	}
	cat(pt)
	println()
	// x.AppendBytes([]byte("More Text Here:"))
	// x.AppendBytes([]byte("\n\n"))
	// x.AppendBytes([]byte("\tMore Text Over Here"))
	// x.Display()

}
