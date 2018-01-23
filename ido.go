/* interkonekto serves a RESTful API to define and translate words in the Ido language.

    /<query>/<data>

query:

    Interpreted as a regex if preceeded by a tilde

data:

    Any |-dilimited list of:
*/
package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

//type Db struct {
//db *sql.DB
//}
var db *sql.DB

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	query := r.URL.Path // trouble: %#
	//_, err := db.Begin()
	//if err != nil {
	//log.Fatal(err)
	//}

	//var result sql.Result

	request := strings.Split(query, "/")[1]
	//    switch request := strings.Split(query, "/")[1]; request {
	//case "krear": //create
	//statement, err := db.Prepare("INSERT INTO ido (io) VALUES (?)")
	//if err != nil {
	//log.Fatal(err)
	//}
	//statement.Exec("en", query)

	//case "lektar": //read
	stmt, err := db.Prepare("SELECT en FROM ido WHERE io = ?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	var name string
	err = stmt.QueryRow(request).Scan(&name)
	if err != nil {
		log.Fatal(err)
	}
	//            fmt.Println(name)
	json.NewEncoder(w).Encode(name)

	//result, err = stmt.Exec("linguo")
	//if err != nil {
	//log.Fatal(err)
	//}

	//case "efacar": //delete
	//statement, _ := db.Prepare("DELETE FROM ido WHERE io='?'")
	//statement.Exec("en", query)
	//
	//case "skribar": //write / update
	//statement, _ := db.Prepare("UPDATE ido SET en='?' WHERE io='?'")
	//statement.Exec("en", query)
	//
	//default:
	//json.NewEncoder(w).Encode(query)
	//json.NewEncoder(w).Encode(request)
	//json.NewEncoder(w).Encode(result)
	//            log.Fatal("unknown request: ", request)
	//}
}

func main() {
	http.HandleFunc("/", handler)

	err := http.ListenAndServeTLS(":4443", "/etc/letsencrypt/live/linguo.io/fullchain.pem", "/etc/letsencrypt/live/linguo.io/privkey.pem", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

	db, err = sql.Open("sqlite3", "./ido.sql")
	if err != nil {
		log.Fatal(err)
	}
	//    defer db.Close()

}
