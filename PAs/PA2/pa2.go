package main

import "fmt"

const MARKER = 'a'

type TrieNode struct {
	children   map[byte]*TrieNode
	endOfWords bool
}

func createNode() *TrieNode {
	node := &TrieNode{}
	node.endOfWords = false

	node.children = make(map[byte]*TrieNode)

	return node
}

func display(root *TrieNode) {
	displayPreOrder(root, "");
}

func displayPreOrder(node *TrieNode, str string) {
	if node.endOfWords == true {
		fmt.Println(str)
	}

	for char, childNode := range node.children {
		displayPreOrder(childNode, str + string(byte(char)))
	}
}

func getIndex(ch byte) int {
	return int(ch - MARKER)
}

func insert(root *TrieNode, key string) {
	tmp := root

	for i := 0; i < len(key); i++ {
		if tmp.children[key[i]] == nil {
			tmp.children[key[i]] = createNode();
		}
		tmp = tmp.children[key[i]]
	}

	tmp.endOfWords = true
}

/**
	Structure:
	root {
		children: map {
			'a': node {
				children: map {
					'n': node {
						children: map {
							't': node {
								children: map,
								endOfWords: true
							}
						}
					}
				}
			}
		}
		endOfWords: false
	}
	
	Would delete the 't' nodes if key is "ant" and key 't' node is an end of word
	and would move down until finding endOfWord = true. and then delete/move up.
*/
func remove(root *TrieNode, key string) {
	if (root == nil) {
		return;
	}

	tmp := root;
	// will be keeping track of path down,
	// so can delete during move up
	wordPath := make([]*TrieNode, len(key)+1);
	wordPath[0] = tmp;

	for i := 0; i < len(key); i++ {
		if (tmp.children[key[i]] == nil) {
			return;
		}

		wordPath[i + 1] = tmp.children[key[i]];
		tmp = tmp.children[key[i]];
	}

	// Now tmp is at the node with true endOfWords
	tmp.endOfWords = false;

	for i := len(key); i > 0; i-- {
		if(len(wordPath[i].children) == 0 && !wordPath[i].endOfWords){
			delete(wordPath[i - 1].children, key[i - 1]);
		} else {
			return;
		}
	}
}

func search(root *TrieNode, key string) bool {
	tmp := root

	for i := 0; i < len(key); i++ {
		if tmp.children[key[i]] != nil {
			tmp = tmp.children[key[i]]
		} else {
			return false
		}
	}

	return (tmp != nil && tmp.endOfWords)
}

func main() {
	/**
		Structure:
		root {
			children: map {
				'a': node {
					children: map {
						'n': node {
							children: map {
								't': node {
									children: map,
									endOfWords: true
								}
							}
						}
					}
				}
			}
			endOfWords: false
		}
	*/
	words := []string{"cat", "bee", "apple", "ant", "antler", "beewax", "car"}
	root := createNode();

	for i := 0; i < len(words); i++ {
		insert(root, words[i])
	}

	fmt.Println("contains the word ant", search(root, "ant"))
	fmt.Println("contains the word antler", search(root, "antler"))
	fmt.Println("contains the word a", search(root, "a"))
	fmt.Println("contains the word bee", search(root, "bee"))
	fmt.Println("contains the word be", search(root, "be"))
	fmt.Println("contains the word car", search(root, "car"))
	fmt.Println("contains the word cart", search(root, "cart"))

	// Added antler for testing
	// should not delete any nodes from map since t has children
	remove(root, "ant");
	fmt.Println("\nTest 1:");
	fmt.Println("contains the word ant", search(root, "ant"));
	fmt.Println("contains the word antler", search(root, "antler"));

	remove(root, "cart");
	fmt.Println("\nTest 2:");
	fmt.Println("contains the word cart", search(root, "cart"));
	fmt.Println("contains the word car", search(root, "car"));
}