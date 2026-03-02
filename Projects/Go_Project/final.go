package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"

	_ "github.com/go-sql-driver/mysql"
)

// Define a struct for the deeply nested item data
// Also need a struct for each idividual course
type Course struct {
	Number 				 string `json:"number"`
	Credit 				 string `json:"credit"`
	Openseats 			 string `json:"openseats"`
	Days 				 string `json:"days"`
	Times 				 string `json:"times"`

	Instructor_fname 	 string `json:"instructor_fname"`
	Instructor_lname 	 string `json:"instructor_lname"`
	Description 		 string `json:"description"`
	Room 				 string `json:"room"`
	Subject 			 string `json:"subject"`

	Type 				 string `json:"type"`
	Prereq 				 string `json:"prereq"`
	Title 				 string `json:"title"`
	Start_date 			 string `json:"start_date"`
	End_date 			 string `json:"end_date"`
	Id 					    int `json:"id"`
}

type Semester struct {
	Id 					 string `json:"id"`
	Year 				    int `json:"year"`
	Display_date         string `json:"display_date"`
	Can_register           bool `json:"can_register"`
	Type                 string `json:"type"`
	Children           []Course `json:"children"`
} 

// Define a struct for the inner 'data' object which holds an array
type DataObject struct {
	Identififer    		string `json:"identifier"`
	Items       	[]Semester `json:"items"`
}

const USER = "root"
const PASSWD = "Deris123" // put in your password here
const DATABASE = "coursesList"
const CONNECTION = "tcp"
const HOST = "127.0.0.1"
const PORT = "3306"

func connectToADatabase() (bool, *sql.DB) {
	var db *sql.DB
	var err error

	query := fmt.Sprintf("%s:%s@%s(%s:%s)/?parseTime=true", USER, PASSWD, CONNECTION, HOST, PORT);

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

	fmt.Println("\nConnected to database: ", DATABASE);

	return true, db;
}

func createDatabase(db *sql.DB) {
	_, err := db.Exec("CREATE DATABASE IF NOT EXISTS coursesList");
	if err != nil {
		fmt.Println("Create Database Error: ", err)
	}

    _, err = db.Exec("USE coursesList");

	createTableQuery := `
	CREATE TABLE courses (
		number VARCHAR(15) NOT NULL,
		credit VARCHAR(2) NOT NULL,
		openSeats VARCHAR(10) NOT NULL,
		days VARCHAR(15) NOT NULL,
		times VARCHAR(40) NOT NULL,
		instructorFname VARCHAR(50) NOT NULL,
		instructorLname VARCHAR(30) NOT NULL,
		description LONGTEXT NOT NULL,
		room VARCHAR(20) NOT NULL,
		subject VARCHAR(50) NOT NULL,
		courseType VARCHAR(20) NOT NULL,
		prereq LONGTEXT NOT NULL,
		title VARCHAR(100) NOT NULL,
		startDate VARCHAR(15) NOT NULL,
		endDate VARCHAR(15) NOT NULL
	);`

	_, err = db.Exec(createTableQuery);
	if err != nil {
		fmt.Printf("Error Creating Table: %v", err);
		return;
	}
}

func loadAllDataToDatabase(coursesData Semester, db *sql.DB) {
	transaction, err := db.Begin();
	if err != nil {
		fmt.Printf("loadAllDataToDatabase: begin transaction: %v", err);
		return;
	}

	insertQuery := `INSERT INTO courses (number, credit, openSeats, days, times, instructorFname, instructorLname, description, room, subject, courseType, prereq, title, startDate, endDate)
					VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`;

	for _, course := range coursesData.Children {
		fmt.Println(course.Title);

		_, err = transaction.Exec(insertQuery,
			course.Number, course.Credit, course.Openseats, course.Days, course.Times, course.Instructor_fname, course.Instructor_lname,
			course.Description, course.Room, course.Subject, course.Type, course.Prereq, course.Title, course.Start_date, course.End_date);

		if err != nil {
			transaction.Rollback();
			fmt.Printf("Error Adding Data: %v", err);
			return;
		}
	}

	err = transaction.Commit()
	if err != nil {
		fmt.Printf("Error committing transaction: %v\n", err)
	} else {
		fmt.Println("Successfully saved all courses!")
	}
}

// func displayAllFallCourses(db *sql.DB) error {
// 	var Semesters []Semester;
// 	var tmp AdvisorInfo;
// 	var studentID string;

// 	rows, err := db.Query("SELECT * FROM person JOIN student ON person.id=student.id;");
// 	if err != nil {
// 		return fmt.Errorf("displayAllRecords: %v", err)
// 	}

// 	defer rows.Close()
	
// 	// Loop through rows, using Scan to assign column data to struct fields.
// 	for rows.Next() {
// 		err = rows.Scan(&tmp.id, &tmp.name, &tmp.dob, &tmp.email, &tmp.phone, &studentID, &tmp.major, &tmp.gpa, &tmp.class);
// 		if err != nil {
// 			return fmt.Errorf("displayAllRecords: %v", err)
// 		}
// 		advisors = append(advisors, tmp)
// 	}
// 	err = rows.Err()

// 	if err != nil {
// 		return fmt.Errorf("displayAllRecords: %v", err)
// 	} else {
// 		i := 1
// 		for _, value := range advisors {
// 			fmt.Printf("%-3d %-20s %-20s %-20s %-20s %-20s %6.3f %-65s \n", i, value.name, value.dob, value.email, value.phone, value.major, value.gpa, value.class);
// 			i++
// 		}
// 	}

// 	return nil;
// }

func fetchCourseList(URL string) DataObject {
	// Make an HTTP GET request
	resp, err := http.Get(URL);
	if err != nil {
		log.Fatalf("Error fetching API URL %s: %v", URL, err);
	}
	defer resp.Body.Close();

	var data DataObject;

	// Read the response body
	body, err := io.ReadAll(resp.Body);
	if err != nil {
		log.Fatal(err);
	}

	err = json.Unmarshal(body, &data);
	if err != nil {
		log.Fatal(err);
	}

	return data;
}

const URL = "https://classlist.champlain.edu/api3/courses/semester/fall/type/all/filter/ug";

func main() {
	var db *sql.DB;
	var connected bool = false;
	var semestersData DataObject;
	var wg sync.WaitGroup;

	connected, db = connectToADatabase();
	if connected == false {
		return;
	}

	wg.Add(2);

	// Fetching Routine
	go func () {
		defer wg.Done();
		semestersData = fetchCourseList(URL);
	}()

	go func () {
		defer wg.Done();
		createDatabase(db);
	}()

	wg.Wait();

	// load all data into db
	loadAllDataToDatabase(semestersData.Items[0], db);

	// for inputIsValid != false {

	// }
	// getUserInput();
	// handleUserInput();

	fmt.Println("\nAll Advisor and Their Student Information:");
	//displayAllRecords(db);

	fmt.Println(semestersData.Items[0].Children[0].Times);
}