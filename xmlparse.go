package main

import (
	"encoding/json"
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"reflect"
	"strings"
)

var url = ""
var file = ""

type Node struct {
	name     string
	text     string
	attrs    map[string]string
	children []*Node
}

func init() {
	flag.StringVar(&url, "url", "", "url to an xml file")
	flag.StringVar(&file, "file", "", "name of an xml file")
	flag.Parse()
}

func main() {
	var r *Node
	if url == "" && file == "" {
		fmt.Println("no url of file provided")
		return
	}
	if url != "" {
		resp, err := http.Get(url)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()
		reader := resp.Body
		r = XMLDecode(reader)
	} else {
		f, err := os.Open(file)
		if err != nil {
			log.Fatal(err)
		}
		r = XMLDecode(f)
	}
	dumpnodes(r, 0)
	//var ra []*Node
	//ra = append(ra, r)
	//dumpjson(ra)
}

func dumpnodes(n *Node, x int) {
	printnode(n, x)
	for i := range n.children {
		c := n.children[i]
		dumpnodes(c, x+1)
	}
}

func dumpjson(n []*Node) {
	data, err := json.Marshal(n)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", data)
}

func vtype(v any) any {
	vt := reflect.TypeOf(v)
	fmt.Println("	TypeOf", vt)
	return vt
}

func makeindent(depth int) string {
	return strings.Repeat("	", depth)
}

func printnode(n *Node, depth int) {
	var indent = makeindent(depth)
	fmt.Println(indent, "depth", depth)
	fmt.Println(indent, "Name", n.name)
	fmt.Println(indent, "Text", n.text)
	for k, v := range n.attrs {
		fmt.Println(indent, "key", k, "val", string(v) )
	}
	fmt.Println(indent, "N children", len(n.children) )
        fmt.Println("")
}
func newnode(tok xml.StartElement) *Node {
	var n *Node
	n = new(Node)
	n.attrs = make(map[string]string, 0)
	n.children = make([]*Node, 0)

	n.name = tok.Name.Local
	for _, a := range tok.Attr {
		n.attrs[a.Name.Local] = a.Value
	}
	return n
}

func XMLDecode(rc io.Reader) *Node {

	var root, node *Node
	var Nodestack []*Node // stack of element names

	dec := xml.NewDecoder(rc)

	for {
		tok, err := dec.Token()
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal(err)
		}
		switch tok := tok.(type) {
		case xml.StartElement:

			node = newnode(tok)
			if len(Nodestack) == 0 {
				root = node
			} else {
				var pnode *Node
				pnode = Nodestack[len(Nodestack)-1]
				pnode.children = append(pnode.children, node)
			}
			Nodestack = append(Nodestack, node)

		case xml.EndElement:

			Nodestack = Nodestack[:len(Nodestack)-1] // pop

		case xml.CharData:
			node = Nodestack[len(Nodestack)-1]
			node.text = string(tok)
		case xml.Comment:
		case xml.ProcInst:
		case xml.Directive:
		}
	}
	return root
}

//!-
