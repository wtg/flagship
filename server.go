package main

import (
    "net/http"
    "log"
    "html/template"
    // "database/sql"
	// _ "github.com/lib/pq"
)

type Page struct {
    Categories []string
    Files []string
}

func index(w http.ResponseWriter, r *http.Request) {
    // fmt.Fprintf(w, "hello world")
    t, err := template.ParseFiles("templates/index.html")
    if err != nil {
        // TODO: 500 handler
        log.Fatal(err)
    }

    p := Page{
        Categories: []string{"Category", "Category", "Category", "Category", "Category", "Category", "Category", "Category", "Category", "Category", "Category"},
        Files: []string{"File", "File", "File", "File", "File"},
    }
    t.Execute(w, p)
}

func main() {
    // _, err := sql.Open("postgres", "user=root dbname=flagship")
	// if err != nil {
	// 	log.Fatal(err)
	// }

    http.HandleFunc("/", index)
    http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

    url := "127.0.0.1:8080"
    log.Print("listening on ", url)
    log.Fatal(http.ListenAndServe(url, nil))
}
