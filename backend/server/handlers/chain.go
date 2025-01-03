package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"server/helper"
	"server/link"
	"server/transfer"
	"strings"
)

var allowedExts = []string{".txt", ".pdf", ".word", ".docx"}

func PresentationPromptHandler(w http.ResponseWriter, r *http.Request, connStore *link.ConnectionStore) {
	srcIP, connID := helper.GetClientData(r)
	if !helper.VerifyClient(srcIP, connID, w) {
		return
	}

	if !connStore.AvailableInPool(5) {
		http.Error(w, "Sorry, I'm busy with other requests. Try again later", http.StatusTooManyRequests)
		return		
	}

	validPreset, presetPath := connStore.IsPreset(connID)
	if !validPreset {
		http.Error(w, "Endpoint only for presentation through parthkhanna.me or for presets", http.StatusForbidden)
		return	
	}

	var promptRequest transfer.UserPromptRequest
	decoder := json.NewDecoder(r.Body)
	decodeErr := decoder.Decode(&promptRequest)
	if decodeErr != nil {
		http.Error(w, "Invalid data", http.StatusBadRequest)
		return
	}

	promptValid, limit := connStore.IsPromptValid(&promptRequest.Prompt)
	if !promptValid {
		http.Error(w, fmt.Sprintf("Prompt too large. Limit to %d", limit), http.StatusRequestEntityTooLarge)
		return
	}

	if connStore.IsRequestMaxReached(srcIP) {
		http.Error(w, "Sorry, I'm too tired to continue. Wake me up tomorrow", http.StatusTooManyRequests)
		return
	}

	demoChat := link.NewChatDemo(connID, srcIP, connStore)
	fmt.Printf("Chat: %+v\n", demoChat)
	
	demoChat.Lock.Lock()

	if demoChat.InUse {
		log.Printf("Chat connection for %s is already busy.\n", srcIP)
		http.Error(w, "Request already pending", http.StatusTooManyRequests)
		return	
	}
	demoChat.InUse = true
	
	demoChat.Lock.Unlock()

	fmt.Println("Changed chat to in use")

	defer func() {
		demoChat.Lock.Lock()
		demoChat.InUse = false
		demoChat.Lock.Unlock()
		fmt.Println("Changed chat to free")
	}()
	// defer demoChat.ToggleUse()

	scriptPath := filepath.Join(presetPath, "generator.py")

	demoChat.Executor.APIKey = connStore.Config["OPENAI_API_KEY"]
	demoChat.Executor.Script = scriptPath

	resp, err := demoChat.Executor.Query(promptRequest.Prompt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	connStore.AddRequestCount(srcIP)
	helper.SendResponse(w, resp)
}

func UserPromptHandler(w http.ResponseWriter, r *http.Request, connStore *link.ConnectionStore) {
	srcIP, connID := helper.GetClientData(r)
	if !helper.VerifyClient(srcIP, connID, w) {
		return
	}

	if !connStore.MatchConnection(connID, srcIP) {
		fmt.Printf("No registed connection ID: %s and source IP: %s\n", connID, srcIP)
		http.Error(w, "Connection not recognized, try reloading", http.StatusUnauthorized)
		return
	}

	if connStore.IsRequestMaxReached(srcIP) {
		http.Error(w, "Request limit reached. Your limit will be reset tomorrow", http.StatusTooManyRequests)
		return
	}

	var promptRequest transfer.UserPromptRequest
	decoder := json.NewDecoder(r.Body)
	decodeErr := decoder.Decode(&promptRequest)
	if decodeErr != nil {
		http.Error(w, "Invalid data", http.StatusBadRequest)
		return
	}

	promptValid, limit := connStore.IsPromptValid(&promptRequest.Prompt)
	if !promptValid {
		http.Error(w, fmt.Sprintf("Prompt too large. Limit to %d", limit), http.StatusRequestEntityTooLarge)
		return
	}

	chat := connStore.GetConnection(connID)
	if chat == nil {
		fmt.Printf("DocumentListHandler failing. No chat")
		http.Error(w, "No chat resource found", http.StatusInternalServerError)
		return
	}
	fmt.Printf("DocumentListHandler Chat Found: %+v\n", chat)
	if chat.InUse {
		http.Error(w, "Previous request pending", http.StatusAccepted)
		return
	}
	chat.InUse = true
	defer chat.ToggleUse()

	userDir := filepath.Join(connStore.GetUploadsPath(), connID)
	if chat.Linker.DirPath == "" {
		chat.Linker.SetPath(userDir)
	}

	buildPath, buildErr := chat.Linker.Build(connStore.Script)
	if buildErr != nil {
		http.Error(w, "could not prepare generation script", http.StatusInternalServerError)
		return
	}
	log.Printf("UserPromptHandler buildPath: %s\n", buildPath)

	chat.Executor.Script = buildPath
	chat.Executor.APIKey = connStore.Config["OPENAI_API_KEY"]

	resp, err := chat.Executor.Query(promptRequest.Prompt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	connStore.AddRequestCount(srcIP)
	helper.SendResponse(w, resp)
}


func DocumentListHandler(w http.ResponseWriter, r *http.Request, connStore *link.ConnectionStore) {
	srcIP, connID := helper.GetClientData(r)
	if !helper.VerifyClient(srcIP, connID, w) {
		return
	}

	isPreset, presetPath := connStore.IsPreset(connID)
	if isPreset {
		files := helper.GetDirectoryContents(presetPath, false, ".pdf")
		docs := &transfer.DocumentListResponse{
			Docs: *files,
		}
		helper.SendResponse(w, docs)
		return
	}

	if !connStore.MatchConnection(connID, srcIP) {
		fmt.Printf("No registed connection ID: %s and source IP: %s\n", connID, srcIP)
		http.Error(w, "", http.StatusUnauthorized)
		return
	}

	chat := connStore.GetConnection(connID)
	if chat == nil {
		fmt.Printf("DocumentListHandler failing. No chat")
		http.Error(w, "No chat resource found", http.StatusInternalServerError)
		return
	}
	fmt.Printf("DocumentListHandler Chat Found: %+v\n", chat)
	if chat.InUse {
		http.Error(w, "Previous request pending", http.StatusAccepted)
		return
	}
	chat.InUse = true
	defer chat.ToggleUse()

	docs := transfer.CreateDocumentList()

	chat.Linker.ListUserFiles(docs)
	if docs == nil {
		fmt.Println("No docs recieved. Directory has not been made yet")
		helper.SendResponse(w, docs)
		return
	}
	fmt.Printf("Docs retrieved: %v\n", docs)

	// connStore.AddRequestCount(srcIP)
	helper.SendResponse(w, docs)
}

func DocumentUploadHandler(w http.ResponseWriter, r *http.Request, connStore *link.ConnectionStore) {
	fmt.Printf("File upload endpoint in use\n")

	srcIP, connID := helper.GetClientData(r)
	if !helper.VerifyClient(srcIP, connID, w) {
		return
	}

	isPreset, _ := connStore.IsPreset(connID)
	if isPreset {
		http.Error(w, fmt.Sprintf("%s under static preset files cannot be modified", connID), http.StatusForbidden)
		return
	}

	if connStore.IsRequestMaxReached(srcIP) {
		http.Error(w, "Request limit reached. Your limit will be reset tomorrow", http.StatusTooManyRequests)
		return
	}

	if !connStore.MatchConnection(connID, srcIP) {
		http.Error(w, "Connection not recognized, try reloading", http.StatusUnauthorized)
		return
	}

	chat := connStore.GetConnection(connID)
	if chat == nil {
		http.Error(w, "Error getting chat connection", http.StatusInternalServerError)
		return
	}

	if chat.InUse {
		http.Error(w, "Previous request pending", http.StatusAccepted)
		return
	}
	chat.InUse = true
	defer chat.ToggleUse()

	userDir := filepath.Join(connStore.GetUploadsPath(), connID)
	if chat.Linker.DirPath == "" {
		chat.Linker.SetPath(userDir)
	}

	helper.MakeDir(userDir)
	
	// Change first value, will shift 20 bits to right.
	r.ParseMultipartForm(2 << 20)
	file, handler, err := r.FormFile("myfile")
	if err != nil {
		fmt.Printf("error retrieving file")
		fmt.Println(err)
		return
	}
	defer file.Close()

	// fmt.Printf("Uploaded file: %+v\n", handler.Filename)
	fmt.Printf("File size: %+v\n", handler.Size)
	// fmt.Printf("MIME Header: %+v\n", handler.Header)

	if handler.Size > 20000000 {
		log.Printf("%s's document upload was denied of extreme file size %d\n", connID, handler.Size)
		http.Error(w, "File must not be bigger than 2MB", http.StatusRequestEntityTooLarge)
		return
	}

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

	fileName := fmt.Sprintf("%s_%s", connID, handler.Filename)

	fileSavePath := filepath.Join(connStore.UploadsDir, fileName)

	tempFile, err := os.Create(fileSavePath)
	if err != nil {
		fmt.Println(err)
	}
	defer tempFile.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		fmt.Println(err)
	}
	tempFile.Write(fileBytes)
	fmt.Printf("File saved to: %s\n", tempFile.Name())
	fmt.Fprintf(w, "Successfully uploaded file")

	userFolderPath := filepath.Join(chat.Linker.ScanDir, connID)
	chat.Linker.SetPath(userFolderPath)
	chat.Linker.Scan()

	buildPath, buildErr := chat.Linker.Build(connStore.Script)
	if buildErr != nil {
		http.Error(w, "could not prepare generation script", http.StatusInternalServerError)
		return
	}
	log.Printf("DocumentUploadHandler build path: %s\n", buildPath)

	chat.Executor.Script = buildPath
	chat.Executor.APIKey = connStore.Config["OPENAI_API_KEY"]

	indexingError := chat.Executor.Index(filepath.Dir(buildPath))
	if indexingError != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// _, testError := chat.Executor.Query("test")
	// if testError != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }

	connStore.AddRequestCount(srcIP)
	helper.SendResponse(w, "Added")
}

func DocumentDeleteHandler(w http.ResponseWriter, r *http.Request, connStore *link.ConnectionStore) {
	log.Printf("Document Delete handler in use\n")

	srcIP, connID := helper.GetClientData(r)
	if !helper.VerifyClient(srcIP, connID, w) {
		return
	}

	isPreset, _ := connStore.IsPreset(connID)
	if isPreset {
		http.Error(w, fmt.Sprintf("%s under static preset files cannot be modified", connID), http.StatusForbidden)
		return
	}

	if connStore.IsRequestMaxReached(srcIP) {
		http.Error(w, "Request limit reached. Your limit will be reset tomorrow", http.StatusTooManyRequests)
		return
	}

	if !connStore.MatchConnection(connID, srcIP) {
		http.Error(w, "Connection not recognized, try reloading", http.StatusUnauthorized)
		return
	}

	chat := connStore.GetConnection(connID)
	if chat == nil {
		http.Error(w, "Error getting chat connection", http.StatusInternalServerError)
		return
	}

	if chat.InUse {
		http.Error(w, "Previous request pending", http.StatusAccepted)
		return
	}
	chat.InUse = true
	defer chat.ToggleUse()

	var deleteRequest transfer.DocumentDeleteRequest
	decoder := json.NewDecoder(r.Body)
	decodeErr := decoder.Decode(&deleteRequest)
	if decodeErr != nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	fmt.Printf("DeleteHandler paths: %v | %v\n", connStore.GetUploadsPath(), connID)

	userDir := filepath.Join(connStore.GetUploadsPath(), connID)
	chat.Linker.SetPath(userDir)

	buildPath, err := chat.Linker.Build(connStore.Script)
	if err != nil {
		http.Error(w, "could not prepare generation script", http.StatusInternalServerError)
		return
	}

	chat.Executor.Script = buildPath
	chat.Executor.APIKey = connStore.Config["OPENAI_API_KEY"]

	chat.Linker.Delete(deleteRequest.Document)
	indexErr := chat.Executor.Index(userDir)
	if indexErr != nil {
		http.Error(w, indexErr.Error(), http.StatusInternalServerError)
		return
	}

	connStore.AddRequestCount(srcIP)
	helper.SendResponse(w, "Deleted")
}

func isAllowedExtension(filename string, ) bool {
	ext := filepath.Ext(filename)
	for _, allowedExt := range allowedExts {
		if strings.ToLower(ext) == allowedExt {
			return true
		}
	}
	return false
}