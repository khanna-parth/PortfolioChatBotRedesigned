package link

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"server/helper"
	"server/transfer"
	"strings"
)

type Linker struct {
	UID string
	ScanDir string
	DirPath string
}

func CreateLinker(uid string, sourceDir string) (*Linker) {
	_, err := os.Stat(sourceDir)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("CreateLinker failed. %s does not exist", sourceDir)
			return nil
		}
		fmt.Println("Error checking directory:", err)
		return nil
	}

	absPath, err := filepath.Abs(sourceDir)
	if err != nil {
		fmt.Printf("Error getting absolute path for directory: %s\n", err)
	}

	return &Linker{
		UID: uid,
		ScanDir: absPath,
	}
}

func (linker *Linker) SetPath(path string) (error) {
	helper.MakeDir(path)
	fmt.Printf("SetPath running. Path given: %s\n", path)
	linker.DirPath = path
	fmt.Printf("Set DirPath to %s\n", linker.DirPath)
	currPath, _ := os.Getwd()
	os.Chdir(linker.ScanDir)
	defer os.Chdir(currPath)
	_, err := os.Stat(path)
	if err != nil {
		err := os.Mkdir(path, 0755)
		if err != nil {
			return fmt.Errorf("could not make user folder")
		}
		return nil
	}
	return nil
}


func (linker *Linker) Scan() {
	os.Chdir(linker.ScanDir)
	cwd, _ := os.Getwd()
	files, _ := os.ReadDir(cwd)
	for _, file := range files {
		originalFile := file.Name()
		fmt.Printf("Found file: %v\n", originalFile)
		fileName := filepath.Base(originalFile)

		delimiterIndex := strings.Index(fileName, "_")
		if delimiterIndex != -1 && strings.HasSuffix(originalFile, ".pdf") {

			newFileName := fileName[delimiterIndex+1:]
			err := linker.moveFile(originalFile, newFileName)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
			}
		}
	}
}

func (linker *Linker) Delete(filename string) (error) {
	os.Chdir(linker.DirPath)
	files, readErr := os.ReadDir(linker.DirPath)
	if readErr != nil {
		readErr := fmt.Errorf("could not read user folder")
		log.Printf("%v\n", readErr)
		return readErr
	}

	for _, file := range files {
		if file.Name() == filename {
			err := os.Remove(file.Name())
			if err != nil {
				remErr := fmt.Errorf("found file but could not delete it")
				log.Println(remErr)
				return remErr
			}

			return nil
		}
	}
	missingErr := fmt.Errorf("could not find %s", filename)
	log.Printf("Linker Delete: %v\n", missingErr)
	return missingErr
}

func (linker *Linker) ListUserFiles(ref *transfer.DocumentList) {
	if linker.DirPath == "" {
		fmt.Println("Linker path not set");
		return
	}
	fmt.Printf("Going to dir path: %v\n", linker.DirPath)
	savePath, _ := os.Getwd()
	os.Chdir(linker.DirPath)
	defer os.Chdir(savePath)
	currPath, _ := os.Getwd()
	files, _ := os.ReadDir(currPath)
	fileNames := []string{}
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".pdf" {
		fileNames = append(fileNames, file.Name())
		}
	}
	
	ref.Docs = fileNames
}

// func (linker *Linker) Execute() {
// 	os.Get
// }

func (linker *Linker) moveFile(srcFile string, newFileName string) (error) {
	newFilePath := filepath.Join(linker.DirPath, newFileName)

	err := os.Rename(srcFile, newFilePath)
	if err != nil {
		return fmt.Errorf("error moving file %v", err)
	}

	fmt.Printf("Moved file to %v\n", newFilePath)
	return nil
}

func (linker *Linker) Build(script string) (buildPath string, err error) {
	dirPath := helper.MakeDir(linker.DirPath)
	fmt.Printf("Linker build reading files from DirPath: %s\n", dirPath)
	currPath, _ := os.Getwd()
	os.Chdir(dirPath)
	defer os.Chdir(currPath)
	files, err := os.ReadDir(dirPath)
	if err != nil {
		listError := fmt.Errorf("cannot build generator for linker with id: %s", linker.UID)
		fmt.Printf("%v\n", listError)
		return "", listError
	}

	fmt.Printf("%+v\n", files)

	for _, file := range files {
		if file.Name() == "generator.py" {
			existPath, err := filepath.Abs(file.Name())
			if err != nil {
				fmt.Printf("Linker build script output exists but abs path failed on file: %s\n", file.Name())
				return "", err
			}
			fmt.Printf("Linker script '%s' exists for id %s\n", existPath, linker.UID)
			return existPath, nil
		}
	}
	data, err := os.ReadFile(script)
	if err != nil {
		readErr := fmt.Errorf("could not read generator file for linker build with id %s", linker.UID)
		fmt.Printf("%v\n", readErr)
		return "", readErr
	}
	finalScriptPath := filepath.Join(linker.DirPath, "generator.py")
	writeOpError := os.WriteFile(finalScriptPath, data, 0644)
	if writeOpError != nil {
		writeErr := fmt.Errorf("could not write generator file for linker build with id %s", linker.UID)
		fmt.Printf("%v\n", writeErr)
		return "", writeErr
	}

	return finalScriptPath, nil
}

func (linker *Linker) Clean() (error) {
	fmt.Printf("Deleting directory: %v\n", linker.DirPath)
	os.RemoveAll(linker.DirPath)
	return nil
}
