package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var allowedExts = []string{".txt", ".pdf", ".word", ".docx"}

func isAllowedExtension(filename string) bool {
	ext := filepath.Ext(filename)
	for _, allowedExt := range allowedExts {
		if strings.ToLower(ext) == allowedExt {
			return true
		}
	}
	return false
}

func enableCors(w http.ResponseWriter) {
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Methods", "POST")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

func uploadFile(w http.ResponseWriter, r *http.Request) {
	enableCors(w)
	fmt.Printf("File upload endpoint in use\n")

	r.ParseMultipartForm(10 << 20)
	file, handler, err := r.FormFile("myfile")
	if err != nil {
		fmt.Printf("error retrieving file")
		fmt.Println(err)
		return
	}
	defer file.Close()

	// fmt.Printf("Uploaded file: %+v\n", handler.Filename)
	// fmt.Printf("File size: %+v\n", handler.Size)
	// fmt.Printf("MIME Header: %+v\n", handler.Header)

	okFile := isAllowedExtension(handler.Filename)
	if !okFile {
		log.Printf("Forbidden file upload attempt: %s\n", handler.Filename)
		fmt.Fprintf(w, "Forbidden file type")
		return
	}
	fmt.Printf("File check passed for %s [%d bytes]\n", handler.Filename, handler.Size)
	ext := filepath.Ext(handler.Filename)
	if ext == "" {
		ext = ".txt"
	}

	fileName := fmt.Sprintf("upload-*.%s", ext)

	tempFile, err := os.CreateTemp("uploaded", fileName)
	if err != nil {
		fmt.Println(err)
	}
	defer tempFile.Close()

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
	}
	tempFile.Write(fileBytes)
	fmt.Fprintf(w, "Successfully uploaded file\n")
}

func setupRoutes() {
	http.HandleFunc("/upload", uploadFile)
	http.ListenAndServe(":8080", nil)
}

func main() {
	fmt.Println("Hello world")
	setupRoutes()
}