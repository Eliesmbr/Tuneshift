package migration

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

type ProgressEvent struct {
	Type    string `json:"type"`
	Message string `json:"message"`
	Current int    `json:"current,omitempty"`
	Total   int    `json:"total,omitempty"`
}

type ProgressReporter struct {
	mu       sync.Mutex
	ch       chan ProgressEvent
	closed   bool
}

func NewProgressReporter() *ProgressReporter {
	return &ProgressReporter{
		ch: make(chan ProgressEvent, 100),
	}
}

func (pr *ProgressReporter) Send(event ProgressEvent) {
	pr.mu.Lock()
	defer pr.mu.Unlock()
	if !pr.closed {
		select {
		case pr.ch <- event:
		default:
			// Drop event if buffer is full
		}
	}
}

func (pr *ProgressReporter) Close() {
	pr.mu.Lock()
	defer pr.mu.Unlock()
	if !pr.closed {
		pr.closed = true
		close(pr.ch)
	}
}

func (pr *ProgressReporter) ServeSSE(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming not supported", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Accel-Buffering", "no")
	flusher.Flush()

	ctx := r.Context()
	for {
		select {
		case <-ctx.Done():
			return
		case event, ok := <-pr.ch:
			if !ok {
				return
			}
			data, _ := json.Marshal(event)
			fmt.Fprintf(w, "data: %s\n\n", data)
			flusher.Flush()
		}
	}
}
