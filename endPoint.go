package main

import (
	"fmt"
	"io"
	"net/http"
	"os"

	// "strconv"
	"time"
)

func main() {
	http.HandleFunc("/api/process-file/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			file, header, _ := r.FormFile("file")
			fmt.Println("File uploaded successfully", file)
			// parts, _ := strconv.Atoi(r.FormValue("parts"))
			readFile, _ := io.ReadAll(file)

			// create a folder uploads if it does not exist
			if _, err := os.Stat("uploads"); os.IsNotExist(err) {
				os.Mkdir("uploads", 0755)
			}
			// write the file in uploads folder
			print(header.Filename)
			print(os.WriteFile("uploads/"+header.Filename, readFile, 0644))
			// json.NewEncoder(w).Encode(map[string]any{"words": words, "vowels": vowels, "alphabets": alphabets, "spaces": spaces, "timeTaken": timeTaken})
			return
		}
	})

	if err := http.ListenAndServe(":8000", nil); err != nil {
		fmt.Println("Server error:", err)
	}
}

func conc(str []byte, a chan int, v chan int, s chan int) {
	var alphabets, vowels, spaces int
	for i := 0; i < len(str); i++ {
		if (str[i] <= 95 && str[i] > 65) || (str[i] <= 122 && str[i] > 97) {
			alphabets += 1
			if str[i] == 97 || str[i] == 101 || str[i] == 105 || str[i] == 111 || str[i] == 117 {
				vowels += 1
			}
		} else if str[i] == 32 {
			spaces += 1
		}
	}
	v <- vowels
	a <- alphabets
	s <- spaces
}
func processFile(str []byte, noOfPortions int) (int, int, int, int, float64) {
	start := float64(time.Now().UnixNano()) / 1e9
	var words, vowels, alphabets, spaces int = 0, 0, 0, 0
	len := len(str)
	a := make(chan int)
	v := make(chan int)
	s := make(chan int)
	for i := 1; i <= noOfPortions; i++ {
		go conc(str[(i-1)/noOfPortions*len:(i/noOfPortions)*len], a, v, s)
	}

	for i := 1; i <= noOfPortions; i++ {
		vowels += <-v
		alphabets += <-a
		spaces += <-s
	}

	words = spaces + 1
	end := float64(time.Now().UnixNano()) / 1e9

	return words, vowels, alphabets, spaces, end - start
}
