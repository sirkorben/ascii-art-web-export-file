package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"text/template"
)

func startHandler(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/ascii-art" && r.URL.Path != "/" {
		http.Error(w, "Oops! 404 not found.", http.StatusNotFound)
		return
	}
	temp, err := template.ParseFiles("templates/index.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error 500 - Internal server error!"))
		return
	}
	err = temp.Execute(w, "")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error 500 - Internal server error!"))
		return
	}
}

func secondHandler(w http.ResponseWriter, r *http.Request) {

	temp, err := template.ParseFiles("templates/ascii-art.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error 500 - Internal server error!"))
		return
	}
	r.ParseForm()
	err = temp.Execute(w, toAscii(r.FormValue("name"), r.FormValue("submit")))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error 500 - Internal server error!"))
		return
	}
}

func handleRequest() {

	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets/"))))
	http.Handle("/templates/", http.StripPrefix("/templates/", http.FileServer(http.Dir("./templates/"))))
	http.Handle("/scripts/", http.StripPrefix("/scripts/", http.FileServer(http.Dir("./scripts/"))))
	http.HandleFunc("/", startHandler)
	http.HandleFunc("/ascii-art", secondHandler)
	http.ListenAndServe(":8080", nil)
}

func main() {

	fmt.Println("Visit http://localhost:8080 for the result")
	handleRequest()
}

func toAscii(inpt string, bannr string) string {

	data, err := os.ReadFile("fonts/" + bannr)
	if err != nil {
		panic(err)
	}
	lines := strings.Split(string(data), "\n")
	input := inpt
	words := strings.Split(string(input), "\\n")
	result := ""

	if inpt == "" {
		return "\n\n\nYou need to enter something!\n\n\n"
	}
	for _, val := range inpt {
		if val < 32 || val > 126 {
			return "\n\n\nOnly English letters and printable symbols are accepted!\n\n\n"
		}
	}
	for i := 0; i < len(words); i++ {
		word := string(words[i])
		var store [][]string
		for _, char := range word {
			ascii := int(char)
			lineNbr := (ascii-32)*9 + 1
			store = append(store, lines[lineNbr:lineNbr+9])
		}
		for i := 0; i < 8; i++ {
			for _, value := range store {
				result += value[i]
			}
			result += "\n"
		}
	}
	return result
}
