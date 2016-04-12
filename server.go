package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
)

type Category struct {
	Id          int64
	Name        string
	Description string
	Private     bool
	Writable    bool
	Featured    bool
	ParentID    int64
	GroupID     int64
	CreatedAt   string
	UpdatedAt   string
}

type Document struct {
	Id          int64
	Title       string
	Description string
	Private     bool
	Writable    bool
	CategoryID  int64
	UserID      int64
	CreatedAt   string
	UpdatedAt   string
}

type Page struct {
	Categories []Category
}

//
// Database
//

var db *sql.DB

func getDB() *sql.DB {
	db, err := sql.Open("postgres", "postgres://vagrant:vagrant@localhost/vagrant")
	if err != nil {
		log.Fatal(err)
	}
	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}

	db.Exec(`
        CREATE TABLE IF NOT EXISTS documents (
            id bigserial PRIMARY KEY,
            title text,
            description text,
            private boolean,
            writable boolean,
            category_id bigint,
            user_id bigint,
            created_at timestamp NOT NULL DEFAULT now(),
            updated_at timestamp
        );`)

	return db
}

func (doc *Document) DBInsert() {
	err := db.QueryRow("INSERT INTO documents (title, description, private, writable, category_id, user_id) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id", doc.Title, doc.Description, doc.Private, doc.Writable, doc.CategoryID, doc.UserID).Scan(&doc.Id)
	if err != nil {
		log.Fatal(err)
	}
}

//
// Routes
//

func index(w http.ResponseWriter, r *http.Request) {
	// fmt.Fprintf(w, "hello world")
	t, err := template.ParseFiles("templates/index.html")
	if err != nil {
		// TODO: 500 handler
		log.Fatal(err)
	}

	p := Page{
		Categories: []Category{
			Category{Id: 1, Name: "Test"},
			Category{Id: 2, Name: "Test"},
			Category{Id: 3, Name: "Test"},
			Category{Id: 4, Name: "Test"},
			Category{Id: 5, Name: "Test"},
			Category{Id: 6, Name: "Test"},
			Category{Id: 7, Name: "Test"},
			Category{Id: 8, Name: "Test"},
			Category{Id: 9, Name: "Test"},
			Category{Id: 10, Name: "Test"},
			Category{Id: 11, Name: "Test"},
		},
	}
	t.Execute(w, p)
}

func upload(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t, err := template.ParseFiles("templates/upload.html")
		if err != nil {
			// TODO: 500 handler
			log.Fatal(err)
		}

		p := Page{
			Categories: []Category{
				Category{Id: 1, Name: "Test"},
				Category{Id: 2, Name: "Test"},
				Category{Id: 3, Name: "Test"},
				Category{Id: 4, Name: "Test"},
				Category{Id: 5, Name: "Test"},
				Category{Id: 6, Name: "Test"},
				Category{Id: 7, Name: "Test"},
				Category{Id: 8, Name: "Test"},
				Category{Id: 9, Name: "Test"},
				Category{Id: 10, Name: "Test"},
				Category{Id: 11, Name: "Test"},
			},
		}
		t.Execute(w, p)
	} else if r.Method == "POST" {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = r.ParseMultipartForm(10000) // 10 MB in memory
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		files := r.MultipartForm.File["files"]
		if len(files) == 0 {
			fmt.Fprintf(w, "no files")
			return
		}
		for _, fileHeader := range files {
			file, err := fileHeader.Open()
			defer file.Close()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			doc := Document{
				Title:       r.Form.Get("title"),
				Description: r.Form.Get("description"),
			}
			doc.Private, err = strconv.ParseBool(r.FormValue("private"))
			doc.Writable, err = strconv.ParseBool(r.FormValue("writeable"))
			doc.CategoryID, err = strconv.ParseInt(r.FormValue("category_id"), 10, 64)
			doc.UserID, err = strconv.ParseInt(r.FormValue("user_id"), 10, 64)
			doc.DBInsert()

			// create uploaded file
			dest, err := os.Create("files/" + strconv.FormatInt(doc.Id, 10))
			defer dest.Close()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// copy to uploaded file
			if _, err := io.Copy(dest, file); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
		http.Redirect(w, r, "/", http.StatusFound)
	} else {
		// TODO: return unsupported method message
		log.Fatal("unsupported method")
	}
}

func main() {
	db = getDB()

	http.HandleFunc("/", index)
	http.HandleFunc("/upload", upload)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	url := "127.0.0.1:8080"
	log.Print("listening on ", url)
	log.Fatal(http.ListenAndServe(url, nil))
}
