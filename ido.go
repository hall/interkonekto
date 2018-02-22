/* interkonekto serves a RESTful API to define and translate words in the Ido language. */
package main

import (
	"database/sql"
	"encoding/json"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
	"strings"
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

	rows, err := db.Query("SELECT * FROM ido WHERE io = 'lingu.o'")
	panicOnErr(err)
	datumi, err := rows.Columns()
	panicOnErr(err)

	http.HandleFunc("/favicon.ico", nil)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, PUT, OPTIONS")
		w.Header().Set("Access-Control-Max-Age", "86400")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		rendimento := json.NewEncoder(w)

		path := strings.Split(r.URL.Path, "/") // trouble characters: %#

		switch path[1] {
		case "":
			rendimento.Encode("https://github.com/linguo-io/api")
			return
		case "*":
			rendimento.Encode(datumi)
			return
		}

		demando := path[1]

		var datumo string
		var linguo string

		splitQuery := strings.Split(demando, ":")
		if len(splitQuery) == 1 {
			linguo = "io" // search with Ido by default
		} else {
			linguo = splitQuery[0]
			demando = splitQuery[1]
		}

		if !validColumn(linguo, datumi) {
			w.WriteHeader(http.StatusBadRequest) //400
			return
		}

		// return all data (for matching queries) if no specific set is requested
		if len(path) < 3 || path[2] == "" {
			datumo = "*"
		} else {
			datumo = path[2]
		}

		if !validColumn(datumo, datumi) {
			w.WriteHeader(http.StatusBadRequest) //400
			return
		}

		switch r.Method {
		case "GET":
			rows, err := db.Query("SELECT "+datumo+" FROM ido WHERE replace("+linguo+",'.','') = ?", demando)
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
				master[demando] = result
				rendimento.Encode(master)
			}
			err = rows.Err()
			panicOnErr(err)

			//w.WriteHeader(http.StatusOK) //200
			//w.WriteHeader(http.StatusNotFound) //404
			//w.WriteHeader(http.StatusBadRequest) //400

		case "PUT":
			if datumo == "*" {
				return
			}
			err := r.ParseForm()
			panicOnErr(err)
			//nova := r.PostFormValue("nova")
			//if (len(nova) < 1) { return }

			//statement, err := db.Prepare("UPDATE ido SET "+data+" = ?")
			//if err != nil {
			//log.Fatal(err)
			//}
			//statement.Exec(nova)
		//	rendimento.Encode("UPDATE ido SET '" + datumo + "' = '" + nova + "' WHERE io = " + demando + ";")

		case "POST":
		//statement, _ := db.Prepare("UPDATE ido SET en='?' WHERE io='?'")
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
	err = http.ListenAndServeTLS(":4443", "cert.pem", "key.pem", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}

// validColumn checks if column requests are valid (to prevent SQL injections).
func validColumn(datumo string, datumi []string) bool {
	var found bool
	for _, c := range strings.Split(datumo, ",") {
		found = false
		for _, b := range datumi {
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
