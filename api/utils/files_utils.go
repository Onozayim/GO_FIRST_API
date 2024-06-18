package utils

import (
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"runtime"
)

var (
	_, b, _, _ = runtime.Caller(0)
	rootPath   = filepath.Join(filepath.Dir(b), "../..")
)

func fileNameWithoutExtension(filename string) string {
	return filename[:len(filename)-len(filepath.Ext(filename))]
}

func SaveFile(file multipart.File, handler *multipart.FileHeader, folder string, delete bool) error {
	defer file.Close()

	folderPath := fmt.Sprintf("%s%s", rootPath, "/api/public")

	if folder != "" {
		folderPath = folderPath + "/" + folder
	}

	fileName := fmt.Sprintf("upload-%s-*%s", fileNameWithoutExtension(handler.Filename), filepath.Ext(handler.Filename))

	_, err := os.Stat(folderPath)

	if !os.IsNotExist(err) && delete && folder != "" {
		os.RemoveAll(folderPath)
	}

	if err := os.Mkdir(folderPath, os.ModePerm); err != nil {
		log.Println(err.Error())
		return fmt.Errorf("algo inexperado ah ocurrido")
	}

	tempFile, err := os.CreateTemp(folderPath, fileName)

	if err != nil {
		log.Println(err.Error())
		return fmt.Errorf("algo inexperado ah ocurrido")
	}

	defer tempFile.Close()

	filebytes, err := io.ReadAll(file)
	if err != nil {
		log.Println(err.Error())
		return fmt.Errorf("algo inexperado ah ocurrido")
	}

	tempFile.Write(filebytes)
	return nil
}
