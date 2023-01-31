package main

import (
	"fmt"
	"log"

	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Tree map[interface{}]interface{}

func (t Tree) leaves() []interface{} {
	var leaves []interface{}
	for k, v := range t {
		switch node := v.(type) {
		case string:
			if strings.HasPrefix(node, "secret ") {
				leaves = append(leaves, strings.TrimPrefix(node, "secret "))
			}
		case Tree:
			leaves = append(leaves, node.leaves()...)
		}
	}
	return leaves
}

func main() {
	file, err := ioutil.ReadFile("tree.yaml")
	if err != nil {
		log.Fatalf("error reading file: %v", err)
	}

	var t Tree
	err = yaml.Unmarshal(file, &t)
	if err != nil {
		log.Fatalf("error unmarshaling file: %v", err)
	}

	for _, leaf := range t.leaves() {
		fmt.Println(leaf)
	}
}

