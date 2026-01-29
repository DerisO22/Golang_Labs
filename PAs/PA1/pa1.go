package main

import (
	"bufio"
	"crypto/rand"
	"fmt"
	"math/big"
	"os"
	"time"
)

const FILE_TO_READ = "messages.txt";
const DATE_LAYOUT = "01/02/2006";

// reads and add to list
func readFromFile(messageArray *[]string, messageArrayLength *int) {
	// From Lab 3 Example of Reading File
	readFile, err := os.Open(FILE_TO_READ);
	
	if err != nil {
		fmt.Println("Error opening file:", err);
		return;
	}
	
	defer readFile.Close();

	scanner := bufio.NewScanner(readFile);
	fmt.Println("Reading from file:");

	for scanner.Scan() {
		*messageArray = append(*messageArray, scanner.Text());
		*messageArrayLength++;
	}
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

func calculateAge(inputBirthDate time.Time) int {
	currentTime := time.Now();

	// Age in years
	age := currentTime.Year() - inputBirthDate.Year();

	// "For example, if today is 1/20/2026 and the user is born on 1/21/2006, 
	// then the user is still 19 until the next day. On the other hand, if the 
	// user is born on 1/19/2006, the user is 20".
	// Also leap years
	if currentTime.Month() < inputBirthDate.Month() || (inputBirthDate.Month() == currentTime.Month() && currentTime.Day() < currentTime.Day()){
		age--;
	}

	return age;
}

func getZodiacSign(inputBirthDate time.Time) string {
	// All the ranges and zodiac signs from Appendix A:
	/*  January 21 February 19 Aquarius
		February 20 March 20 Pisces
		March 21 April 20 Aries
		April 21 May 21 Taurus
		May 22 June 21 Gemini
		June 22 July 22 Cancer
		July 23 August 21 Leo
		August 22 September 23 Virgo
		September 24 October 23 Libra
		October 24 November 22 Scorpio
		November 23 December 21 Sagittarius
		December 22 January 20 Capricorn
	*/	
	day := inputBirthDate.Day();
	month := inputBirthDate.Month();

	switch {
		case (month == 1 && day >= 21) || (month == 2 && day <= 19):
			return "Aquaris";
		case (month == 2 && day >= 20) || (month == 3 && day <= 20):
			return "Pisces";
		case (month == 3 && day >= 21) || (month == 4 && day <= 20):
			return "Aries";
		case (month == 4 && day >= 21) || (month == 5 && day <= 21):
			return "Taurus";
		case (month == 5 && day >= 22) || (month == 6 && day <= 21):
			return "Gemini";
		case (month == 6 && day >= 22) || (month == 7 && day <= 22):
			return "Cancer";
		case (month == 7 && day >= 23) || (month == 8 && day <= 21):
			return "Leo";
		case (month == 8 && day >= 22) || (month == 9 && day <= 23):
			return "Virgo";
		case (month == 9 && day >= 24) || (month == 10 && day <= 23):
			return "Libra";
		case (month == 10 && day >= 24) || (month == 11 && day <= 22):
			return "Scorpio";
		case (month == 11 && day >= 23) || (month == 12 && day <= 21):
			return "Sagittarius";
		case (month == 12 && day >= 22) || (month == 1 && day <= 20):
			return "Capricorn";
		default:
			return "";
	}
}

func main() {
	// using slice like the generateArray func from lab3 example
	var messages []string;
	var name string;
	var birthDate string;
	messageArrayLength := 0;

	// I could just return the new populated message slice
	// in readFromFile, but wanted to try out pointers/pass-by-ref
	readFromFile(&messages, &messageArrayLength);

	fmt.Println("\nI'm am the Almighty Psychic Fortune Teller!");
	fmt.Println("I know what's in your future!!!\n");

	// get info
	fmt.Print("What is your name?: ");
	fmt.Scan(&name);

	fmt.Print("When were you born (mm/dd/yyyy)?: ");
	fmt.Scan(&birthDate);

	// format date: https://www.geeksforgeeks.org/go-language/time-parse-function-in-golang-with-examples/
	formattedBirthDate, err := time.Parse(DATE_LAYOUT, birthDate);
	if err != nil {
		fmt.Println("Error formatting birth date: ", err);
		return;
	}

	// get info
	userAge := calculateAge(formattedBirthDate);
	zodiacSign := getZodiacSign(formattedBirthDate);
	randomFortuneMessage := messages[generateRandNumber(messageArrayLength)];

	// output info
	fmt.Printf("\n%s, you are %d years old\n", name, userAge);
	fmt.Printf("You are a: %s\n", zodiacSign);
	fmt.Printf("Fortune Message: \n %s\n\n", randomFortuneMessage);
}