package main

import (
	"fmt"
	"log"
	"strings"

	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Tree map[interface{}]interface{}

func (t Tree) printLeaves() {
	for k, v := range t {
		switch node := v.(type) {
		case string:
			if strings.HasPrefix(node, "secret ") {
				fmt.Println(k.(string) + "=" + strings.TrimPrefix(node, "secret "))
			}
		case Tree:
			node.printLeaves()
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

func (t Tree) removeEmptyBranches() Tree {
	for k, v := range t {
		switch node := v.(type) {
		case Tree:
			if len(node) == 0 {
				delete(t, k)
			} else {
				t[k] = node.removeEmptyBranches()
			}
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

	t.printLeaves()
	t = t.removeSecrets()
	t = t.removeEmptyBranches()

	output, err := yaml.Marshal(t)
	if err != nil {
		log.Fatalf("error marshaling output: %v", err)
	}
	fmt.Println("\n\nTree as YAML:\n", string(output))
}

