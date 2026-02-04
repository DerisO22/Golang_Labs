package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type AddressStruct struct {
	street string
	city string
	zip string
	state string
}
type PersonStruct struct {
	name string
	address AddressStruct
	phoneNumber string
}
type StudentStruct struct {
	person PersonStruct
	id string
	major string
	gpa float32
}

const NUM_OF_STUDENTS = 5;

func display(student map[string]StudentStruct) {
	fmt.Println("Students:")
	
	for _, value := range student {
		fmt.Println("ID:", value.id)
		fmt.Println("Name:", value.person.name)
		fmt.Println("Major:", value.major)
		fmt.Println("GPA:", value.gpa)
		fmt.Println("Address:", value.person.address.street,
		value.person.address.city, value.person.address.state,
		value.person.address.zip)
		fmt.Println("Phone Number:", value.person.phoneNumber)
		fmt.Println()
	}
}

func getUserInput() {
	// student info - student struct
	var studentInfo map[string]StudentStruct;
	
	// using code 3 from lab 2
	reader := bufio.NewReader(os.Stdin);

	studentInfo = make(map[string]StudentStruct);

	for i := 0; i < NUM_OF_STUDENTS; i++ {
		fmt.Print("\nWhat is your name?: ");
		nameInput, err := reader.ReadString('\n');
		if err != nil {
			log.Fatal(err);
		}
		// Trim the newline character from the end of the string
		nameInput = strings.TrimSpace(nameInput);

		fmt.Print("What state do you live in?: ");
		stateInput, err := reader.ReadString('\n');
		if err != nil {
			log.Fatal(err);
		}
		// Trim the newline character from the end of the string
		stateInput = strings.TrimSpace(stateInput);

		fmt.Print("What city do you live in?: ");
		cityInput, err := reader.ReadString('\n');
		if err != nil {
			log.Fatal(err);
		}
		// Trim the newline character from the end of the string
		cityInput = strings.TrimSpace(cityInput);

		fmt.Print("What is your ZIP Code?: ");
		zipInput, err := reader.ReadString('\n');
		if err != nil {
			log.Fatal(err);
		}
		// Trim the newline character from the end of the string
		zipInput = strings.TrimSpace(zipInput);

		fmt.Print("What is your street address?: ");
		streetAddressInput, err := reader.ReadString('\n');
		if err != nil {
			log.Fatal(err);
		}
		// Trim the newline character from the end of the string
		streetAddressInput = strings.TrimSpace(streetAddressInput);

		fmt.Print("What is your phone number?: ");
		phoneNumberInput, err := reader.ReadString('\n');
		if err != nil {
			log.Fatal(err);
		}
		// Trim the newline character from the end of the string
		phoneNumberInput = strings.TrimSpace(phoneNumberInput);

		fmt.Print("What is your major?: ");
		majorInput, err := reader.ReadString('\n');
		if err != nil {
			log.Fatal(err);
		}
		// Trim the newline character from the end of the string
		majorInput = strings.TrimSpace(majorInput);

		var gpaInput = -1.0;

		for (gpaInput < 0.0 || gpaInput > 4.0) {
			fmt.Print("What is your gpa?: ");
			gpaInputString, err := reader.ReadString('\n');

			if err != nil {
				log.Fatal(err);
			}

			gpaInputString = strings.TrimSpace(gpaInputString);

			// Converting it to float32
			gpaInputFloat, errFloat := strconv.ParseFloat(gpaInputString, 32);
			
			gpaInput = gpaInputFloat;

			if errFloat != nil {
				log.Fatal(err);
			}

			if (gpaInputFloat < 0.0 || gpaInputFloat > 4.0) {
				fmt.Println("Invalid GPA Input\n");
			}
		}

		studentID := fmt.Sprintf("%08.2f", float32(i + 1));

		studentInfo[studentID] = StudentStruct{
			person: PersonStruct{
				name: nameInput,
				address: AddressStruct{
					street: streetAddressInput,
					city: cityInput,
					zip: zipInput,
					state: stateInput,
				},
				phoneNumber: phoneNumberInput,
			},
			id: studentID,
			major: majorInput,
			gpa: float32(gpaInput),
		}

		display(studentInfo);
	}
}

func main() {
	// Get User input for the 5 users
	getUserInput();
}
