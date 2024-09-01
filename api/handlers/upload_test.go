package handlers

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/go-audio/audio"
	"github.com/go-audio/wav"
)

func TestUploadWithBadHttpMethod(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/upload", nil)
	w := httptest.NewRecorder()

	UploadHandler(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Fatalf("Expected status code %d, got %d", http.StatusMethodNotAllowed, resp.StatusCode)
	}
}

func TestUploadWithoutRequestBody(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/upload", nil)
	w := httptest.NewRecorder()

	UploadHandler(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("Expected status code %d, got %d", http.StatusBadRequest, resp.StatusCode)
	}
}

func TestUploadWithInvalidFileFormat(t *testing.T) {
	var buffer bytes.Buffer
	writer := multipart.NewWriter(&buffer)
	fileWriter, err := writer.CreateFormFile("file", "mockfile.txt")
	if err != nil {
		t.Fatalf("Failed to create form file: %v", err)
	}
	mockFileContent := []byte("This is some mock file content.")
	_, err = io.Copy(fileWriter, bytes.NewReader(mockFileContent))
	if err != nil {
		t.Fatalf("Failed to copy mock file content: %v", err)
	}
	writer.Close()
	req := httptest.NewRequest(http.MethodPost, "/upload", &buffer)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()

	UploadHandler(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("Expected status code %d, got %d", http.StatusBadRequest, resp.StatusCode)
	}
}

func TestUploadWithValidAudioFile(t *testing.T) {
	var buffer bytes.Buffer
	writer := multipart.NewWriter(&buffer)
	fileWriter, err := writer.CreateFormFile("file", "mockfile.wav")
	if err != nil {
		t.Fatalf("Failed to create form file: %v", err)
	}
	tmpFile, err := os.CreateTemp("", "mockfile*.wav")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	encoder := wav.NewEncoder(tmpFile, 44100, 16, 1, 1)
	audioBuffer := &audio.IntBuffer{
		Format: &audio.Format{
			SampleRate:  44100,
			NumChannels: 1,
		},
		Data:           make([]int, 44100), // 1 second of silence
		SourceBitDepth: 16,
	}
	if err := encoder.Write(audioBuffer); err != nil {
		t.Fatalf("Failed to write audio buffer: %v", err)
	}
	if err := encoder.Close(); err != nil {
		t.Fatalf("Failed to close encoder: %v", err)
	}
	if _, err := tmpFile.Seek(0, io.SeekStart); err != nil {
		t.Fatalf("Failed to seek in temp file: %v", err)
	}
	_, err = io.Copy(fileWriter, tmpFile)
	if err != nil {
		t.Fatalf("Failed to copy WAV file content: %v", err)
	}

	if err := encoder.Write(audioBuffer); err != nil {
		t.Fatalf("Failed to write audio buffer: %v", err)
	}
	if err := encoder.Close(); err != nil {
		t.Fatalf("Failed to close encoder: %v", err)
	}
	writer.Close()
	req := httptest.NewRequest(http.MethodPost, "/upload", &buffer)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()

	UploadHandler(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("Expected status code %d, got %d", http.StatusCreated, resp.StatusCode)
	}

}
