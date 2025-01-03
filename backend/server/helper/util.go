package helper

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

func SendResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func SendWebsocket(conn *websocket.Conn, data interface{}) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Println("Error marshalling JSON:", err)
		return err
	}

	err = conn.WriteMessage(websocket.TextMessage, jsonData)
	if err != nil {
		log.Println("Error sending message:", err)
		return err
	}
	fmt.Println("Message sent:", string(jsonData))
	return nil
}

func MultipartExtractJSON(w http.ResponseWriter, r *http.Request) (interface{}, error){
	var data interface{}

	jsonPart, _, err := r.FormFile("jsonData")
	if err == nil {
		jsonBytes, err := io.ReadAll(jsonPart)
		if err != nil {
			http.Error(w, "Error reading JSON data", http.StatusInternalServerError)
			// return nil, fmt.Errorf("error reading json data")
			return nil, nil
		}

		err = json.Unmarshal(jsonBytes, &data)
		if err != nil {
			http.Error(w, "Error parsing JSON data", http.StatusInternalServerError)
			return nil, nil
		}

		fmt.Printf("Received JSON metadata: %+v\n", data)
	} else {
		fmt.Println("No JSON data provided in request.")
	}

	return data, nil
}

func getClientIP(r *http.Request) string {
    host, _, err := net.SplitHostPort(r.RemoteAddr)
    if err != nil {
        return ""
    }
    return host
}

func GetClientData(r *http.Request) (host string, connID string) {
	host = getClientIP(r)
	connID = r.Header.Get("X-Connection-ID")
	
	return host, connID
}

func VerifyClient(host string, connID string, w http.ResponseWriter) bool {
	if host == "" {
		fmt.Println("Could not read ip")
		http.Error(w, "Could not read IP", http.StatusInternalServerError)
		return false
	}
	if connID == "" {
		fmt.Println("Invalid. No Conn-ID in request")
		http.Error(w, "", http.StatusUnauthorized)
		return false
	}

	return true
}

func MakeDir(path string) (completedPath string) {
	if strings.Contains(path, ".") {
		path = filepath.Dir(path)
	}
	_, err := os.Stat(path)
	if err != nil {
		os.MkdirAll(path, 0755)
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		return path
	}

	return absPath
}

func GetDirectoryContents(path string, typeIsDir bool, acceptType string) *[]string {
	var scanned []string
	objs, err := os.ReadDir(path)
	if err != nil {
		return nil
	}

	for _, obj := range objs {
		if typeIsDir {
			if obj.IsDir() {
				scanned = append(scanned, obj.Name())
			}
		} else {
			if acceptType == "" {
				scanned = append(scanned, obj.Name())	
			} else {
				if filepath.Ext(obj.Name()) == acceptType {
					scanned = append(scanned, obj.Name())	
				}
			}
		}
	}

	return &scanned
}

func ExtractPresets(config map[string]string, mainPath string) map[string]string {
	presetMappings := map[string]string{}

	for key, val := range config {
		if strings.Contains(key, "PRESET") && !strings.Contains(key, "#") {
			presetElems := strings.Split(key, "_")
			// fmt.Printf("Preset elems: %v\n", presetElems)
			if len(presetElems) == 2 {
				// fmt.Printf("Key: %v, Val: %v\n", key, val)
				staticPath := filepath.Join(mainPath, "STATIC")
				presetPath := filepath.Join(staticPath, val)
				// fmt.Printf("Preset path: %v\n", presetPath)
				_, err := os.Stat(presetPath)
				if err != nil {
					fmt.Printf("%v is not valid\n", presetPath)
				} else {
					presetMappings[val] = presetPath
				}
			}
		}
	}

	return presetMappings
}

func IsNewDay(t1 time.Time, t2 time.Time) bool {
	return t1.YearDay() != t2.YearDay()
}

