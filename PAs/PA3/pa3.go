package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type Advisor struct {
	id     string
	name   string
	dob    string
	email  string
	phone  string
}

type Student struct {
	major string
	gpa   float32
	class string
}

type AdvisorInfo struct {
	Advisor
	Student
}

const USER = "root"
const PASSWD = "Deris123" // put in your password here
const DATABASE = "advising"
const CONNECTION = "tcp"
const HOST = "127.0.0.1"
const PORT = "3306"

func connectToADatabase() (bool, *sql.DB) {
	var db *sql.DB
	var err error

	query := fmt.Sprintf("%s:%s@%s(%s:%s)/%s?parseTime=true", USER, PASSWD, CONNECTION, HOST, PORT, DATABASE)

	// Get a database handle.
	db, err = sql.Open("mysql", query)
	if err != nil {
		fmt.Println("MySQL Open Error: ", err)
		return false, db
	}

	err = db.Ping()
	if err != nil {
		fmt.Println("MySQL Ping Error: ", err)
		return false, db
	}

	fmt.Println("\nConnected to database: ", DATABASE)

	return true, db
}

func displayAllRecords(db *sql.DB) error {
	var advisors []AdvisorInfo;
	var tmp AdvisorInfo;
	var studentID string;

	rows, err := db.Query("SELECT * FROM person JOIN student ON person.id=student.id;");
	if err != nil {
		return fmt.Errorf("displayAllRecords: %v", err)
	}

	defer rows.Close()
	
	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		err = rows.Scan(&tmp.id, &tmp.name, &tmp.dob, &tmp.email, &tmp.phone, &studentID, &tmp.major, &tmp.gpa, &tmp.class);
		if err != nil {
			return fmt.Errorf("displayAllRecords: %v", err)
		}
		advisors = append(advisors, tmp)
	}
	err = rows.Err()

	if err != nil {
		return fmt.Errorf("displayAllRecords: %v", err)
	} else {
		i := 1
		for _, value := range advisors {
			fmt.Printf("%-3d %-20s %-20s %-20s %-20s %-20s %6.3f %-65s \n", i, value.name, value.dob, value.email, value.phone, value.major, value.gpa, value.class);
			i++
		}
	}

	return nil;
}

func insertNewRecord(db *sql.DB) (string, error) {
	var newAdvisor AdvisorInfo;

	// Getting input to add advisor
	reader := bufio.NewReader(os.Stdin);
	fmt.Println("\nInsert A New Advisor Into Database:")

	fmt.Print("Enter advisor id?: ");
	idInput, err := reader.ReadString('\n');
	if err != nil {
		log.Fatal(err);
	}
	newAdvisor.id = strings.TrimSpace(idInput);

	fmt.Print("Enter advisor name?: ");
	nameInput, err := reader.ReadString('\n');
	if err != nil {
		log.Fatal(err);
	}
	newAdvisor.name = strings.TrimSpace(nameInput);

	fmt.Print("Enter advisor dob(dd/mm/yyyy)?: ");
	dob, err := reader.ReadString('\n');
	if err != nil {
		log.Fatal(err);
	}
	newAdvisor.dob = strings.TrimSpace(dob);

	fmt.Print("Enter advisor email?: ");
	email, err := reader.ReadString('\n');
	if err != nil {
		log.Fatal(err);
	}
	newAdvisor.email = strings.TrimSpace(email);

	fmt.Print("Enter advisor phone_number?: ");
	phone_number, err := reader.ReadString('\n');
	if err != nil {
		log.Fatal(err);
	}
	newAdvisor.phone = strings.TrimSpace(phone_number);

	fmt.Print("Enter student major?: ");
	major, err := reader.ReadString('\n');
	if err != nil {
		log.Fatal(err);
	}
	newAdvisor.major = strings.TrimSpace(major);

	fmt.Print("Enter student class?: ");
	class, err := reader.ReadString('\n');
	if err != nil {
		log.Fatal(err);
	}
	newAdvisor.class = strings.TrimSpace(class);

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

	newAdvisor.gpa = float32(gpaInput);

	fmt.Println("");
	fmt.Println(newAdvisor);

	_, advisorErr := db.Exec("INSERT INTO person (id, name, dob, email, phone) VALUES (?, ?, ?, ?, ?);", newAdvisor.id, newAdvisor.name, newAdvisor.dob, newAdvisor.email, newAdvisor.phone);
	
	if advisorErr != nil {
		return "0", fmt.Errorf("Error Inserting a New Advisor: %v", advisorErr);
	}

	/*
	// This doesn't work since the primary keys aren't auto increment
	id, advisorErr := advisorResult.LastInsertId();
	if advisorErr != nil {
		return "0", fmt.Errorf("Error Inserting a New Advisor: %v", advisorErr);
	}
	*/

	_, studentErr := db.Exec("INSERT INTO student (id, major, gpa, class) VALUES (?, ?, ?, ?)", newAdvisor.id, newAdvisor.major, newAdvisor.gpa, newAdvisor.class);

	if studentErr != nil {
		return "0", fmt.Errorf("Error Inserting a New Student: %v", studentErr);
	}

	return newAdvisor.id, nil;
}

func searchAdviseeByName(db *sql.DB) error {
	var advisors []Advisor;
	var tmp Advisor;

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter name to search for?: ")
	
	nameToSearch, _ := reader.ReadString('\n')
	nameToSearch = strings.TrimSpace(nameToSearch)

	rows, err := db.Query("SELECT * FROM person WHERE name = ?", nameToSearch);
	if err != nil {
		return fmt.Errorf("AdviseeByName %q: %v", nameToSearch, err);
	}
	defer rows.Close();

	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		err := rows.Scan(&tmp.id, &tmp.name, &tmp.dob, &tmp.email, &tmp.phone);

		if err != nil {
			return fmt.Errorf("AdviseeByName %q: %v", nameToSearch, err);
		}
		advisors = append(advisors, tmp);
	}

	err = rows.Err();

	if err != nil {
		return fmt.Errorf("AdviseeByName %q: %v", nameToSearch, err);
	} else {
		i := 1;
		for _, value := range advisors {
			fmt.Printf("%-3d %-20s %-20s %-20s %-20s\n", i, value.name, value.dob, value.email, value.phone);
			i++
		}
	}

	return nil;
}

func deleteARecord(db *sql.DB) error {
	// Getting input to add advisor
	reader := bufio.NewReader(os.Stdin);
	fmt.Println("\nDelete an Advisor From Database:")

	fmt.Print("Enter advisor id to delete?: ");
	idInput, errInput := reader.ReadString('\n');
	if errInput != nil {
		log.Fatal(errInput);
	}
	idInput = strings.TrimSpace(idInput);

	delete_student_query := "DELETE FROM student WHERE id = ?";
	delete_person_query := "DELETE FROM person WHERE id = ?"

	_, err := db.Exec(delete_student_query, idInput)
	if err != nil {
		return err
	}
	result, err := db.Exec(delete_person_query, idInput)

	if err != nil {
		return err;
	}

	// Optional: Check how many rows were affected
	rowsAffected, err := result.RowsAffected();
	if err != nil {
		return fmt.Errorf("Could not get rows affected: %v\n", err);
	}

	fmt.Printf("\nDeleted %d row(s)\n", rowsAffected);

	return nil;
}

func main() {
	var db *sql.DB;
	var connected bool = false;

	connected, db = connectToADatabase();

	if connected == false {
		return;
	}

	fmt.Println("\nAll Advisor and Their Student Information:");
	displayAllRecords(db);

	err := searchAdviseeByName(db);
	if err != nil {
		fmt.Println(err);
	}

	id, err := insertNewRecord(db);

	if err != nil {
		fmt.Println("Error in Adding a New Record");
	} else {
		fmt.Println("Added a New Record with id = ", id);
	}

	// Deleting
	deleteErr := deleteARecord(db);

	if deleteErr != nil {
		fmt.Println("Error deleting record")
	} else {
		fmt.Println("Successfully deleted record")
	}

	fmt.Println("\nAll Advisor and Their Student Information after CR(U)D Operations:");
	displayAllRecords(db);
}