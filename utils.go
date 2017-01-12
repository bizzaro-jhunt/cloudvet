package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"math/rand"
	"strings"
	"time"

	fmt "github.com/starkandwayne/goutils/ansi"
	"io"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func ParseVcap(s string) ([]interface{}, error) {
	var v interface{}
	fmt.Printf("parsing vcap:\n%s\n", s)
	err := json.Unmarshal([]byte(s), &v)
	if err != nil {
		return nil, err
	}

	if m, ok := v.(map[string]interface{}); ok {
		ll := make([]interface{}, 0)
		for _, vv := range m {
			if ss := vv.([]interface{}); ok {
				for _, s := range ss {
					ll = append(ll, s)
				}
			}
		}
		return ll, nil
	}

	return nil, fmt.Errorf("VCAP_SERVICES is corrupt:\n%s\n", s)
}

func Tagged(service interface{}, want string) bool {
	if svc, ok := service.(map[string]interface{}); ok {
		if tags, ok := svc["tags"]; ok {
			if ll, ok := tags.([]interface{}); ok {
				for _, tag := range ll {
					if fmt.Sprintf("%v", tag) == want {
						return true
					}
				}
			}
		}
	}
	return false
}

func extract(service interface{}, keys []string, orig string) (string, error) {
	if m, ok := service.(map[string]interface{}); ok {
		if v, ok := m[keys[0]]; ok {
			if len(keys) > 1 {
				return extract(v, keys[1:], orig)
			}
			switch v.(type) {
			case map[interface{}]interface{}:
				return "", fmt.Errorf("%s is a map, not a scalar", orig)
			case []interface{}:
				return "", fmt.Errorf("%s is a list, not a scalar", orig)
			}
			return fmt.Sprintf("%v", v), nil
		}
	}
	return "", fmt.Errorf("%s not found in VCAP_SERVICES", orig)
}
func Extract(service interface{}, keys ...string) (string, error) {
	return extract(service, keys, strings.Join(keys, "."))
}

func Step(out io.Writer, s string, args ...interface{}) {
	fmt.Fprintf(out, s+"... ", args...)
}

func Info(out io.Writer, s string, args ...interface{}) {
	fmt.Fprintf(out, s+"\n", args...)
}

func Final(w http.ResponseWriter, b bytes.Buffer, err error) {
	w.Header().Set("Content-type", "text/plain")
	if err == nil {
		w.WriteHeader(200)
	} else {
		w.WriteHeader(500)
	}
	fmt.Fprintf(w, b.String())
	if err != nil {
		fmt.Fprintf(w, "@R{FAILED}\n")
		fmt.Fprintf(w, "@R{=========================================================}\n")
		fmt.Fprintf(w, "@R{%s}\n", err)
		fmt.Fprintf(w, "@R{=========================================================}\n\n\n")
	}
}

func OK(out io.Writer) {
	fmt.Fprintf(out, "@G{OK}\n")
}

const randos = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandomString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = randos[rand.Intn(len(randos))]
	}
	return string(b)
}
