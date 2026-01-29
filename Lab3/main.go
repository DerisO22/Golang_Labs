package main

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"
)

// so 0-10
const MAX_NUM = 10;
const OPERATORS = "+-*";

// generate rand operator
func generateRandOperator() byte {
	// Geeks For Geeks (rune/char/byte count in string):
	randNum := generateRandNumber(len(OPERATORS));

	return OPERATORS[randNum];
}

// generate rand num
func generateRandNumber(max_num int) int {
	// from lab 3 specifications
	max := big.NewInt(int64(max_num));

	tmp, err := rand.Int(rand.Reader, max);

	if err != nil {
		panic(err);
	}

	return int(tmp.Int64());
}

// check answer
func getRealAnswer(num1 int, num2 int, operator byte) int {
	switch operator {
	case '+':
		return num1 + num2;
	case '-':
		return num1 - num2;
	case '*':
		return num1 * num2;
	default:
		return 0;
	}
}

func main() {
	keepPlaying := "y";
	var studentAnswer int;

	fmt.Println("Welcome to the Math Test Game\n");

	// loop for math game
	for (strings.ToLower(keepPlaying)) != "n" {
		num1, num2 := generateRandNumber(MAX_NUM) + 1, generateRandNumber(MAX_NUM) + 1;
		operator := generateRandOperator();

		// helpful chart I found for printf in go: https://www.geeksforgeeks.org/go-language/fmt-printf-function-in-golang-with-examples/
		fmt.Printf("\nWhat is %d %c %d?", num1, operator, num2);
		
		fmt.Print("\n\nEnter your answer: ");
		fmt.Scan(&studentAnswer);

		// check answer
		realAnswer := getRealAnswer(num1, num2, operator);

		if studentAnswer == realAnswer {
			fmt.Println("\nNice Job, You Are Correct\n")
		} else {
			fmt.Printf("Incorrect Answer, The Correct Answer is %d", realAnswer);
		}

		// check if want to keep playing
		fmt.Print("\nWould You Like to Continue(y/n)?: ");
		fmt.Scan(&keepPlaying);
	}
}