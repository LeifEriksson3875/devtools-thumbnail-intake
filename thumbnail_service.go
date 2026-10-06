package main

import (
	"encoding/json"
	"net/http"
)

type thumbnailRequest struct {
	File     string `json:"file"`
	Filename string `json:"filename"`
	Width    int    `json:"width"`
	Height   int    `json:"height"`
}

func chooseSize(width, height int) (int, int) {
	if width <= 0 || height <= 0 {
		return 320, 180
	}
	if width > 640 {
		return 640, height * 640 / width
	}
	return width, height
}

func thumbnailHandler(client *ImageClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var in thumbnailRequest
		if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}
		width, height := chooseSize(in.Width, in.Height)
		id, err := client.Upload(in.File, in.Filename)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}
		result, err := client.Process(id, width, height)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"image": id, "width": width, "height": height, "thumbnail": json.RawMessage(result)})
	}
}

func main() { http.ListenAndServe(":8080", thumbnailHandler(NewImageClient())) }
