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
	columns, err := rows.Columns()
	panicOnErr(err)

	http.HandleFunc("/favicon.ico", nil)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		output := json.NewEncoder(w)

		path := strings.Split(r.URL.Path, "/") // trouble characters: %#

		switch path[1] {
		case "":
			output.Encode("https://github.com/linguo-io/api")
			return
		case "*":
			output.Encode(columns)
			return
		}

		query := path[1]

		var data string
		// return all data (for matching queries) if no specific set is requested
		if len(path) < 3 {
			data = "*"
		} else {
			data = path[2]
		}

		if !validColumn(data, columns) {
			w.WriteHeader(http.StatusBadRequest) //400
			return
		}

		switch r.Method {
		case "GET":
			splitQuery := strings.Split(query, ":")
			if (len(splitQuery) == 1) {
				splitQuery[0] = "io"
				splitQuery = append(splitQuery, query)
			}

			if !validColumn(splitQuery[0], columns) {
				w.WriteHeader(http.StatusBadRequest) //400
				return
			}


			rows, err := db.Query("SELECT "+data+" FROM ido WHERE replace("+splitQuery[0]+",'.','') = ?", splitQuery[1])
			panicOnErr(err)
			defer rows.Close()

			cols, err := rows.Columns()
			panicOnErr(err)

			colNum := len(cols)

			dest := make([]interface{}, colNum)
			raw := make([][]byte, colNum)
			result := make(map[string]interface{}, colNum)
			for i, _ := range raw {
				dest[i] = &raw[i]
			}

			for rows.Next() {
				err = rows.Scan(dest...)
				panicOnErr(err)
				for i, r := range raw {
					if r == nil {
						result[cols[i]] = ""
					} else {
						entries := strings.Split(string(r), "\n")
						if len(entries) > 1 {
							result[cols[i]] = strings.Split(string(r), "\n")
						} else {
							result[cols[i]] = string(r)
						}
					}

				}
				master := make(map[string]map[string]interface{}, colNum)
				master[splitQuery[1]] = result
				output.Encode(master)
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
func validColumn(a string, columns []string) bool {
	var found bool
	for _, c := range strings.Split(a, ",") {
		found = false
		for _, b := range columns {
			if b == c || c == "*" {
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
