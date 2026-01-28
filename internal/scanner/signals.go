package scanner

import "sync"

// RepoSignals holds scanned information
type RepoSignals struct {
	mu sync.RWMutex

	Files           map[string]bool   // tracks file existence
	FileContent     map[string]string // scanned file content (code only)
	BoolSignals     map[string]bool
	StringSignals   map[string]string
	IntSignals      map[string]int
	DetectedRegions map[string]bool
}

func (s *RepoSignals) SetFile(path string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Files[path] = true
}

func (s *RepoSignals) SetContent(path, content string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.FileContent[path] = content
}

func (s *RepoSignals) SetBool(key string, val bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.BoolSignals[key] = val
}

func (s *RepoSignals) SetString(key, val string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.StringSignals[key] = val
}

func (s *RepoSignals) SetInt(key string, val int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.IntSignals[key] = val
}

func (s *RepoSignals) SetRegion(region string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.DetectedRegions[region] = true
}

func (s *RepoSignals) GetBool(key string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.BoolSignals[key]
}

func (s *RepoSignals) GetBoolSignal(key string) (val, ok bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	val, ok = s.BoolSignals[key]
	return val, ok
}

func (s *RepoSignals) GetString(key string) string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.StringSignals[key]
}

func (s *RepoSignals) GetStringSignal(key string) (val string, ok bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	val, ok = s.StringSignals[key]
	return val, ok
}

func (s *RepoSignals) GetInt(key string) int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.IntSignals[key]
}

func (s *RepoSignals) GetIntSignal(key string) (val int, ok bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	val, ok = s.IntSignals[key]
	return val, ok
}

func (s *RepoSignals) HasFile(path string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.Files[path]
}

func (s *RepoSignals) GetContent(path string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	content, ok := s.FileContent[path]
	return content, ok
}

func (s *RepoSignals) GetRegionCount() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.DetectedRegions)
}

func (s *RepoSignals) GetFiles() map[string]bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	// Return a copy to avoid race on map iteration
	res := make(map[string]bool, len(s.Files))
	for k, v := range s.Files {
		res[k] = v
	}
	return res
}

func (s *RepoSignals) GetFileContentMap() map[string]string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	res := make(map[string]string, len(s.FileContent))
	for k, v := range s.FileContent {
		res[k] = v
	}
	return res
}
