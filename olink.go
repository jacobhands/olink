package olink

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func init() {
	file, e := ioutil.ReadFile("./config.json")
	if e != nil {
		fmt.Printf("File error: %v\n", e)
		c := Config{}
		c.Urls = make([]url, 20)
		conf, _ := json.MarshalIndent(c, "", " ")
		ioutil.WriteFile("./config.json", conf, 0644)
		os.Exit(1)
	}
	var c Config
	json.Unmarshal(file, &c)

	http.ListenAndServe(":8080", &handler{c})
}

type handler struct{ Config }

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		h.redirectShortUrl(w, r)
	}
}

func (h *handler) redirectShortUrl(w http.ResponseWriter, r *http.Request) {
	var url string
	fmt.Println("Path: " + r.URL.Path)
	for i := range h.Urls {
		if h.Urls[i].Key == r.URL.Path[1:] {
			url = h.Urls[i].Value
			break
		}
	}
	if url != "" {
		http.Redirect(w, r, url, 302)
		fmt.Println("Redirected: " + url)
	} else {
		fmt.Println("Not found: " + r.URL.Path)
		http.Error(w, "Not found.", 404)
	}
}

type Config struct {
	Urls []url
}

type url struct {
	Key   string
	Value string
}

// Exists reports whether the named file or directory exists.
func Exists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}
