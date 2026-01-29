/*
Author:                   Deris C. O’Malley
Class:                    CSI-380-01
Assignment:               PA 1
Date Assigned:            29, January
Due Date:                 5, February, 11:00AM

Description:
This program will ask the users for their name and date of
birth, then give them a fortune message. This program will
also tell the user’s age, zodiac sign, and fortune message.

Certification of Authenticity:
I certify that this is entirely my own work,except where I have given fully-documented
References to the work of others. I understand the definition and consequences of
Plagiarism and acknowledge that the assessor of this assignment may, for the purpose
of assessing this assignment:
-Reproduce this assignment and provide a copy to another member of academic staff;
and/or
- Communicate a copy of this assignment to a plagiarism checking service (which
May then retain a copy of this assignment on its database for the purpose of future
Plagiarism checking)
*/

package main

import (
	"bufio"
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"os"
	"strings"
	"time"
)

const FILE_TO_READ = "messages.txt";
const DATE_LAYOUT = "01/02/2006";

// reads and add to list
func readFromFile(messageArray *[]string) {
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
	if currentTime.Month() < inputBirthDate.Month() || (inputBirthDate.Month() == currentTime.Month() && currentTime.Day() < inputBirthDate.Day()){
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

func getUserInputs(name *string, birthDate *time.Time) {
	// get input
	var birthDateInput string;
	validBirthDate := false;

	// using code 3 from lab 2
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("What is your name?: ")
	
	// ReadString reads until the delimiter '\n' (newline)
	nameInput, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	// Trim the newline character from the end of the string
	nameInput = strings.TrimSpace(nameInput);
	*name = nameInput;

	for validBirthDate == false {
		fmt.Print("When were you born (mm/dd/yyyy)?: ");
		birthDateInput, _= reader.ReadString('\n');
		birthDateInput = strings.TrimSpace(birthDateInput);
		
		// format date: https://www.geeksforgeeks.org/go-language/time-parse-function-in-golang-with-examples/
		formattedBirthDate, err := time.Parse(DATE_LAYOUT, birthDateInput);
		if err != nil {
			fmt.Println("Error formatting birth date: ", err);
			continue;
		}

		// extra checks on valid date input
		if formattedBirthDate.After(time.Now()) {
			fmt.Println("Invalid! Birthdate cannot be in the future.\n");
			continue;
		}

		// Oldest person ever was 122, so anything over that is somewhat unrealistic
		if calculateAge(formattedBirthDate) > 122 {
			fmt.Println("Invalid! Enter a reasonable birthday\n");
			continue;
		}

		validBirthDate = true;
		*birthDate = formattedBirthDate;
	}
}

func main() {
	// using slice like the generateArray func from lab3 example
	var messages []string;
	var birthDate time.Time;
	var name string;

	// I could just return the new populated messages slice
	// in readFromFile, but wanted to try out pointers/pass-by-ref
	readFromFile(&messages);

	fmt.Println("\nI'm am the Almighty Psychic Fortune Teller!");
	fmt.Println("I know what's in your future!!!\n");

	// get input
	getUserInputs(&name, &birthDate);

	// get info
	userAge := calculateAge(birthDate);
	zodiacSign := getZodiacSign(birthDate);
	randomFortuneMessage := messages[generateRandNumber(len(messages))];

	// output info
	fmt.Printf("\n%s, you are %d years old\n", name, userAge);
	fmt.Printf("Your Zodiac Sign: %s\n", zodiacSign);
	fmt.Printf("Fortune Message: \n %s\n\n", randomFortuneMessage);
}