package main

// an attempt at making a piecetable/piecemap in go but right now im using list(doubly) from container/list
// and honestly a bunch of other random stuff but im ready to start recording my eventual implementation

import (
	"container/list"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
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
	original []byte
	added    []byte
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

//AppendBytes allows append only nodes to be added to the piece table
func (PT *PieceTable) AppendBytes(data []byte) {
	dataLen := len(data)
	dataStart := len(PT.added)

	PT.added = append(PT.added, data...)
	//calculate line offsets
	los := getLineOffsets(data)
	PT.newNode(Added, dataStart, dataLen, true, los)
}

// Display currently displays the []bytes to the terminal ( there will be read functions instead)
func (PT *PieceTable) Display() {
	for e := PT.nodes.Front(); e != nil; e = e.Next() {
		n := e.Value.(*Node)
		if n.typ == Original && n.visible {
			for i := n.start; i < n.start+n.length; i++ {
				fmt.Print(string(PT.original[i]))
			}
		} else if n.typ == Added && n.visible {
			for i := n.start; i < n.start+n.length; i++ {
				fmt.Print(string(PT.added[i]))
			}
		}
	}
}
func (PT *PieceTable) printNodes() {
	for e := PT.nodes.Front(); e != nil; e = e.Next() {
		n := e.Value.(*Node)
		fmt.Println(n)
	}
}

func getLineOffsets(buf []byte) []int {
	var bucket []int
	for i := 0; i < len(buf); i++ {
		if buf[i] == '\n' {
			bucket = append(bucket, i)
		}
	}
	return bucket
}

// NewPT will eventually return a piecetable/map and will probably have a separate
// new function for the optional buffer (this would be starting a new buffer for instance)
func NewPT(optBuf []byte) *PieceTable {
	if optBuf != nil && len(optBuf) < (32*1024) {
		pt := &PieceTable{original: optBuf, added: []byte(""), nodes: list.New()}
		//calculate lineoffsets
		los := getLineOffsets(optBuf)
		pt.newNode(Original, 0, len(optBuf), true, los)
		return pt
	}
	return &PieceTable{original: []byte(""), added: []byte(""), nodes: list.New()}
}

func openAndReadFile(f string) ([]byte, error) {
	var fd *os.File
	if _, err := os.Stat(f); err == nil {

		fd, err = os.Open(f)
		// eventually include this as part of a shutdown
		defer fd.Close()
		if err != nil {
			// TODO Handle failing to open file
			// need to figure out how i will display that to the user
			// or make a generic die function
			return nil, err
		}
		fi, err := fd.Stat()
		if err != nil {
			return nil, err

		}
		//1 MB?
		if fi.Size() > (1 * 1024 * 1024) {
			//do a buffered read
			buf := make([]byte, 32*1024) // define your buffer size here.

			for {
				n, err := fd.Read(buf)

				if n > 0 {
					return buf[:n], nil // your read buffer.
				}

				if err == io.EOF {
					break
				}
				if err != nil {
					log.Printf("read %d bytes: %v", n, err)
					break
				}
			}
		} else {
			//otherwise read all with ioutil
			data, err := ioutil.ReadAll(fd)
			if err != nil {
				return nil, err
			}
			return data, nil
		}
	}
	return nil, fmt.Errorf("Need to call with a file that exists for now")

}

func main() {

	data, err := openAndReadFile("test.go")
	if err != nil {
		fmt.Println(err)
	}

	x := NewPT(data)
	x.AppendBytes([]byte("More Text Here:"))
	x.AppendBytes([]byte("\n\n"))
	x.AppendBytes([]byte("\tMore Text Over Here"))
	x.Display()
	x.DeleteBytesAt(1, 10, []byte("faslj"))

}
