package link

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

type Executor struct {
	UID string
	Script string
	IsRunning bool
	PID int
	APIKey string

}

func NewExecutor(uid string) *Executor {
	return &Executor{
		UID: uid,
	}
}

type ChainResponse struct {
	Prompt string `json:"prompt"`
	Response string `json:"response"`
	Elapsed float64 `json:"elapsed"`
}

func (ex *Executor) ToggleUse() {
	ex.IsRunning = !ex.IsRunning
}

func (ex *Executor) Query(prompt string) (*ChainResponse, error) {
	log.Printf("executor chain prompt started with query '%s' for user %s at path %s\n", prompt, ex.UID, ex.Script)
	
	ex.IsRunning = true
	defer ex.ToggleUse()

	scriptPath := filepath.Dir(ex.Script)
	cmd := exec.Command("python3", ex.Script, "--query", prompt, "--dir", scriptPath, "--key", ex.APIKey)

	var output bytes.Buffer
	cmd.Stdout = &output

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		fmt.Println("STDOUT:")
		fmt.Println(output.String())

		fmt.Println("STDERR:")
		fmt.Println(stderr.String())
		log.Printf("error running executor script for uid %s: %v\n", ex.UID, err)
		return nil, err
	}

	// fmt.Println("STDOUT:")
	// fmt.Println(output.String())
	// fmt.Println("STDERR:")
	// fmt.Println(stderr.String())

	var chainOutput ChainResponse
	err = json.Unmarshal(output.Bytes(), &chainOutput)
	if err != nil {
		log.Printf("error parsing executor script for uid %s: %v\n", ex.UID, err)
		return nil, err
	}

	log.Printf("Executor chain successful. Response: %s\n", chainOutput.Response)
	return &chainOutput, nil
}

func (ex *Executor) Index(userDirPath string) error {
	currPath, err := os.Getwd()
	if err != nil {
		wdErr := fmt.Errorf("executor could not get starting working directory under uid: %s", ex.UID)
		log.Print(wdErr)
		return wdErr
	}
	fmt.Printf("Starting dir: %s\n", currPath)
	fmt.Printf("UserDirPath: %s\n", userDirPath)
	os.Chdir(userDirPath)
	defer os.Chdir(currPath)
	files, err := os.ReadDir(userDirPath)
	if err != nil {
		readErr := fmt.Errorf("executor could not read files for indexing under uid dir: %s", ex.UID)
		log.Print(readErr)
		return readErr
	}

	for _, file := range files {
		if file.Type().IsDir() {
			filePath, _ := filepath.Abs(file.Name())
			log.Printf("SRemoving stored chroma files under directory %s for user %s\n", filePath, ex.UID)
			os.RemoveAll(filePath)
		}
		if strings.HasSuffix(file.Name(), ".sqlite3") {
			os.Remove(file.Name())
			fmt.Printf("Simulated removing stored chroma sqlite file %s for user %s\n", file.Name(), ex.UID)
		}
	}
	_, failed := ex.Query("test")
	if failed != nil {
		return failed
	}
	return nil
}

func (ex *Executor) Kill() error {
	kill := exec.Command("kill", strconv.Itoa(ex.PID))
	err := kill.Run()
	if err != nil {
		killError := fmt.Errorf("Executor with PID %d could not be killed", ex.PID)
		fmt.Printf("%v\n", killError)
		return killError
	}
	return nil
}

