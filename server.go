package main

import (
    "net/http"
    "log"
    "html/template"
    // "database/sql"
	// _ "github.com/lib/pq"
)

type Category struct {
    Id int
    Name string
    Description string
    IsPrivate bool
    IsWritable bool
    IsFeatured bool
    ParentId int
    GroupId int
    CreatedAt string
    UpdatedAt string
}

type Document struct {
    Id int
    Title string
    Description string
    IsPrivate bool
    IsWriteable bool
    CategoryId int
    UserId int
    CreatedAt string
    UpdatedAt string
}

type Page struct {
    Categories []Category
}

func index(w http.ResponseWriter, r *http.Request) {
    // fmt.Fprintf(w, "hello world")
    t, err := template.ParseFiles("templates/index.html")
    if err != nil {
        // TODO: 500 handler
        log.Fatal(err)
    }

    p := Page{
        Categories: []Category{
            Category{Id:  1, Name:"Test"},
            Category{Id:  2, Name:"Test"},
            Category{Id:  3, Name:"Test"},
            Category{Id:  4, Name:"Test"},
            Category{Id:  5, Name:"Test"},
            Category{Id:  6, Name:"Test"},
            Category{Id:  7, Name:"Test"},
            Category{Id:  8, Name:"Test"},
            Category{Id:  9, Name:"Test"},
            Category{Id: 10, Name:"Test"},
            Category{Id: 11, Name:"Test"},
        },
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
