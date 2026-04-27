package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/ValeriyAlexeyev/go1fl-sprint6-final/internal/service"
)

const maxUploadSize = 10 << 20 // 10 MB

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	http.ServeFile(w, r, "index.html")
}

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseMultipartForm(maxUploadSize)
	if err != nil {
		http.Error(w, "parse form error", http.StatusInternalServerError)
		return
	}

	var file io.ReadCloser
	var filename string

	// Берём первый загруженный файл независимо от имени поля
	for _, headers := range r.MultipartForm.File {
		if len(headers) > 0 {
			f, err := headers[0].Open()
			if err != nil {
				http.Error(w, "file open error", http.StatusInternalServerError)
				return
			}
			file = f
			filename = headers[0].Filename
			break
		}
	}

	if file == nil {
		http.Error(w, "file not found", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "read file error", http.StatusInternalServerError)
		return
	}

	result, err := service.Convert(string(data))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ext := filepath.Ext(filename)
	if ext == "" {
		ext = ".txt"
	}

	outName := time.Now().UTC().String() + ext

	outFile, err := os.Create(outName)
	if err != nil {
		http.Error(w, "create file error", http.StatusInternalServerError)
		return
	}
	defer outFile.Close()

	_, err = outFile.WriteString(result)
	if err != nil {
		http.Error(w, "write file error", http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, result)
}
	
