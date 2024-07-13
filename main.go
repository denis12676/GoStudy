package main

import (
	"net/http"
)

func main() {
	// Обработка статических файлов
	fs := http.FileServer(http.Dir("./static"))

	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Маршрутизация для главной страницы
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/index.html")
	})

	// Запуск сервера на порту 8080
	http.ListenAndServe(":8080", nil)
}
