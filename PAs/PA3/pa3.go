package main

import (
	"database/sql"
	"fmt"

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

	fmt.Println("Connected to database: ", DATABASE)

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
			fmt.Printf("%-3d %-20s %-20s %-20s %-20s %-20s $%6.2f %-65s \n", i, value.name, value.dob, value.email, value.phone, value.major, value.gpa, value.class);
			i++
		}
	}

	return fmt.Errorf("")
}

func main() {
	var db *sql.DB
	var connected bool = false

	connected, db = connectToADatabase()

	if connected == false {
		return;
	}

	fmt.Println("Hello")
	displayAllRecords(db)
}
