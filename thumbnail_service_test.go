package main

import "testing"

func TestChooseSize(t *testing.T) {
	tests := []struct {
		name               string
		w, h, wantW, wantH int
	}{
		{"default", 0, 0, 320, 180}, {"cap wide image", 1200, 600, 640, 320}, {"keep small", 400, 300, 400, 300},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w, h := chooseSize(tt.w, tt.h)
			if w != tt.wantW || h != tt.wantH {
				t.Fatalf("got %dx%d, want %dx%d", w, h, tt.wantW, tt.wantH)
			}
		})
	}
}
