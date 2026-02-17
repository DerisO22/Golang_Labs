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
	id     int64
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
const DATABASE = "pa3"
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
	var studentID int8;

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
			fmt.Printf("%-3d %-20s %-20s %-20s %-20s %-20s %6.2f %-65s \n", i, value.name, value.dob, value.email, value.phone, value.major, value.gpa, value.class);
			i++
		}
	}

	return fmt.Errorf("")
}

func insertNewRecord(db *sql.DB) (int64, error) {
	var newAdvisor AdvisorInfo;
	var idInput int;

	// Getting input to add advisor
	reader := bufio.NewReader(os.Stdin);
	fmt.Println("Insert A New Advisor Into Database:")

	fmt.Print("Enter advisor id?: ");
	fmt.Scanln(&idInput);

	newAdvisor.id = int64(idInput);
	
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

	fmt.Println(newAdvisor);

	advisorResult, advisorErr := db.Exec("INSERT INTO person (id, name, dob, email, phone) VALUES (?, ?, ?, ?, ?);", newAdvisor.id, newAdvisor.name, newAdvisor.dob, newAdvisor.email, newAdvisor.phone);
	
	if advisorErr != nil {
		return 0, fmt.Errorf("Error Inserting a New Advisor: %v", advisorErr);
	}

	id, advisorErr := advisorResult.LastInsertId();
	if advisorErr != nil {
		return 0, fmt.Errorf("Error Inserting a New Advisor: %v", advisorErr);
	}

	_, studentErr := db.Exec("INSERT INTO student (id, major, gpa, class) VALUES (?, ?, ?, ?)", newAdvisor.id, newAdvisor.major, newAdvisor.gpa, newAdvisor.class);

	if studentErr != nil {
		return 0, fmt.Errorf("Error Inserting a New Student: %v", studentErr);
	}

	return id, nil;
}

func searchAdviseeByName(db *sql.DB, name string) error {
	var advisors []Advisor;
	var tmp Advisor;

	rows, err := db.Query("SELECT * FROM person WHERE name = ?", name);
	if err != nil {
		return fmt.Errorf("AdviseeByName %q: %v", name, err);
	}
	defer rows.Close();

	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		err := rows.Scan(&tmp.id, &tmp.name, &tmp.dob, &tmp.email, &tmp.phone);

		if err != nil {
			return fmt.Errorf("AdviseeByName %q: %v", name, err);
		}
		advisors = append(advisors, tmp);
	}

	err = rows.Err();

	if err != nil {
		return fmt.Errorf("AdviseeByName %q: %v", name, err);
	} else {
		i := 1;
		for _, value := range advisors {
			fmt.Printf("%-3d %-20s %-20s %-20s %-20s\n", i, value.name, value.dob, value.email, value.phone);
			i++
		}
	}

	return nil
}

func deleteARecord(db *sql.DB, name string) error {
	query := "DELETE FROM person WHERE name = ?";

	// Execute the query using db.Exec()
	result, err := db.Exec(query, name);

	if err != nil {
		return err;
	}

	// Optional: Check how many rows were affected
	rowsAffected, err := result.RowsAffected();
	if err != nil {
		return fmt.Errorf("Could not get rows affected: %v\n", err);
	}

	fmt.Printf("Deleted %d row(s)\n", rowsAffected);

	return fmt.Errorf("");
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

	nameToSearch := "Test1";
	fmt.Printf("Searching for Advisor - %s:\n", nameToSearch)
	err := searchAdviseeByName(db, nameToSearch);
	if err != nil {
		fmt.Println(err);
	}

	id, err := insertNewRecord(db);

	if err != nil {
		fmt.Println("Added a New Record with id = ", id);
	} else {
		fmt.Println("Error in Adding a New Record");
	}

	deleteErr := deleteARecord(db, "Deris");

	if deleteErr != nil {
		fmt.Println("Deleted Record with name = Deris");
	}
}
