package main
import (
	"fmt"
	"os"
	"net/http"
	"bytes"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type MongoDoc struct {
	Key string
	Value string
}

func init() {
	http.HandleFunc("/mongo", func(w http.ResponseWriter, r *http.Request) {
		var b bytes.Buffer

		fmt.Fprintf(&b, "starting @M{mongo} smoke tests...\n")
		Step(&b, "parsing VCAP_SERVICES env var to find our MongoDB endpoint")
		vcap, err := ParseVcap(os.Getenv("VCAP_SERVICES"))
		if err != nil {
			Final(w, b, err)
			return
		}
		OK(&b)

		Step(&b, "searching VCAP_SERVICES for our 'mongo' service")
		var url string
		for _, service := range vcap {
			if !Tagged(service, "mongo") {
				continue
			}
			url, err = Extract(service, "credentials", "url")
			if err != nil {
				Final(w, b, err)
				return
			}
			break
		}
		if url == "" {
			Final(w, b, fmt.Errorf("No service tagged 'mongo' was found in VCAP_SERVICES"))
			return
		}
		OK(&b)

		Step(&b, "connecting to @C{%s}", url)
		session, err := mgo.Dial(url)
		if err != nil {
			Final(w, b, err)
			return
		}
		defer session.Close()
		OK(&b)

		dbName := "vet-db-" + RandomString(16)
		cName := "vet-col-" + RandomString(16)
		Info(&b, "using %s/%s", dbName, cName)
		c := session.DB(dbName).C(cName)

		key := "key-" + RandomString(4)
		val := "val-" + RandomString(44)

		Step(&b, "inserting a tracer document")
		err = c.Insert(&MongoDoc{
			Key: key,
			Value: val,
		})
		if err != nil {
			Final(w, b, err)
			return
		}
		OK(&b)

		Step(&b, "retrieving tracer document")
		var d MongoDoc
		err = c.Find(bson.M{"key":key}).One(&d)
		if err != nil {
			Final(w, b, err)
			return
		}
		if d.Value != val {
			err = fmt.Errorf("got '%s' from mongo, but we expected '%s'", d.Value, val)
			Final(w, b, err)
			return
		}
		OK(&b)

		Step(&b, "removing tracer document")
		err = c.Remove(bson.M{"key":key})
		if err != nil {
			Final(w, b, err)
			return
		}
		d.Value = "(none)"
		err = c.Find(bson.M{"key":key}).One(&r)
		if err == nil {
			err = fmt.Errorf("was able to retrieve our document after removing it")
			Final(w, b, err)
			return
		}
		if d.Value == val {
			err = fmt.Errorf("was able to retrieve our document after removing it")
			Final(w, b, err)
			return
		}
		OK(&b)

		fmt.Fprintf(&b, "\n\n@G{MONGODB TESTS PASSED!}\n\n")
		Final(w, b, nil)
	})
}
