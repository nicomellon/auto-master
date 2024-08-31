package main

import (
	"io/ioutil"
	"net/http/httptest"
	"testing"
)

func TestUploadHandler(t *testing.T) {
	req := httptest.NewRequest("POST", "/upload", nil)
	w := httptest.NewRecorder()
	uploadHandler(w, req)
	resp := w.Result()
	if resp.StatusCode != 400 {
		t.Fail()
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	if string(data) != "Error parsing request body" {
		t.Errorf("expected Error parsing request body got %v", string(data))
	}
}
