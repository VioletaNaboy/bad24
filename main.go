package main

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

func processFile(filePath string) string {
	f, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}
	lines := strings.Split(string(f), "\n")
	var numbers []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if len(line) > 0 {
			numbers = append(numbers, line)
		}
	}

	graph := buildGraph(numbers)

	var longestSeq []string
	memo := make(map[string][]string)

	for _, num := range numbers {
		visited := make(map[string]bool)
		seq := dfs(num, graph, visited, memo)
		if len(seq) > len(longestSeq) {
			longestSeq = seq
		}
	}

	result := ""
	for i, num := range longestSeq {
		if i == 0 {
			result += num
		} else {
			result += num[2:]
		}
	}

	fmt.Println("Longest sequence:", longestSeq)
	fmt.Println("Longest sequence:", result)

	return result
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tmpl := template.Must(template.ParseFiles("index.html"))
		tmpl.Execute(w, nil)
		return
	}

	if r.Method == http.MethodPost {
		file, header, err := r.FormFile("file")
		if err != nil {
			http.Error(w, "Error reading file", http.StatusBadRequest)
			return
		}
		defer file.Close()

		filePath := "./uploaded-" + header.Filename
		out, err := os.Create(filePath)
		if err != nil {
			http.Error(w, "Error saving file", http.StatusInternalServerError)
			return
		}
		defer out.Close()
		io.Copy(out, file)

		result := processFile(filePath)
		os.Remove(filePath)

		tmpl := template.Must(template.ParseFiles("index.html"))
		data := struct {
			Result string
		}{
			Result: result,
		}
		tmpl.Execute(w, data)
	}
}

func main() {
	http.HandleFunc("/", uploadHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Starting server on port %s\n", port)
	http.ListenAndServe(":"+port, nil)
	fmt.Println("Server is running on port", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
