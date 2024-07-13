package main

import (
	"html/template"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

type Post struct {
	ID      int
	Title   string
	Content string
	Image   string
}

var posts = []Post{}
var tmpl = template.Must(template.ParseGlob("templates/*.html"))

func main() {
	// Обработка статических файлов
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Обработка изображений
	http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("./images"))))

	// Маршрутизация
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/post/", postHandler)
	http.HandleFunc("/new", newHandler)

	// Запуск сервера
	http.ListenAndServe(":8080", nil)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "index.html", posts)
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Path[len("/post/"):])
	if err != nil || id < 0 || id >= len(posts) {
		http.NotFound(w, r)
		return
	}
	tmpl.ExecuteTemplate(w, "post.html", posts[id])
}

func newHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		title := r.FormValue("title")
		content := r.FormValue("content")

		var imagePath string
		file, handler, err := r.FormFile("image")
		if err == nil {
			defer file.Close()
			imagePath = filepath.Join("images", handler.Filename)
			f, err := os.OpenFile(imagePath, os.O_WRONLY|os.O_CREATE, 0666)
			if err != nil {
				http.Error(w, "Unable to save image", http.StatusInternalServerError)
				return
			}
			defer f.Close()
			_, err = io.Copy(f, file)
			if err != nil {
				http.Error(w, "Unable to save image", http.StatusInternalServerError)
				return
			}
		}

		imagePath = "/" + imagePath // добавляем слэш для корректного отображения пути
		posts = append(posts, Post{ID: len(posts), Title: title, Content: content, Image: imagePath})
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	tmpl.ExecuteTemplate(w, "new.html", nil)
}
