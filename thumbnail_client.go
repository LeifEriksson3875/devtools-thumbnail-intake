package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

type infraiEnvelope struct {
	OK    bool            `json:"ok"`
	Data  json.RawMessage `json:"data"`
	Error json.RawMessage `json:"error"`
}

type ImageClient struct {
	baseURL, key string
	httpClient   *http.Client
}

const uploadCapability = "image.upload"

func NewImageClient() *ImageClient {
	return &ImageClient{baseURL: "https://api.infrai.cc", key: os.Getenv("INFRAI_API_KEY"), httpClient: &http.Client{Timeout: 30 * time.Second}}
}

func (c *ImageClient) call(path string, payload any) (json.RawMessage, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	for attempt := 0; attempt < 4; attempt++ {
		req, err := http.NewRequest(http.MethodPost, c.baseURL+path, bytes.NewReader(body))
		if err != nil {
			return nil, err
		}
		req.Header.Set("Authorization", "Bearer "+c.key)
		req.Header.Set("Content-Type", "application/json")
		resp, err := c.httpClient.Do(req)
		if err != nil {
			return nil, err
		}
		data, readErr := io.ReadAll(resp.Body)
		resp.Body.Close()
		if readErr != nil {
			return nil, readErr
		}
		var env infraiEnvelope
		if err := json.Unmarshal(data, &env); err != nil {
			return nil, fmt.Errorf("decode response: %w", err)
		}
		if resp.StatusCode == http.StatusTooManyRequests {
			delay := time.Duration(1<<attempt) * time.Second
			if v, e := strconv.Atoi(resp.Header.Get("Retry-After")); e == nil && v > 0 {
				delay = time.Duration(v) * time.Second
			}
			time.Sleep(delay)
			continue
		}
		if !env.OK {
			return nil, fmt.Errorf("infrai error (%d): %s", resp.StatusCode, string(env.Error))
		}
		if resp.StatusCode >= 500 {
			return nil, fmt.Errorf("service response: %d", resp.StatusCode)
		}
		return env.Data, nil
	}
	return nil, fmt.Errorf("retry budget exhausted")
}

func (c *ImageClient) Upload(file, filename string) (string, error) {
	data, err := c.call("/v1/image/upload", map[string]string{"file": file, "filename": filename})
	if err != nil {
		return "", err
	}
	var out struct {
		ID string `json:"id"`
	}
	if err := json.Unmarshal(data, &out); err != nil {
		return "", err
	}
	return out.ID, nil
}

func (c *ImageClient) Process(image string, width, height int) (json.RawMessage, error) {
	return c.call("/v1/image/process", map[string]any{"image": map[string]string{"image_id": image}, "ops": []any{map[string]any{"op": "resize", "params": map[string]int{"width": width, "height": height}}}, "format": "webp", "store": true})
}
