package main

// an attempt at making a piecetable/piecemap in go but right now im using list(doubly) from container/list
// and honestly a bunch of other random stuff but im ready to start recording my eventual implementation

import (
	"container/list"
	"fmt"
	"io/ioutil"
	"log"
)

// NodeType Enum they are both readonly/appendonly buffers
type NodeType int

const (
	//Added Buffer NodeType descriptor
	Added NodeType = iota
	//Original Buffer NodeType descriptor
	Original
)

// PieceTable is currently 2 buffers but will be modified in the future
type PieceTable struct {
	original []rune
	added    []rune
	nodes    *list.List
}

func (n *Node) String() (result string) {
	if n.typ == Original {
		result = fmt.Sprintf("{NodeType: Original, start: %d, length: %d, visible: %v, lineOffsets: %v}",
			n.start, n.length, n.visible, n.lineOffsets)
	} else {
		result = fmt.Sprintf("{NodeType: Added, start: %d, length: %d, visible: %v, lineOffsets: %v}}",
			n.start, n.length, n.visible, n.lineOffsets)
	}
	return
}

// Node is the element in the list that contains some metadata for the contents and the operation
type Node struct {
	typ         NodeType
	start       int
	length      int
	visible     bool
	lineOffsets []int
}

func (PT *PieceTable) newNode(typ NodeType, start, length int, visible bool, lineOffsets []int) {
	PT.nodes.PushBack(&Node{typ: typ, start: start, length: length, visible: visible, lineOffsets: lineOffsets})
}

// //AppendBytes allows append only nodes to be added to the piece table
// func (PT *PieceTable) AppendBytes(data []byte) {
// 	dataLen := len(data)
// 	dataStart := len(PT.added)

// 	PT.added = append(PT.added, data...)
// 	//calculate line offsets
// 	los := getLineOffsets(data)
// 	PT.newNode(Added, dataStart, dataLen, true, los)
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
	PT.newNode(Added, addBufBeforeLen, dLen, true, los)
}

// NewPT will eventually return a piecetable/map and will probably have a separate
// new function for the optional buffer (this would be starting a new buffer for instance)
func NewPT(optBuf []rune) *PieceTable {
	if optBuf != nil {
		optBufLen := len(optBuf)
		pt := &PieceTable{original: optBuf, added: []rune(""), nodes: list.New()}
		//calculate lineoffsets
		los := getLineOffsets(optBuf)
		pt.newNode(Original, 0, optBufLen, true, los)
		return pt
	}
	return &PieceTable{original: []rune(""), added: []rune(""), nodes: list.New()}
}

func openAndReadFile(f string) []rune {
	data, err := ioutil.ReadFile(f)
	if err != nil {
		log.Fatal(err)
	}
	return []rune(string(data))
}

func main() {

	data := openAndReadFile("peace.go")

	pt := NewPT(data)
	pt.AppendString(`//EXTRA
	data to have at the bottom test`)
	for e := pt.nodes.Front(); e != nil; e = e.Next() {
		n := e.Value.(*Node)
		fmt.Println(n)
	}
	// x.AppendBytes([]byte("More Text Here:"))
	// x.AppendBytes([]byte("\n\n"))
	// x.AppendBytes([]byte("\tMore Text Over Here"))
	// x.Display()

}
