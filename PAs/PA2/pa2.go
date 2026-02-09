package main

import "fmt"

const SIZE = 26
const MARKER = 'a'

type TrieNode struct {
	children   [SIZE]*TrieNode
	endOfWords bool
}

func createNode() *TrieNode {
	node := &TrieNode{}
	node.endOfWords = false

	for i := 0; i < SIZE; i++ {
		node.children[i] = nil
	}

	return node
}

func display(wordsMap map[string]TrieNode) {
	displayPreOrder(wordsMap, "")
}

func displayPreOrder(wordsMap map[string]TrieNode, str string) {
	if wordsMap[str].endOfWords == true {
		fmt.Println(str)
	}

	for i := 0; i < SIZE; i++ {
		if wordsMap[str].children[i] != nil {
			displayPreOrder(wordsMap, str + string(byte(i + MARKER)))
		}
	}
}

func getIndex(ch byte) int {
	return int(ch - MARKER)
}

func insert(root *TrieNode, key string) {
	tmp := root

	for i := 0; i < len(key); i++ {
		index := getIndex(key[i])
		if tmp.children[index] == nil {
			tmp.children[index] = createNode()
		}
		tmp = tmp.children[index]
	}

	tmp.endOfWords = true
}

func search(root *TrieNode, key string) bool {
	tmp := root

	for i := 0; i < len(key); i++ {
		index := getIndex(key[i])
		if tmp.children[index] != nil {
			tmp = tmp.children[index]
		} else {
			return false
		}
	}

	return (tmp != nil && tmp.endOfWords)
}

func main() {
	wordsTest := make(map[string]TrieNode);

	/**
		Structure:
		{
			"cat": {children TrieNode, endOfWords bool},
			"bee": {children TrieNode, endOfWords bool},
			"apple": {children TrieNode, endOfWords bool}
		}
	*/
	words := []string{"cat", "bee", "apple", "ant", "beewax", "car"}

	root := createNode();

	for _, word := range words {
		wordsTest[word] = TrieNode{
			root.children,
			false,
		}

		insert(root, word);
	}

	fmt.Println(wordsTest);

	display(wordsTest);

	fmt.Println("contains the word ant", search(root, "ant"))
	fmt.Println("contains the word a", search(root, "a"))
	fmt.Println("contains the word bee", search(root, "bee"))
	fmt.Println("contains the word be", search(root, "be"))
	fmt.Println("contains the word car", search(root, "car"))
	fmt.Println("contains the word cart", search(root, "cart"))
}