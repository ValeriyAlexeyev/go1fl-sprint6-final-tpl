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

// IndexHandler отдаёт index.html
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	http.ServeFile(w, r, "index.html")
}

// UploadHandler принимает данные и конвертирует их
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	defer r.Body.Close()

	data, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "read body error", http.StatusInternalServerError)
		return
	}

	result, err := service.Convert(string(data))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	filename := time.Now().UTC().String() + filepath.Ext("result.txt")

	file, err := os.Create(filename)
	if err != nil {
		http.Error(w, "create file error", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	_, err = file.WriteString(result)
	if err != nil {
		http.Error(w, "write file error", http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, result)
}
