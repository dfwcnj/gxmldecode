
# gxmldecode
#
#

This is a work in progress

xml unmarshal doesn't seem to do what i would like for it to do. for
example, the following doesn't work. of course, it could be that i have
a fundamental misunderstanding of the language.

package main

import (
	"encoding/json"
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
)

var file = ""

type Node struct {
	Name     string
	Text     string
	Attrs    map[string]string
	Children []*Node
}

func init() {
	flag.StringVar(&file, "file", "", "name of an xml file")
	flag.Parse()
}

func main() {
	if file == "" {
		fmt.Println("no file provided")
		return
	}
	f, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}
	unmarshal(f)
}

func unmarshal(rc io.Reader) {
	b, err := io.ReadAll(rc)
	if err != nil {
		log.Fatal(err)
	}
	node := new(Node)
	err = xml.Unmarshal(b, node)
	if err != nil {
		log.Fatal(err)
	}
	data, err := json.MarshalIndent(node, "", "	")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", data)
}

This little program parses either a url of an xml file or a local xml
file to a hierarchical structure. The structure 'looks' like the
structure returned by the go decoder because that is what the program
uشes to extract the structure of the xml bytes. It was inspired by the python
elementtree package. o

I am trying to be a good citizen by making the node structure
'unexported' while making its fields exported so that json.MarshalIndent
is able to produce something that makes senѕe.



