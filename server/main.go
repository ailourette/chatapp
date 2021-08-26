package main

import (
	"html/template"
	"net/http"
	"sync"
)

// Declaration of global variable
var (
	//tablePatient     [4][8]string
	//doctors          = [4]string{"Amy", "Collin", "Jenkins", "Sarah"}
	//doctorTime       = [8]string{"9am to 10am", "10am to 11am", "11am to 12pm", "1pm to 2pm", "2pm to 3pm", "3pm to 4pm", "4pm to 5pm", "5pm to 6pm"}
	tpl   *template.Template
	mutex sync.Mutex // Concurrency
)

func init() {
	tpl = template.Must(template.ParseGlob("templates/*"))
}

// This function initialise sessions from the saved csv file to the dbSession map which is a map of key cookie value, value username
// func initSession() { // Convert CSV file 2D array to patient account map
// 	records := sessionReadCsv()
// 	for i := 0; i < len(records); i++ {
// 		for j := 0; j < 2; j++ {
// 			cookieValue := records[i][0]
// 			un := records[i][1]
// 			dbSessions[cookieValue] = un
// 		}
// 	}
// }

func home(res http.ResponseWriter, req *http.Request) {
	tpl.ExecuteTemplate(res, "homePage.gohtml", nil)
}

func userLoginSuccess(res http.ResponseWriter, req *http.Request) {
	tpl.ExecuteTemplate(res, "patientLoginSuccess.gohtml", nil)
}

// func adminLoginSuccess(res http.ResponseWriter, req *http.Request) {
// 	tpl.ExecuteTemplate(res, "adminLoginSuccess.gohtml", nil)
// }

// func doctorLoginSuccess(res http.ResponseWriter, req *http.Request) {
// 	tpl.ExecuteTemplate(res, "doctorLoginSuccess.gohtml", nil)
// }

func main() {
	http.HandleFunc("/", home)
	//http.HandleFunc("/logout", logout)
	http.HandleFunc("/usersignup", userSignup)
	http.HandleFunc("/userlogin", userLogin)
	http.HandleFunc("/userloginsuccess", userLoginSuccess)
	//http.HandleFunc("/adminLogin", adminLogin)
	//http.HandleFunc("/adminLoginSuccess", adminLoginSuccess)
	//http.HandleFunc("/deletesessions", deleteSessions)
	//http.HandleFunc("/deleteusers", deleteUsers)
	//http.HandleFunc("/docWriteNotes", docWriteNotes)
	//http.HandleFunc("/docReadNotes", docReadNotes)
	http.HandleFunc("/userchangepassword", userChangePassword)
	http.ListenAndServeTLS(":5221", "cert.pem", "key.pem", nil)
}
