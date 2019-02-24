/*interkonekto serves a RESTful API to define and translate words in the ido language. */
package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	_ "github.com/lib/pq"
)

var db *sql.DB
var datumi []string

func main() {

	db = DB()
	defer db.Close()

	//http.HandleFunc("/favicon.ico", nil)
	http.HandleFunc("/", indexHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}

func DB() *sql.DB {

	dbURI := fmt.Sprintf("user=postgres password=%s host=/cloudsql/bryton:us-central1:vortaro dbname=postgres", os.Getenv("DB_PASSWORD"))
	db, err := sql.Open("postgres", dbURI)
	panicOnErr(err)

	rows, err := db.Query("SELECT * FROM vortaro WHERE io = 'lingu.o'")
	panicOnErr(err)
	datumi, err = rows.Columns()
	panicOnErr(err)

	return db

}

func indexHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, PUT, OPTIONS")
	w.Header().Set("Access-Control-Max-Age", "86400")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	rendimento := json.NewEncoder(w)

	path := strings.Split(r.URL.Path, "/") // trouble characters: %#

	switch path[1] {
	case "":
		rendimento.Encode("https://gitlab.com/hall/interkonekto")
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
		linguo = "io" // search with ido by default
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
		rows, err := db.Query("SELECT "+datumo+" FROM vortaro WHERE REPLACE("+linguo+",'.','') = $1", demando)
		panicOnErr(err)
		defer rows.Close()

		cols, err := rows.Columns()
		panicOnErr(err)

		colNum := len(cols)

		dest := make([]interface{}, colNum)
		raw := make([][]byte, colNum)
		result := make(map[string]interface{}, colNum)
		for i := range raw {
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

		//statement, err := db.Prepare("UPDATE vortaro SET "+data+" = ?")
		//if err != nil {
		//log.Fatal(err)
		//}
		//statement.Exec(nova)
	//	rendimento.Encode("UPDATE vortaro SET '" + datumo + "' = '" + nova + "' WHERE io = " + demando + ";")

	case "POST":
	//statement, _ := db.Prepare("UPDATE vortaro SET en='?' WHERE io='?'")
	//statement.Exec("en", query)

	case "DELETE": //delete
		//statement, _ := db.Prepare("DELETE FROM vortaro WHERE io='?'")
		//statement.Exec("en", query)

		//default:
		//json.NewEncoder(w).Encode(query)
		//json.NewEncoder(w).Encode(request)
		//json.NewEncoder(w).Encode(result)
		//            log.Fatal("unknown request: ", request)
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
