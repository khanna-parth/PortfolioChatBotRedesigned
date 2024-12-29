package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"server/link"
	"server/middleware"
	"strings"
	"time"

	"github.com/google/uuid"
)

var allowedExts = []string{".txt", ".pdf", ".word", ".docx"}
var cookieStore = make(map[int]string)

func isAllowedExtension(filename string) bool {
	ext := filepath.Ext(filename)
	for _, allowedExt := range allowedExts {
		if strings.ToLower(ext) == allowedExt {
			return true
		}
	}
	return false
}

// func enableCors(w http.ResponseWriter) {
//     w.Header().Set("Access-Control-Allow-Origin", "*")
//     w.Header().Set("Access-Control-Allow-Methods", "POST")
//     w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
// }

func uploadFile(w http.ResponseWriter, r *http.Request) {
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


func main() {
	fmt.Println("Hello world")
	mux := http.NewServeMux()

	registerHandler := middleware.ApplyMiddlewares(http.HandlerFunc(registerUser), middleware.LoggingMiddleware, middleware.CORSMiddleware)

	mux.Handle("/register", registerHandler)
	mux.HandleFunc("/set", SetCookieHandler)
	mux.HandleFunc("/get", GetCookieHandler)
	mux.HandleFunc("/add-doc", DocumentModificationHandler)
	mux.HandleFunc("/list-docs", DocumentListHandler)
	mux.HandleFunc("/test", handler)

	http.ListenAndServe(":3000", mux)
}

func handler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	for i := 0; i < 20; i++ {
		select {
		case <-ctx.Done():
			http.Error(w, "Request was canceled", http.StatusRequestTimeout)
			return
		default:
			time.Sleep(500 * time.Millisecond)
			fmt.Printf("Processing... %d\n", i)
		}
	}

	// Send the response after completing the work
	fmt.Fprintf(w, "Processing completed.")
}

func DocumentModificationHandler(w http.ResponseWriter, r *http.Request) {
	l := link.CreateLink()
	user := link.CreateUser()
	user.UID = "pkhanna1"
	resp, err := l.AddDocument(user, "")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, resp)
}

func DocumentListHandler(w http.ResponseWriter, r *http.Request) {
	l := link.CreateLink()
	user := link.CreateUser()
	user.UID = "pkhanna"
	docs, err := l.ListDocuments(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	SendResponse(w, docs)
}

func SendResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func registerUser(w http.ResponseWriter, r *http.Request) {
	cookie := http.Cookie{
		Name: "regCookie",
		Value: uuid.New().String(),
		Path: "/",
		MaxAge: 3600,
		HttpOnly: true,
		Secure: true,
		SameSite: http.SameSiteLaxMode,
	}

	http.SetCookie(w, &cookie)

}

func GetCookieHandler(w http.ResponseWriter, r *http.Request) {
	cookie := GetCookie(w, r)
	fmt.Println("Get cookie handler hit")
	fmt.Printf("Cookie: %v\n", cookie)
	if cookie != "" {
		for key, val := range cookieStore {
			fmt.Printf("Key: %v, val: %v\n", key, val)
			if val == cookie {
				fmt.Printf("Welcome back user: %v\n", key)
			}
		}
	}
}

func GetCookie(w http.ResponseWriter, r *http.Request) string {
	cookie, err := r.Cookie("regCookie")
	if err != nil {
		switch {
			case errors.Is(err, http.ErrNoCookie):
				http.Error(w, "cookie not found", http.StatusBadRequest)
			default:
				log.Println(err)
				http.Error(w, "server error", http.StatusInternalServerError)
		}

		return ""
	}

	return cookie.Value
}

func SetCookieHandler(w http.ResponseWriter, r *http.Request) {
	SetCookie(w, uuid.New())
	fmt.Println("Cookie handler set completed")
}

func SetCookie(w http.ResponseWriter, uid uuid.UUID) {
	cookie := http.Cookie{
		Name: "regCookie",
		Value: uid.String(),
		Path: "/",
		MaxAge: 3600,
		HttpOnly: true,
		Secure: true,
		SameSite: http.SameSiteLaxMode,
	}

	fmt.Printf("Set cookie %+v\n", cookie)

	cookieStore[len(cookieStore)+1] = cookie.Value

	http.SetCookie(w, &cookie)
}

