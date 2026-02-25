package herald

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// Store persists Herald rules to a JSON file.
type Store struct {
	mu   sync.RWMutex
	path string
}

// NewStore creates a store backed by ~/.ghabricator/herald-rules.json.
func NewStore() *Store {
	home, _ := os.UserHomeDir()
	dir := filepath.Join(home, ".ghabricator")
	os.MkdirAll(dir, 0o755)
	return &Store{path: filepath.Join(dir, "herald-rules.json")}
}

// List returns all rules.
func (s *Store) List() ([]Rule, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.readAll()
}

// Get returns a single rule by ID.
func (s *Store) Get(id string) (*Rule, error) {
	rules, err := s.List()
	if err != nil {
		return nil, err
	}
	for i := range rules {
		if rules[i].ID == id {
			return &rules[i], nil
		}
	}
	return nil, nil
}

// Save creates or updates a rule. If r.ID is empty, a new ID is assigned.
func (s *Store) Save(r *Rule) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	rules, err := s.readAll()
	if err != nil {
		rules = nil // treat as empty
	}

	now := time.Now()
	if r.ID == "" {
		r.ID = randomHex(8)
		r.CreatedAt = now
		r.UpdatedAt = now
		rules = append(rules, *r)
	} else {
		r.UpdatedAt = now
		found := false
		for i := range rules {
			if rules[i].ID == r.ID {
				r.CreatedAt = rules[i].CreatedAt
				rules[i] = *r
				found = true
				break
			}
		}
		if !found {
			r.CreatedAt = now
			rules = append(rules, *r)
		}
	}
	return s.writeAll(rules)
}

// Delete removes a rule by ID.
func (s *Store) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	rules, err := s.readAll()
	if err != nil {
		return err
	}
	filtered := rules[:0]
	for _, r := range rules {
		if r.ID != id {
			filtered = append(filtered, r)
		}
	}
	return s.writeAll(filtered)
}

func (s *Store) readAll() ([]Rule, error) {
	data, err := os.ReadFile(s.path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var rules []Rule
	if err := json.Unmarshal(data, &rules); err != nil {
		return nil, err
	}
	return rules, nil
}

func (s *Store) writeAll(rules []Rule) error {
	data, err := json.MarshalIndent(rules, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, data, 0o644)
}

func randomHex(n int) string {
	b := make([]byte, n)
	rand.Read(b)
	return hex.EncodeToString(b)
}
