package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
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
	data := fetchCourseList(URL);

	fmt.Println(data.Items[0].Children[0].Times);
}