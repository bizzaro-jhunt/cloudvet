package main

import (
	"bytes"
	fmt "github.com/starkandwayne/goutils/ansi"
	"gopkg.in/redis.v5"
	"net/http"
	"os"
)

func init() {
	http.HandleFunc("/redis", func(w http.ResponseWriter, r *http.Request) {
		var expect, got string
		var b bytes.Buffer

		fmt.Fprintf(&b, "starting @M{redis} smoke tests...\n")
		Step(&b, "parsing VCAP_SERVICES env var to find our Redis endpoint")
		vcap, err := ParseVcap(os.Getenv("VCAP_SERVICES"))
		if err != nil {
			Final(w, b, err)
			return
		}
		OK(&b)

		Step(&b, "searching VCAP_SERVICES for our 'redis' service")
		var hostname, port, password string
		for _, service := range vcap {
			if !Tagged(service, "redis") {
				continue
			}
			hostname, err = Extract(service, "credentials", "host")
			if err != nil {
				Final(w, b, err)
				return
			}
			port, err = Extract(service, "credentials", "port")
			if err != nil {
				Final(w, b, err)
				return
			}
			password, err = Extract(service, "credentials", "password")
			if err != nil {
				Final(w, b, err)
				return
			}
			break
		}
		if hostname == "" {
			Final(w, b, fmt.Errorf("No service tagged 'redis' was found in VCAP_SERVICES"))
			return
		}
		OK(&b)

		Step(&b, "connecting to @C{%s:%s}", hostname, port)
		client := redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%s", hostname, port),
			Password: password,
		})
		_, err = client.Ping().Result()
		if err != nil {
			Final(w, b, err)
			return
		}
		OK(&b)

		key := "cloud.vet.redis.test.key"
		expect = "the first value"
		Step(&b, "storing a value")
		err = client.Set(key, expect, 0).Err()
		if err != nil {
			Final(w, b, err)
			return
		}
		OK(&b)

		Step(&b, "retrieving the stored value")
		got, err = client.Get(key).Result()
		if err != nil {
			Final(w, b, err)
			return
		}
		if got != expect {
			Final(w, b, fmt.Errorf("We wrote '%s' to the key-value store, but\n"+
				"got back '%s'", expect, got))
			return
		}
		OK(&b)

		Step(&b, "updating the stored a value")
		err = client.Set(key, expect, 0).Err()
		if err != nil {
			Final(w, b, err)
			return
		}
		OK(&b)

		Step(&b, "retrieving the updated value")
		got, err = client.Get(key).Result()
		if err != nil {
			Final(w, b, err)
			return
		}
		if got != expect {
			Final(w, b, fmt.Errorf("We wrote '%s' to the key-value store, but\n"+
				"got back '%s'", expect, got))
			return
		}
		OK(&b)

		fmt.Fprintf(&b, "\n\n@G{REDIS TESTS PASSED!}\n\n")
		Final(w, b, nil)
	})
}
