package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestProcessRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/v1/image/process" {
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		var body struct {
			Image map[string]string `json:"image"`
			Ops   []struct {
				Op     string         `json:"op"`
				Params map[string]int `json:"params"`
			} `json:"ops"`
			Format string `json:"format"`
			Store  bool   `json:"store"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		if body.Image["image_id"] != "pim_test" || len(body.Ops) != 1 || body.Ops[0].Op != "resize" || body.Ops[0].Params["width"] != 320 || body.Ops[0].Params["height"] != 180 || body.Format != "webp" || !body.Store {
			t.Errorf("unexpected process payload: %+v", body)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"ok":true,"data":{}}`))
	}))
	defer server.Close()
	client := &ImageClient{baseURL: server.URL, httpClient: server.Client()}
	if _, err := client.Process("pim_test", 320, 180); err != nil {
		t.Fatal(err)
	}
}
