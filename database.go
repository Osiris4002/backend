package main

import (
	"math/rand"
	"strconv"
	"time"
    "database/sql"
    "fmt"
    _ "github.com/lib/pq"
)

var db *sql.DB

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}

// WARNING: THE PASSWORD IS STORED IN PLAIN-TEXT IN THE DATABASE WHICH IS HIGHLY INSECURE
// MUST BE CHANGED BEFORE THE BACKEND GETS RELEASED TO THE PUBLIC
func AddDataToDatabase(student_data students_database) {
	insertdata := `insert into "students_data"("roll_no", "name", "hostel_code", "password") values($1, $2, $3, $4)`
    _, err := db.Exec(insertdata, student_data.Roll_No, student_data.Name, student_data.Hostel_Code, student_data.Password)
    CheckError(err)
    Students_Data = append(Students_Data, student_data)
}

func AddComplaintToDatabase(complaint_data students_complaint) {
	// Randomly generate the unique id of the complaint
	rand.Seed(time.Now().UnixNano())
	complaint_data.Uid = strconv.Itoa(rand.Intn(10000000 - 1 + 1) + 1)

	insertdata := `insert into "complaint_data" values($1, $2, $3, $4)`
    _, err := db.Exec(insertdata, complaint_data.Uid, complaint_data.Complaint_Text, complaint_data.Complaint_Text_Title, complaint_data.Roll_No, complaint_data.Hostel_Code)
    CheckError(err)
}

func CheckUsernameAndPassword (username string, password string) bool {
	// WARNING: THE PASSWORD IS STORED IN PLAIN-TEXT IN THE DATABASE WHICH IS HIGHLY INSECURE
	// MUST BE CHANGED BEFORE THE BACKEND GETS RELEASED TO THE PUBLIC
	var pass string
	query := fmt.Sprintf("SELECT password from admin_data WHERE username='%s'", username)
	rows, err := db.Query(query)
	CheckError(err)
	defer rows.Close()

	// First check whether the roll no. exists in the database
	rows.Next()
	err = rows.Scan(&pass)
	CheckError(err)
	if pass == "" {
		return false
	} else if password != pass {
		return false
	}
	return true
}

func CheckRollNoAndPassword (roll_no string, password string) bool {
	// WARNING: THE PASSWORD IS STORED IN PLAIN-TEXT IN THE DATABASE WHICH IS HIGHLY INSECURE
	// MUST BE CHANGED BEFORE THE BACKEND GETS RELEASED TO THE PUBLIC
	var pass string
	query := fmt.Sprintf("SELECT password from students_data WHERE roll_no='%s'", roll_no)
	rows, err := db.Query(query)
	CheckError(err)
	defer rows.Close()

	// First check whether the roll no. exists in the database
	rows.Next()
	err = rows.Scan(&pass)
	CheckError(err)
	if pass == "" {
		return false
	} else if password != pass {
		return false
	}
	return true
}

func GatherDataFromDatabase() {
	var admin_data admin_database
	var students_data students_database
	var student_complaint_data students_complaint

	rows, err := db.Query("SELECT username, name, hostel_code from admin_data")
	CheckError(err)
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&admin_data.Username, &admin_data.Name, &admin_data.Hostel_Code)
		CheckError(err)
		Admin_Data = append(Admin_Data, admin_data)
	}

	rows, err = db.Query("SELECT roll_no, name, hostel_code from students_data")
	CheckError(err)
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&students_data.Roll_No, &students_data.Name, &students_data.Hostel_Code)
		CheckError(err)
		Students_Data = append(Students_Data, students_data)
	}

	rows, err = db.Query("SELECT * from complaint_data")
	CheckError(err)
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&student_complaint_data.Uid, &student_complaint_data.Complaint_Text, &student_complaint_data.Complaint_Text_Title, &student_complaint_data.Roll_No)
		CheckError(err)
		Complaint_Data = append(Complaint_Data, student_complaint_data)
	}
}

func UserComplaintResolve(query_resolve string, uid string) {
	insertdata := `update "complaint_data" set "query_resolved" $1 where uid = $2' `
    _, err := db.Exec(insertdata, query_resolve, uid)
    CheckError(err)
}

func OpenDatabase() {
	var err error
	psqlconn := fmt.Sprintf("host=localhost port=5432 user=postgres password=12345 dbname=students_database sslmode=disable")
	db, err = sql.Open("postgres", psqlconn)
	if (err != nil) {
		CheckError(err)
	}

	err = db.Ping()
    if (err != nil) {
		CheckError(err)
	}
	
	fmt.Println("Database Connected!")
}
