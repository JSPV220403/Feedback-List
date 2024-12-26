package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"text/template"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func init() {
	var err error
	// Update DSN with your MySQL credentials
	db, err = sql.Open("mysql", "root:jspv@tcp(127.0.0.1:3306)/crud_demo")
	if err != nil {
		log.Fatal("Error connecting to the database:", err)
	}
	if err = db.Ping(); err != nil {
		log.Fatal("Database is not reachable:", err)
	}
}
func send(w http.ResponseWriter, r *http.Request) {
	var id string
	var title string
	var platform string
	var modules string
	var descriptions string
	var tags string
	var votes string
	var reports []map[string]string
	row, err := db.Query("select * from report")
	if err != nil {
		fmt.Println("Unable to Retrive")
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
	defer row.Close()
	for row.Next() {
		row.Scan(&id, &title, &platform, &modules, &descriptions, &tags, &votes)
		reports = append(reports, map[string]string{"ID": id, "Title": title, "Platform": platform, "Modules": modules, "Descriptions": descriptions, "Tags": tags, "Votes": votes})
	}
	temp := template.Must(template.ParseFiles("homepage.html"))
	temp.Execute(w, reports)
}
func add(w http.ResponseWriter, r *http.Request) {
	title := r.FormValue("_title")
	platform := r.FormValue("_platform")
	modules := r.FormValue("_modules")
	descriptions := r.FormValue("_desc")
	tags := r.FormValue("buttonValue")
	_, err := db.Query("insert into report(title, platform, modules, descriptions, tags, votes) values(?,?,?,?,?,1)", title, platform, modules, descriptions, tags)
	if err != nil {
		fmt.Println("Unable to Insert")
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
	fmt.Println("successfully Inserted")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
func addCall(w http.ResponseWriter, r *http.Request) {
	temp := template.Must(template.ParseFiles("add.html"))
	temp.Execute(w, nil)
}
func updateCall(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("_id")
	fmt.Println(id)
	temp := template.Must(template.ParseFiles("update.html"))
	temp.Execute(w, map[string]string{"ID": id})
}
func delete(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("_id")
	_, err := db.Query("delete from report where id=?", id)
	if err != nil {
		fmt.Println("Unable to Delete")
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
	fmt.Println("successfully Deleted")
	http.Redirect(w, r, "/", http.StatusSeeOther)

}
func update(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("_id")
	title := r.FormValue("_title")
	platform := r.FormValue("_platform")
	module := r.FormValue("_modules")
	desc := r.FormValue("_desc")
	buttons := r.FormValue("buttonValue")
	_, err := db.Query("update report set title=?, platform=?,modules=?, descriptions=?, tags=? where id=?", title, platform, module, desc, buttons, id)
	if err != nil {
		fmt.Println("Unable to Insert")
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
	fmt.Println("successfully Inserted")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
func vote(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("_id")
	row, err := db.Query("select votes from report where id=?", id)
	if err != nil {
		fmt.Println("Unable to Rise vote")
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
	var vc int
	row.Next()
	row.Scan(&vc)
	vc++
	_, err1 := db.Query("update report set votes=? where id=?", vc, id)
	if err1 != nil {
		fmt.Println("Unable to Rise vote")
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
	fmt.Println("Rised 1 vote")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
func main() {

	http.HandleFunc("/", send)
	http.HandleFunc("/add", add)
	http.HandleFunc("/addCall", addCall)
	http.HandleFunc("/updateCall", updateCall)
	http.HandleFunc("/update", update)
	http.HandleFunc("/vote", vote)
	http.HandleFunc("/delete", delete)
	http.ListenAndServe(":8080", nil)
}
