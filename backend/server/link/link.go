package link

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Link struct {
	client *http.Client
}

type User struct {
	UID string
	DocsStored int
	QueriesCount int
}

func CreateLink() *Link {
	return &Link{
		client: &http.Client{},
	}
}

func CreateUser() *User {
	return &User{}
}

// func (l *Link) RegisterID(user *User) {
// 	if user.uid == "" {
		
// 	}
// }

type DocumentRequest struct {
	UserID string `json:"userID"`
	Document string `json:"docPath"`
}

type DocumentResponse struct {
	Message string `json:"message"`
}

type DocumentListResponse struct {
	Docs[]string `json:"documents"`
}

func (l *Link) AddDocument(user *User, doc string) (string, error) {
	data := DocumentRequest{
		UserID: user.UID,
		Document: "resume.pdf",
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	
	url := "http://localhost:8000/add-doc"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := l.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		fmt.Printf("Status code was not 200: %v\n", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("internal error")
	}

	var response DocumentResponse
	unwrapError := json.Unmarshal(body, &response)
	if unwrapError != nil {
		return "", fmt.Errorf("internal error")
	}

	fmt.Printf("Response message: %v\n", response.Message)


	return response.Message, nil
}


func (l *Link) ListDocuments(user *User) (*[]string, error) {
	data := DocumentRequest{
		UserID: user.UID,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	
	url := "http://localhost:8000/list-docs"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := l.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		fmt.Printf("Status code was not 200: %v\n", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("internal error")
	}

	var response DocumentListResponse
	unwrapError := json.Unmarshal(body, &response)
	if unwrapError != nil {
		return nil, fmt.Errorf("internal error")
	}

	fmt.Printf("Response message: %v\n", response.Docs)


	return &response.Docs, nil
}