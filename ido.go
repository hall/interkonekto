/* interkonekto serves a RESTful API to define and translate words in the Ido language. */
package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

var allColumns []string

func main() {
	//sql.Register("sqlite3_regexp",
	//&sqlite3.SQLiteDriver{
	//Extensions: []string{
	//"sqlite3_mod_regexp",
	//},
	//})

	db, err := sql.Open("sqlite3", "./ido.db")
	panicOnErr(err)
	defer db.Close()
	err = db.Ping()
	panicOnErr(err)

	rows, err := db.Query("SELECT * FROM ido WHERE io = 'linguo'")
	panicOnErr(err)
	allColumns, err = rows.Columns()
	panicOnErr(err)

	http.HandleFunc("/favicon.ico", nil)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "application/json")
		output := json.NewEncoder(w)

		path := strings.Split(r.URL.Path, "/") // trouble characters: %#

		query := path[1]
		var data string
		// return all data (for matching queries) if no specific set is requested
		if len(path) < 3 || path[2] == "" {
			data = "*"
		} else {
			data = path[2]
		}

		if !validColumn(data) {
			w.WriteHeader(http.StatusBadRequest) //400
			return
		}

		switch r.Method {
		case "GET":
			rows, err := db.Query("SELECT "+data+" FROM ido WHERE io = ?", query)
			panicOnErr(err)
			defer rows.Close()

			columns, err := rows.Columns()
			panicOnErr(err)

			colNum := len(columns)

			cols := make([]interface{}, colNum)
			result := make([]string, colNum)
			for i, _ := range result {
				cols[i] = &result[i]
			}

			for rows.Next() {
				err = rows.Scan(cols...)
				panicOnErr(err)
				output.Encode(result)
			}
			err = rows.Err()
			panicOnErr(err)

			//w.WriteHeader(http.StatusOK) //200
			//w.WriteHeader(http.StatusNotFound) //404
			//w.WriteHeader(http.StatusBadRequest) //400

		case "PUT": //write / update
		//statement, _ := db.Prepare("UPDATE ido SET en='?' WHERE io='?'")
		//statement.Exec("en", query)

		case "POST": //create
		//statement, err := db.Prepare("INSERT INTO ido (io) VALUES (?)")
		//if err != nil {
		//log.Fatal(err)
		//}
		//statement.Exec("en", query)

		case "DELETE": //delete
			//statement, _ := db.Prepare("DELETE FROM ido WHERE io='?'")
			//statement.Exec("en", query)

			//default:
			//json.NewEncoder(w).Encode(query)
			//json.NewEncoder(w).Encode(request)
			//json.NewEncoder(w).Encode(result)
			//            log.Fatal("unknown request: ", request)
		}

	})
	err = http.ListenAndServeTLS(":4443", "/etc/letsencrypt/live/linguo.io/fullchain.pem", "/etc/letsencrypt/live/linguo.io/privkey.pem", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}

// validColumn checks if column requests are valid (to prevent SQL injections).
func validColumn(a string) bool {
	var found bool
	for _, c := range strings.Split(a, ",") {
		found = false
		for _, b := range allColumns {
			if b == c {
				found = true
				continue
			}
		}
		if !found {
			return false
		}
	}
	return true
}

func panicOnErr(err error) {
	if err != nil {
		panic(err.Error())
	}
}
