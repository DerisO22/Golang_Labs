package main

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
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
	// I added the IF NOT EXISTS to prevent errors
	_, err := db.Exec("CREATE DATABASE IF NOT EXISTS coursesList");
	if err != nil {
		fmt.Println("Create Database Error: ", err)
	}

    _, err = db.Exec("USE coursesList");

	_, err = db.Exec("DROP TABLE IF EXISTS courses;")
	if err != nil {
		fmt.Println("Error dropping table: ", err)
	}

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

	// semester.children has the courses
	for _, course := range coursesData.Children {
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
		fmt.Printf("Error committing transaction: %v\n", err);
	} else {
		fmt.Println("Successfully saved all courses!");
	}
}

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

// Help function, so I don't need to repeat the row iteration for each search func
func printCourses(rows *sql.Rows) {
	defer rows.Close()
	var courses []Course

	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		var tmp Course
		err := rows.Scan(&tmp.Number, &tmp.Days, &tmp.Times, &tmp.Room, &tmp.Instructor_fname, &tmp.Instructor_lname, &tmp.Openseats)
		if err != nil {
			fmt.Printf("Scan error: %v\n", err)
			return
		}
		courses = append(courses, tmp)
	}

	if len(courses) == 0 {
		fmt.Println("No courses found.")
		return
	}

	// col headers. but will prob not be seen if data is too large
	fmt.Printf("\n%-10s %-8s %-20s %-15s %-25s %-5s\n", "NUMBER", "DAYS", "TIME", "ROOM", "INSTRUCTOR", "SEATS")
	for _, c := range courses {
		name := c.Instructor_fname + " " + c.Instructor_lname

		// *** Printing things out based on their varchar sizes overflowed terminal on my laptop, so just abritary vals***
		fmt.Printf("%-10s %-8s %-20s %-15s %-25s %-5s\n", c.Number, c.Days, c.Times, c.Room, name, c.Openseats)
	}
}

/*
	All the User Operations
*/
func displayFallCourses(userInput int, db *sql.DB) error {
	var query string;

	fmt.Println("\nAll Fall Courses:\n");

	switch(userInput){
		case 1:
			query = "SELECT number, days, times, room, instructorFname, instructorLname, openSeats FROM courses;"
		case 2:
			query = "SELECT number, days, times, room, instructorFname, instructorLname, openSeats FROM courses WHERE courseType = 'Day/Evening';"
		case 3:
			query = "SELECT number, days, times, room, instructorFname, instructorLname, openSeats FROM courses WHERE courseType != 'Day/Evening';";
	}

	rows, err := db.Query(query);
	if err != nil {
		return fmt.Errorf("displayFallCourses: %v", err);
	}

	printCourses(rows);

	return nil;
}

func searchByCoursePrefix(db *sql.DB) {
	var prefixToSearch string
	fmt.Print("\nEnter a course-prefix to search (CSI, MTH, etc)?: ");
	
	fmt.Scan(&prefixToSearch)

	query := "SELECT number, days, times, room, instructorFname, instructorLname, openSeats FROM courses WHERE number LIKE ?;";
	rows, err := db.Query(query, prefixToSearch + "%");
	if err != nil {
		fmt.Printf("courseByPrefix %q: %v", prefixToSearch, err);
		return;
	}
	
	printCourses(rows);
}

func searchByLevel(db *sql.DB) {
	var levelToSearch string
	fmt.Print("\nEnter a course-level to search (1xx, 2xx, etc)?: ");
	
	fmt.Scan(&levelToSearch)

	query := "SELECT number, days, times, room, instructorFname, instructorLname, openSeats FROM courses WHERE number LIKE ?;";
	rows, err := db.Query(query, "%" + levelToSearch + "%");
	if err != nil {
		fmt.Printf("courseByLevel %q: %v", levelToSearch, err);
		return;
	}
	
	printCourses(rows);
}

func searchByPrefixAndLevel(db *sql.DB) {
	var courseToSearch string
	reader := bufio.NewReader(os.Stdin);

	fmt.Print("Enter Course Prefix and Level (e.g CSI 300, MTH 270): ");
	class, err := reader.ReadString('\n');
	if err != nil {
		log.Fatal(err);
	}
	courseToSearch = strings.TrimSpace(class);

	query := "SELECT number, days, times, room, instructorFname, instructorLname, openSeats FROM courses WHERE number LIKE ?;";
	rows, err := db.Query(query, courseToSearch + "%");
	if err != nil {
		fmt.Printf("courseByPrefixAndLevel %q: %v", courseToSearch, err);
		return;
	}
	
	printCourses(rows);
}

func searchForOpenCourses(db *sql.DB) {
	fmt.Println("\nAll Open Courses");

	query := "SELECT number, days, times, room, instructorFname, instructorLname, openSeats FROM courses WHERE openSeats != '0';";
	rows, err := db.Query(query);
	if err != nil {
		fmt.Printf("openCourses : %v", err);
		return;
	}
	
	printCourses(rows);
}

func searchForSpecificCourse(db *sql.DB) {
	var tmp Course;
	var courseToSearch string
	reader := bufio.NewReader(os.Stdin);

	fmt.Print("Enter Course to Search (e.g CSI 230-02, MTH 270-01): ");
	class, err := reader.ReadString('\n');
	if err != nil {
		log.Fatal(err);
	}
	courseToSearch = strings.TrimSpace(class);

	query := "SELECT * FROM courses WHERE number = ? LIMIT 1;";
	row := db.QueryRow(query, courseToSearch);
	err = row.Scan(&tmp.Number, &tmp.Credit, &tmp.Openseats, &tmp.Days, &tmp.Times, &tmp.Instructor_fname, 
					&tmp.Instructor_lname, &tmp.Description, &tmp.Room, &tmp.Subject, &tmp.Type, &tmp.Prereq,
					&tmp.Title, &tmp.Start_date, &tmp.End_date);
	
	if err == sql.ErrNoRows {
		fmt.Printf("Course not found: %v", err);
		return;
	} else if err != nil {
		fmt.Printf("Error: %v", err);
		return;
	}

	const courseTemplate = `
	------ %s ------
	Course Number: %-20s 
	Credits:       %-20s
	Open Seats:    %-20s  
	Days:          %-20s
	Times:         %-20s  
	Room:          %-20s
	Instructor:    %s %s
	Subject:       %-20s  
	Type:          %-20s
	Title:         %s
	Description:   %s
	Prereqs:       %s
	Dates:         %s to %s
	`
	
	fmt.Printf(courseTemplate, 
		courseToSearch, tmp.Number, tmp.Credit, 
		tmp.Openseats, tmp.Days, tmp.Times, tmp.Room,
		tmp.Instructor_fname, tmp.Instructor_lname,
		tmp.Subject, tmp.Type, tmp.Title,
		tmp.Description, tmp.Prereq, 
		tmp.Start_date, tmp.End_date,
	)	
}

// user input funcs
func getUserInput() int {
	userInput := 1

	for {
		// bad indentation, but outputs better
		fmt.Println(`
	Operations:
	  1. Display all courses
	  2. Display all on-campus courses
	  3. Display all CCO courses
	  4. Search for courses with a specified prefix, like CSI, GPR, MTH, etc
	  5. Search for courses with a specified prefix, like CSI, GPR, MTH, etc, and 
	     specific level, like 1xx, 2xx, etc.
	  6. Search for all courses with a specific level, like 1xx, 2xx, etc
	  7. Search for courses that are still open
	  8. Search for the details for a specific course. And, your system should return all the
	     information about the course in a nice format.
	  9. Exit Program
		`);

		fmt.Print("Enter Operation(1-9): ");
		fmt.Scan(&userInput);

		if(userInput >= 1 && userInput <= 9) {
			return userInput;
		} 
		
		fmt.Println("\nInvalid Input! Try Again");
	}
}

func handleUserInput(userInput int, db *sql.DB) {
	switch(userInput) {
		// All display related operations
		// decided to make the first 3 displays fall through
		case 1, 2, 3:
			displayFallCourses(userInput, db);			
		case 4:
			searchByCoursePrefix(db);
		case 5:
			searchByPrefixAndLevel(db);
		case 6:
			searchByLevel(db);
		case 7:
			searchForOpenCourses(db);
		case 8:
			searchForSpecificCourse(db);
		case 9:
			fmt.Println("\nExiting Program. Goodbye");
	}
}

const URL = "https://classlist.champlain.edu/api3/courses/semester/fall/type/all/filter/ug";

func main() {
	var db *sql.DB;
	var connected bool = false;
	var semestersData DataObject;
	var userInput int;
	var wg sync.WaitGroup;

	connected, db = connectToADatabase();
	if connected == false {
		return;
	}

	wg.Add(2);

	// go routines
	go func () {
		defer wg.Done();
		semestersData = fetchCourseList(URL);
	}()

	go func () {
		defer wg.Done();
		createDatabase(db);
	}()

	wg.Wait();

	if len(semestersData.Items) == 0 {
		log.Fatal("No semester data returned");
	}

	// load all data into db
	loadAllDataToDatabase(semestersData.Items[0], db);

	fmt.Println("\nWelcome to Fall 2026 Course Management System");
	for userInput != 9 {
		userInput = getUserInput();
		handleUserInput(userInput, db);
	}
}