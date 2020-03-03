package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"
	"html/template"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

var DSN = "root:1234@tcp(localhost:3306)/post?charset=utf8"

var tmpl = template.Must(template.New("Blog").ParseFiles("./public/index.html"))

type Server struct {
	db *sql.DB
}

type posts struct {
	ID int
	Date time.Duration
	Name string
	Post []post
}

type post struct {
	ID int
	PostID int
	Text string
}

func main(){
	db, err := sql.Open("mysql", DSN)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	s := Server{
		db: db,
	}

	router := mux.NewRouter()
	router.HandleFunc("/", s.homeHandle)
	router.HandleFunc("/posts/{id:[0-9]+}", s.postHandle)
}

func (s *Server) homeHandle(w http.ResponseWriter, r *http.Request){
	posts, err := getPosts(s.db)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := tmpl.ExecuteTemplate(w, "allposts", posts); err != nil {
		log.Println(err)
	}
}

func (s *Server) postHandle(w http.ResponseWriter, r *http.Request) {
	post, err := getPost(s.db, r.URL.Query().Get("id"))
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if err := tmpl.ExecuteTemplate(w, "post", post); err != nil {
		log.Println(err)
	}
}

func getPosts (db *sql.DB) ([]posts, error) {
	res := make([]posts, 0, 1)

	rows, err := db.Query("select * from posts.posts")
	if err != nil {
		return res, err
	}
	defer rows.Close()

	for rows.Next() {
		post := posts{}

		if err := rows.Scan(&post.ID, &post.Date, &post.Name); err != nil {
			log.Println(err)
			continue
		}

		res = append(res, post)
	}

	return res, nil
}

func getPost (db *sql.DB, id string) (posts, error) {
	posts := posts{}

	row := db.QueryRow(fmt.Sprintf("select * from posts.posts where posts.id = %v", id))
	err := row.Scan(&posts.ID, &posts.Date, &posts.Name)
	if err != nil {
		return posts, err
	}

	rows, err := db.Query(fmt.Sprintf("select * from posts.post WHERE post.post_id = %v", id))
	if err != nil {
		return posts, err
	}
	defer rows.Close()

	for rows.Next() {
		post := post{}

		err := rows.Scan(&post.ID, new(int), &post.Text)
		if err != nil {
			log.Println(err)
			continue
		}

		posts.Post = append(posts.Post, post)
	}

	return posts, nil
}