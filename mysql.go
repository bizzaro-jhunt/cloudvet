package main

import (
	"bytes"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	fmt "github.com/starkandwayne/goutils/ansi"
	"net/http"
	"os"
)

func init() {
	http.HandleFunc("/mysql", func(w http.ResponseWriter, r *http.Request) {
		var expect, got string
		var b bytes.Buffer

		fmt.Fprintf(&b, "starting @M{mysql} smoke tests...\n")
		Step(&b, "parsing VCAP_SERVICES env var to find our MySQL endpoint")
		vcap, err := ParseVcap(os.Getenv("VCAP_SERVICES"))
		if err != nil {
			Final(w, b, err)
			return
		}
		OK(&b)

		Step(&b, "searching VCAP_SERVICES for our 'mysql' service")
		var url, hostname, port, username, password, dbname string
		for _, service := range vcap {
			if !Tagged(service, "mysql") {
				continue
			}
			hostname, err = Extract(service, "credentials", "hostname")
			if err != nil {
				Final(w, b, err)
				return
			}
			port, err = Extract(service, "credentials", "port")
			if err != nil {
				Final(w, b, err)
				return
			}
			username, err = Extract(service, "credentials", "username")
			if err != nil {
				Final(w, b, err)
				return
			}
			password, err = Extract(service, "credentials", "password")
			if err != nil {
				Final(w, b, err)
				return
			}
			dbname, err = Extract(service, "credentials", "name")
			if err != nil {
				Final(w, b, err)
				return
			}
			url = fmt.Sprintf("%s:%s@(%s:%s)/%s", username, password, hostname, port, dbname)
			break
		}
		if url == "" {
			Final(w, b, fmt.Errorf("No service tagged 'mysql' was found in VCAP_SERVICES"))
			return
		}
		OK(&b)

		Step(&b, "connecting to @C{%s}", url)
		db, err := sql.Open("mysql", url)
		if err != nil {
			Final(w, b, err)
			return
		}
		OK(&b)

		Step(&b, "creating test database schema")
		_, err = db.Exec(`CREATE TABLE test_data (x TEXT)`)
		if err != nil {
			Final(w, b, err)
			return
		}
		OK(&b)

		expect = "the first value"
		Step(&b, "inserting a value into the database")
		_, err = db.Exec(`INSERT INTO test_data (x) VALUES (?)`, expect)
		if err != nil {
			Final(w, b, err)
			return
		}
		OK(&b)

		Step(&b, "retrieving the stored value from the database")
		row := db.QueryRow(`SELECT x FROM test_data LIMIT 1`)
		if row == nil {
			Final(w, b, fmt.Errorf("Unable to find any data after our INSERT statement"))
			return
		}
		err = row.Scan(&got)
		if err != nil {
			Final(w, b, err)
			return
		}
		if got != expect {
			Final(w, b, fmt.Errorf("We wrote '%s' to the database, but\n"+
				"got back '%s'", expect, got))
			return
		}
		OK(&b)

		expect = "an updated value"
		Step(&b, "updating a value in the database")
		_, err = db.Exec(`UPDATE test_data SET x = ?`, expect)
		if err != nil {
			Final(w, b, err)
			return
		}
		OK(&b)

		Step(&b, "retrieving the updated value from the database")
		row = db.QueryRow(`SELECT x FROM test_data LIMIT 1`)
		if row == nil {
			Final(w, b, fmt.Errorf("Unable to find any data after our UPDATE statement"))
			return
		}
		err = row.Scan(&got)
		if err != nil {
			Final(w, b, err)
			return
		}
		if got != expect {
			Final(w, b, fmt.Errorf("We wrote '%s' to the database, but\n"+
				"got back '%s'", expect, got))
			return
		}
		OK(&b)

		Step(&b, "dropping test database schema")
		_, err = db.Exec(`DROP TABLE test_data`)
		if err != nil {
			Final(w, b, err)
			return
		}
		OK(&b)

		fmt.Fprintf(&b, "\n\n@G{MYSQL TESTS PASSED!}\n\n")
		Final(w, b, nil)
	})
}
