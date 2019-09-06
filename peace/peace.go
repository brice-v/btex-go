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
	original []rune
	added    []rune
	nodes    *list.List
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

// //AppendBytes allows append only nodes to be added to the piece table
// func (PT *PieceTable) AppendBytes(data []byte) {
// 	dataLen := len(data)
// 	dataStart := len(PT.added)

// 	PT.added = append(PT.added, data...)
// 	//calculate line offsets
// 	los := getLineOffsets(data)
// 	PT.newNodeAppendOnly(Added, dataStart, dataLen, true, los)
// }

// // Display currently displays the []bytes to the terminal ( there will be read functions instead)
// func (PT *PieceTable) Display() {
// 	for e := PT.nodes.Front(); e != nil; e = e.Next() {
// 		n := e.Value.(*Node)
// 		if n.typ == Original && n.visible {
// 			for i := n.start; i < n.start+n.length; i++ {
// 				fmt.Print(string(PT.original[i]))
// 			}
// 		} else if n.typ == Added && n.visible {
// 			for i := n.start; i < n.start+n.length; i++ {
// 				fmt.Print(string(PT.added[i]))
// 			}
// 		}
// 	}
// }
// func (PT *PieceTable) printNodes() {
// 	for e := PT.nodes.Front(); e != nil; e = e.Next() {
// 		n := e.Value.(*Node)
// 		fmt.Println(n)
// 	}
// }

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
	addBufBeforeLen := len(PT.added)
	d := []rune(data)
	dLen := len(d)
	los := getLineOffsets(d)
	PT.added = append(PT.added, d...)
	PT.newNodeAppendOnly(Added, addBufBeforeLen, dLen, los)
}

//InsertStringAt will insert a string into the piece table at an offset
// this offset can be considered the byte location from the beginning of
// the visible buffers.
// data is the string to (append to the add buffer) be added to the
// PieceTable
func (PT *PieceTable) InsertStringAt(offset int, data string) {
	newNodeStart := len(PT.added)
	newNodeLength := len(data)
	PT.added = append(PT.added, []rune(data)...)
	los := getLineOffsets([]rune(data))

	totLen := 0
	// looop through the nodes and find out where the offset is gonna be, use the length += next length to
	for e := PT.nodes.Front(); e != nil; e = e.Next() {
		n, ok := e.Value.(*Node)
		if !ok {
			panic("Not unrwapping a node")
		}
		//Skip the sentinel nodes
		if n.typ == Sentinel {
			if e.Next() == nil {
				//this means were in the end node
			}
			continue
		}
		totLen = totLen + n.length

		if offset < totLen {
			// this is all for the node that goes before
			currentNodeType := n.typ
			currentStart := n.start
			lengthToOffset := n.start + offset
			//lets just imagine  the above will work

			//not sure if any of these ranges make sense and if thats the cause
			var recalculatedLineOffsets []int
			if currentNodeType == Original {
				newBuf := PT.original[currentStart:lengthToOffset]
				println(string(newBuf))
				recalculatedLineOffsets = getLineOffsets(newBuf)
			} else {
				newBuf := PT.added[currentStart:lengthToOffset]
				recalculatedLineOffsets = getLineOffsets(newBuf)
			}

			//this is the original data and were fixing the view on it
			PT.newNodeBefore(currentNodeType, currentStart, lengthToOffset, recalculatedLineOffsets, e.Prev())

			// this is the new data insertion
			PT.newNodeBefore(Added, newNodeStart, newNodeLength, los, e.Prev())

			//fixing the new view continued
			newStart := n.start + (offset - n.start)

			if currentNodeType == Original {
				newBuf := PT.original[newStart:n.length]
				recalculatedLineOffsets = getLineOffsets(newBuf)
			} else if currentNodeType == Added {
				newBuf := PT.added[newStart:n.length]
				recalculatedLineOffsets = getLineOffsets(newBuf)
			} else {
				continue
			}
			PT.newNodeAfter(currentNodeType, newStart, n.length, recalculatedLineOffsets, e)
			// dont know if its possible but delete the node were standing on
			abc := e.Next()
			if abc != nil {
				PT.nodes.Remove(e.Prev())
			}

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
		pt := &PieceTable{original: optBuf, added: []rune(""), nodes: newEmptyList()}
		//calculate lineoffsets
		los := getLineOffsets(optBuf)
		pt.newNodeAppendOnly(Original, 0, optBufLen, los)
		return pt
	}
	return &PieceTable{original: []rune(""), added: []rune(""), nodes: newEmptyList()}
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
			fmt.Print(string(pt.original[n.start : n.start+n.length]))
		} else if n.typ == Added {
			fmt.Print(string(pt.added[n.start : n.start+n.length]))
		} else {
			// e = e.Next()
			continue
		}
	}
}

func main() {

	// data := openAndReadFile("peace.go")

	pt := NewPT([]rune(`The quick 
	brown`))
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
