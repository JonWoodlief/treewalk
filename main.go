package main

import (
	"fmt"
	"log"
	"strings"

	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Tree map[interface{}]interface{}

func (t Tree) printLeaves(prefix string) {
	for k, v := range t {
		switch node := v.(type) {
		case string:
			if !strings.HasPrefix(node, "secret ") {
				fmt.Println(prefix+k.(string), "=", node)
			}
		case Tree:
			node.printLeaves(prefix + k.(string) + ".")
		}
	}
}

func (t Tree) removeSecrets() Tree {
	for k, v := range t {
		switch node := v.(type) {
		case string:
			if strings.HasPrefix(node, "secret ") {
				delete(t, k)
			}
		case Tree:
			t[k] = node.removeSecrets()
		}
	}
	return t
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

	t = t.removeSecrets()

	t.printLeaves("")

	output, err := yaml.Marshal(t)
	if err != nil {
		log.Fatalf("error marshaling output: %v", err)
	}
	fmt.Println("\n\nTree as YAML:\n", string(output))
}
