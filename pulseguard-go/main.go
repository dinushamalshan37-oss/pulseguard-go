package main

import (
	"context"
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
)

//go:embed web/*
var website embed.FS

type Check struct {
	At        time.Time `json:"at"`
	Status    int       `json:"status"`
	LatencyMS int64     `json:"latency_ms"`
	Online    bool      `json:"online"`
	Error     string    `json:"error,omitempty"`
}
type Target struct {
	ID     string  `json:"id"`
	Name   string  `json:"name"`
	URL    string  `json:"url"`
	Checks []Check `json:"checks"`
}
type Store struct {
	mu      sync.RWMutex
	Targets []Target `json:"targets"`
	path    string
}

func loadStore(path string) (*Store, error) {
	s := &Store{path: path, Targets: []Target{}}
	b, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return s, nil
	}
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(b, s); err != nil {
		return nil, fmt.Errorf("read state: %w", err)
	}
	if s.Targets == nil {
		s.Targets = []Target{}
	}
	return s, nil
}
func (s *Store) saveLocked() error {
	if err := os.MkdirAll(filepath.Dir(s.path), 0755); err != nil {
		return err
	}
	b, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	tmp := s.path + ".tmp"
	if err := os.WriteFile(tmp, b, 0600); err != nil {
		return err
	}
	return os.Rename(tmp, s.path)
}
func (s *Store) all() []Target {
	s.mu.RLock()
	defer s.mu.RUnlock()
	b, _ := json.Marshal(s.Targets)
	var copy []Target
	_ = json.Unmarshal(b, &copy)
	if copy == nil {
		return []Target{}
	}
	return copy
}
func validURL(raw string) bool {
	u, err := url.ParseRequestURI(raw)
	return err == nil && (u.Scheme == "https" || u.Scheme == "http") && u.Host != "" && u.User == nil
}
func (s *Store) add(name, raw string) (Target, error) {
	name = strings.TrimSpace(name)
	raw = strings.TrimSpace(raw)
	if len(name) < 2 || len(name) > 60 || len(raw) > 2000 || !validURL(raw) {
		return Target{}, errors.New("enter a 2–60 character name and an absolute http(s) URL")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(s.Targets) >= 25 {
		return Target{}, errors.New("maximum 25 websites per local instance")
	}
	t := Target{ID: strconv.FormatInt(time.Now().UnixNano(), 36), Name: name, URL: raw, Checks: []Check{}}
	s.Targets = append(s.Targets, t)
	if err := s.saveLocked(); err != nil {
		s.Targets = s.Targets[:len(s.Targets)-1]
		return Target{}, err
	}
	return t, nil
}
func (s *Store) remove(id string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, t := range s.Targets {
		if t.ID == id {
			s.Targets = append(s.Targets[:i], s.Targets[i+1:]...)
			if err := s.saveLocked(); err != nil {
				log.Printf("save: %v", err)
			}
			return true
		}
	}
	return false
}
func (s *Store) record(id string, c Check) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i := range s.Targets {
		if s.Targets[i].ID == id {
			s.Targets[i].Checks = append(s.Targets[i].Checks, c)
			if len(s.Targets[i].Checks) > 48 {
				s.Targets[i].Checks = s.Targets[i].Checks[len(s.Targets[i].Checks)-48:]
			}
			if err := s.saveLocked(); err != nil {
				log.Printf("save: %v", err)
			}
			return
		}
	}
}
func checkSite(ctx context.Context, client *http.Client, raw string) Check {
	begin := time.Now()
	result := Check{At: begin}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, raw, nil)
	if err != nil {
		result.Error = err.Error()
		return result
	}
	req.Header.Set("User-Agent", "PulseGuard/1.0 (+local health monitor)")
	res, err := client.Do(req)
	result.LatencyMS = time.Since(begin).Milliseconds()
	if err != nil {
		result.Error = err.Error()
		return result
	}
	defer res.Body.Close()
	result.Status = res.StatusCode
	result.Online = res.StatusCode >= 200 && res.StatusCode < 400
	return result
}
func runChecks(ctx context.Context, s *Store, client *http.Client) {
	targets := s.all()
	limiter := make(chan struct{}, 5)
	var wg sync.WaitGroup
	for _, t := range targets {
		wg.Add(1)
		go func(t Target) {
			defer wg.Done()
			select {
			case limiter <- struct{}{}:
			case <-ctx.Done():
				return
			}
			defer func() { <-limiter }()
			c := checkSite(ctx, client, t.URL)
			s.record(t.ID, c)
			if !c.Online {
				log.Printf("ALERT: %s (%s): status=%d error=%s", t.Name, t.URL, c.Status, c.Error)
			}
		}(t)
	}
	wg.Wait()
}
func replyJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
func envInt(key string, defaultValue int) int {
	v, err := strconv.Atoi(os.Getenv(key))
	if err != nil || v < 5 || v > 3600 {
		return defaultValue
	}
	return v
}
func main() {
	dataFile := os.Getenv("DATA_FILE")
	if dataFile == "" {
		dataFile = "data/sites.json"
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	s, err := loadStore(dataFile)
	if err != nil {
		log.Fatal(err)
	}
	interval := time.Duration(envInt("CHECK_INTERVAL_SECONDS", 30)) * time.Second
	client := &http.Client{Timeout: 7 * time.Second, CheckRedirect: func(req *http.Request, via []*http.Request) error {
		if len(via) >= 3 {
			return errors.New("too many redirects")
		}
		return nil
	}}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		runChecks(ctx, s, client)
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				runChecks(ctx, s, client)
			}
		}
	}()
	mux := http.NewServeMux()
	webRoot, err := fs.Sub(website, "web")
	if err != nil {
		log.Fatal(err)
	}
	mux.Handle("GET /", http.FileServer(http.FS(webRoot)))
	mux.HandleFunc("GET /api/targets", func(w http.ResponseWriter, r *http.Request) { replyJSON(w, 200, s.all()) })
	mux.HandleFunc("POST /api/targets", func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, 4096)
		var input struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		}
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			replyJSON(w, 400, map[string]string{"error": "invalid JSON"})
			return
		}
		t, err := s.add(input.Name, input.URL)
		if err != nil {
			replyJSON(w, 400, map[string]string{"error": err.Error()})
			return
		}
		replyJSON(w, 201, t)
		// A fresh monitor cycle will pick up the new target on the next interval.
	})
	mux.HandleFunc("DELETE /api/targets/{id}", func(w http.ResponseWriter, r *http.Request) {
		if !s.remove(r.PathValue("id")) {
			replyJSON(w, 404, map[string]string{"error": "not found"})
			return
		}
		w.WriteHeader(204)
	})
	server := &http.Server{Addr: "127.0.0.1:" + port, Handler: mux, ReadHeaderTimeout: 5 * time.Second}
	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = server.Shutdown(shutdownCtx)
	}()
	log.Printf("PulseGuard: http://127.0.0.1:%s (every %s)", port, interval)
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}
