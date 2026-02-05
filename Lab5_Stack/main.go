package main

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"slices"
)

const MAX_VALUE = 100

type Stack struct {
    elements []int
}


func display(stack *Stack) {
    for i := range stack.elements {
        fmt.Printf("%d ", stack.elements[i])
    }
    fmt.Printf("\n")
}


func generateRandNumber() int {
    max := big.NewInt(MAX_VALUE)

    tmp, err := rand.Int(rand.Reader, max)
    if err != nil {
        panic(err)
    }

    return int(tmp.Int64())
}


func (stack *Stack) getLen() int {
    return len(stack.elements)
}


func isExists(stack *Stack, searchKey int) bool {
	return slices.Contains(stack.elements, searchKey);
}


func (stack *Stack) pop() int {
	if len(stack.elements) < 1 {
		return 0;
	}

	topElement := stack.elements[len(stack.elements)-1];
	stack.elements = stack.elements[:len(stack.elements)-1];

	return topElement;
}


func (stack *Stack) push(value int) {
    stack.elements = append(stack.elements, value)
}


func main() {
    stack := Stack{}

    for i := 0; i < 10; i++ {
        stack.push(generateRandNumber())
    }

    display(&stack)

    length := stack.getLen()

    fmt.Println("Stack len =", length)

    for length > 0 {
        fmt.Println("Popped ", stack.pop())
        display(&stack)
        length--
    }
}
