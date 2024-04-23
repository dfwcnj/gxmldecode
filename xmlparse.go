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

type node struct {
	Name     string				`json:"name"`
	Text     string				`json:"text"`
	Attrs    map[string]string		`json:"attrs"`
	Children []*node			`json:"children"`
}

func init() {
	flag.StringVar(&url, "url", "", "url to an xml file")
	flag.StringVar(&file, "file", "", "name of an xml file")
	flag.Parse()
}

func main() {
	var r *node
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
	//dumpnodes(r, 0)
	dumpjson(r)
}

func dumpnodes(n *node, x int) {
	printnode(n, x)
	for i := range n.Children {
		c := n.Children[i]
		dumpnodes(c, x+1)
	}
}

func dumpjson(n *node) {
	data, err := json.MarshalIndent(n, "", "	")
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

func printnode(n *node, depth int) {
	var indent = makeindent(depth)
	fmt.Println(indent, "Depth:", depth)
	fmt.Println(indent, "Name:", n.Name)
	fmt.Println(indent, "Text:", n.Text)
	for k, v := range n.Attrs {
		fmt.Println(indent, "Key:", k, "Val:", string(v))
	}
	fmt.Println(indent, "Children:", len(n.Children))
	fmt.Println("")
}
func newnode(tok xml.StartElement) *node {
	var n *node
	n = new(node)
	n.Attrs = make(map[string]string, 0)
	n.Children = make([]*node, 0)

	n.Name = tok.Name.Local
	for _, a := range tok.Attr {
		n.Attrs[a.Name.Local] = a.Value
	}
	return n
}

func XMLDecode(rc io.Reader) *node {

	var root, nod *node
	var Nodestack []*node // stack of element names

	dec := xml.NewDecoder(rc)
	dec.Strict = false

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

			nod = newnode(tok)
			if len(Nodestack) == 0 {
				root = nod
			} else {
				var pnode *node
				pnode = Nodestack[len(Nodestack)-1]
				pnode.Children = append(pnode.Children, nod)
			}
			Nodestack = append(Nodestack, nod)

		case xml.EndElement:

			Nodestack = Nodestack[:len(Nodestack)-1] // pop

		case xml.CharData:
			if len(Nodestack) > 0 {
				nod = Nodestack[len(Nodestack)-1]
				nod.Text = string(tok)
			} else {
                            //fmt.Println("bare CharData", string(tok) )
                        }
		case xml.Comment:
		case xml.ProcInst:
		case xml.Directive:
		}
	}
	return root
}

//!-
