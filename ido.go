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

//type Vortaro struct {
//Name         string 'json:"name"'
//Tradukuro  []string 'json:"tradukuro"'
//}

//type Db struct {
//db *sql.DB
//}
//var db *sql.DB

func main() {
	//sql.Register("sqlite3_regexp",
	//&sqlite3.SQLiteDriver{
	//Extensions: []string{
	//"sqlite3_mod_regexp",
	//},
	//})

	db, err := sql.Open("sqlite3", "./ido.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/favicon.ico", nil)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "application/json")

		path := r.URL.Path // trouble characters: %#
		p := strings.Split(path, "/")

		query := p[1]
		var data string
		// return all data (for matching queries) by default
		if len(p) < 3 || p[2] == "" {
			data = "*"
		} else {
			data = p[2]
		}

		if ! validColumn(data) {
			w.WriteHeader(http.StatusBadRequest) //400
			return
		}
		

		//log.Println(query + " " + data)

		//_, err := db.Begin()
		//if err != nil {
		//log.Fatal(err)
		//}

		switch r.Method {
		case "GET":
			rows, err := db.Query("SELECT "+data+" FROM ido WHERE io = ?", query)
			if err != nil {
				log.Fatal(err)
			}
			defer rows.Close()

			for rows.Next() {
				var name string
				if err := rows.Scan(&name); err != nil {
					log.Fatal(err)
				}
				json.NewEncoder(w).Encode(name)
			}
			if err := rows.Err(); err != nil {
				log.Fatal(err)
			}

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

// validColumn checks if a column request is valid (to prevent SQL injections).
func validColumn(a string) bool {
	switch a {
	case
		"*", // all columns
		"af","am","an","ar","ast","ay","az","bar","be","bg","br","bs","ca","ce","ckb","cmn","co","cs","cv","cy","da","de","dv","el","en","eo","es","et","eu","fa","fi","fo","fr","fy","ga","gd","gl","gn","grc","gv","haw","he","hi","hr","ht","hu","hy","ia","id","ie","io","is","it","ja","jbo","jv","ka","ko","ku","kv","kw","la","lb","lbe","lt","lv","mdf","mg","mhr","mk","mn","ms","mt","mwl","myv","nl","nn","no","nov","nso","oc","om","pl","pms","ps","pt","qu","rm","ro","ru","rup","scn","sco","se","sh","sk","sl","so","sq","sr","sv","sw","th","tl","tok","tr","tt","udm","uk","vec","vi","vo","wa","wo","xh","yi","yo","zh","zu",
		"semantiko",
		"morfologio",
		"exemplaro",
		"sinomino",
		"antonimo",
		"kompundi",
		"kategorio":
		return true
	}
	return false
}
