package main

import (
	"encoding/json"
	"fmt"
	"github/muhammad-usaman13/processFile"
	"io"
	"net/http"
	"os"
	"strconv"
)

func main() {
	http.HandleFunc("/api/process-file/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			file, header, _ := r.FormFile("file")
			fmt.Println("File uploaded successfully", file)
			parts, _ := strconv.Atoi(r.FormValue("parts"))
			readFile, _ := io.ReadAll(file)

			returningValues := processFile.Process(readFile, parts)
			fmt.Println(returningValues)
			// create a folder uploads if it does not exist
			if _, err := os.Stat("uploads"); os.IsNotExist(err) {
				os.Mkdir("uploads", 0755)
			}
			// write the file in uploads folder
			print(header.Filename)
			print(os.WriteFile("uploads/"+header.Filename, readFile, 0644))
			json.NewEncoder(w).Encode(returningValues)
			return
		}
	})

	if err := http.ListenAndServe(":8000", nil); err != nil {
		fmt.Println("Server error:", err)
	}
}
