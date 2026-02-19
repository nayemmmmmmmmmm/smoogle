package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// Mock apiConfig (add any fields if needed)
type testAPIConfig struct{}

func TestHandlerValidateChirp_Valid(t *testing.T) {
	apiCfg := &apiConfig{}

	body := `{"body":"Hello world"}`
	req := httptest.NewRequest(http.MethodPost, "/api/validate_chirp", strings.NewReader(body))
	w := httptest.NewRecorder()

	apiCfg.handlerValidateChirp(w, req)

	resp := w.Result()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}

	var result validationResponse
	json.NewDecoder(resp.Body).Decode(&result)

	if !result.Valid {
		t.Fatalf("expected valid=true, got false")
	}
}

func TestHandlerValidateChirp_TooLong(t *testing.T) {
	apiCfg := &apiConfig{}

	longBody := strings.Repeat("a", 141)
	body := `{"body":"` + longBody + `"}`

	req := httptest.NewRequest(http.MethodPost, "/api/validate_chirp", strings.NewReader(body))
	w := httptest.NewRecorder()

	apiCfg.handlerValidateChirp(w, req)

	resp := w.Result()

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", resp.StatusCode)
	}

	var result errorResponse
	json.NewDecoder(resp.Body).Decode(&result)

	if result.Error != "Chirp is too long" {
		t.Fatalf("unexpected error message: %s", result.Error)
	}
}

func TestHandlerValidateChirp_InvalidJSON(t *testing.T) {
	apiCfg := &apiConfig{}

	req := httptest.NewRequest(http.MethodPost, "/api/validate_chirp", bytes.NewBuffer([]byte("{invalid json")))
	w := httptest.NewRecorder()

	apiCfg.handlerValidateChirp(w, req)

	resp := w.Result()

	if resp.StatusCode != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", resp.StatusCode)
	}
}

func TestHandlerValidateChirp_ProfanityFiltered(t *testing.T) {
	apiCfg := &apiConfig{}

	// Adjust this depending on what badWordReplacement does
	body := `{"body":"badword"}`

	req := httptest.NewRequest(http.MethodPost, "/api/validate_chirp", strings.NewReader(body))
	w := httptest.NewRecorder()

	apiCfg.handlerValidateChirp(w, req)

	resp := w.Result()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}

	var result validationResponse
	json.NewDecoder(resp.Body).Decode(&result)

	// If badWordReplacement changes the text,
	// cleanedMessage != params.Body → Valid should be false
	if result.Valid {
		t.Fatalf("expected valid=false when profanity is replaced")
	}
}
