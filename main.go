package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"strconv"
	"strings"
)

func getLeaves(node interface{}, path []string) [][2]interface{} {
	var leaves [][2]interface{}
	switch n := node.(type) {
	case map[interface{}]interface{}:
		for key, value := range n {
			keyString, ok := key.(string)
			if !ok {
				keyString = strconv.Itoa(key.(int))
			}
			leaves = append(leaves, getLeaves(value, append(path, keyString))...)
		}
	case []interface{}:
		for i, item := range n {
			leaves = append(leaves, getLeaves(item, append(path, strconv.Itoa(i)))...)
		}
	default:
		leaves = append(leaves, [2]interface{}{strings.Join(path, "."), n})
	}
	return leaves
}

func main() {
	file, err := ioutil.ReadFile("tree.yaml")
	if err != nil {
		panic(err)
	}

	var tree interface{}
	err = yaml.Unmarshal(file, &tree)
	if err != nil {
		panic(err)
	}

	leaves := getLeaves(tree, []string{})
	for _, leaf := range leaves {
		fmt.Printf("%v=%v\n", leaf[0], leaf[1])
	}
}

