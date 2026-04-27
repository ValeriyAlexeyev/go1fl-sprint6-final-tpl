package handlers

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/ValeriyAlexeyev/go1fl-sprint6-final/internal/service"
)

const maxUploadSize = 10 << 20 // 10MB

// IndexHandler отдаёт HTML
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	http.ServeFile(w, r, "index.html")
}

// UploadHandler обрабатывает файл
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseMultipartForm(maxUploadSize); err != nil {
		http.Error(w, "failed to parse form", http.StatusBadRequest)
		return
	}

	data, filename, err := readInput(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	result, err := service.Convert(string(data))
	if err != nil {
		if errors.Is(err, service.ErrEmptyInput) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, "conversion error", http.StatusInternalServerError)
		return
	}

	savedFile, err := saveResult(result, filename)
	if err != nil {
		http.Error(w, "failed to save result", http.StatusInternalServerError)
		return
	}

	// Ответ клиенту
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprintf(w, "Result:\n%s\n\nSaved as: %s\n", result, savedFile)
}

// --- helpers ---

func readInput(r *http.Request) ([]byte, string, error) {
	file, header, err := r.FormFile("uploadfile")
	if err == nil {
		defer file.Close()

		data, err := io.ReadAll(file)
		if err != nil {
			return nil, "", err
		}

		return data, header.Filename, nil
	}

	// fallback — читаем body
	data, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, "", err
	}

	return data, "input.txt", nil
}

func saveResult(result, originalName string) (string, error) {
	ext := filepath.Ext(originalName)
	if ext == "" {
		ext = ".txt"
	}

	filename := fmt.Sprintf("result_%d%s", time.Now().Unix(), ext)

	file, err := os.Create(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	if _, err := file.WriteString(result); err != nil {
		return "", err
	}

	return filename, nil
}
