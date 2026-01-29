package scanner

import (
	"testing"
)

func TestRepoSignals(t *testing.T) {
	s := &RepoSignals{
		Files:           make(map[string]bool),
		FileContent:     make(map[string]string),
		BoolSignals:     make(map[string]bool),
		StringSignals:   make(map[string]string),
		IntSignals:      make(map[string]int),
		DetectedRegions: make(map[string]bool),
	}

	t.Run("Files", func(t *testing.T) { testSignalsFiles(t, s) })
	t.Run("FileContent", func(t *testing.T) { testSignalsFileContent(t, s) })
	t.Run("BoolSignals", func(t *testing.T) { testSignalsBool(t, s) })
	t.Run("StringSignals", func(t *testing.T) { testSignalsString(t, s) })
	t.Run("IntSignals", func(t *testing.T) { testSignalsInt(t, s) })
	t.Run("Regions", func(t *testing.T) { testSignalsRegions(t, s) })
}

func testSignalsFiles(t *testing.T, s *RepoSignals) {
	s.SetFile("main.go")
	if !s.HasFile("main.go") {
		t.Errorf("expected HasFile(main.go) to be true")
	}
	if s.HasFile("nonexistent.go") {
		t.Errorf("expected HasFile(nonexistent.go) to be false")
	}

	files := s.GetFiles()
	if !files["main.go"] {
		t.Errorf("expected GetFiles to contain main.go")
	}
	// Verify copy
	files["hacked.go"] = true
	if s.HasFile("hacked.go") {
		t.Errorf("GetFiles should return a copy, not a reference")
	}
}

func testSignalsFileContent(t *testing.T, s *RepoSignals) {
	s.SetContent("main.go", "package main")
	content, ok := s.GetContent("main.go")
	if !ok || content != "package main" {
		t.Errorf("expected GetContent(main.go) to return 'package main'")
	}
	_, ok = s.GetContent("nonexistent.go")
	if ok {
		t.Errorf("expected GetContent(nonexistent.go) to return ok=false")
	}

	contentMap := s.GetFileContentMap()
	if contentMap["main.go"] != "package main" {
		t.Errorf("expected GetFileContentMap to contain main.go content")
	}
	// Verify copy
	contentMap["hacked.go"] = "hacked"
	_, ok = s.GetContent("hacked.go")
	if ok {
		t.Errorf("GetFileContentMap should return a copy, not a reference")
	}
}

func testSignalsBool(t *testing.T, s *RepoSignals) {
	s.SetBool("test_bool", true)
	if !s.GetBool("test_bool") {
		t.Errorf("expected GetBool(test_bool) to be true")
	}
	valBool, ok := s.GetBoolSignal("test_bool")
	if !ok || !valBool {
		t.Errorf("expected GetBoolSignal(test_bool) to be true, true")
	}
	_, ok = s.GetBoolSignal("nonexistent")
	if ok {
		t.Errorf("expected GetBoolSignal(nonexistent) to return ok=false")
	}
}

func testSignalsString(t *testing.T, s *RepoSignals) {
	s.SetString("test_string", "hello")
	if s.GetString("test_string") != "hello" {
		t.Errorf("expected GetString(test_string) to be 'hello'")
	}
	valString, ok := s.GetStringSignal("test_string")
	if !ok || valString != "hello" {
		t.Errorf("expected GetStringSignal(test_string) to be 'hello', true")
	}
	_, ok = s.GetStringSignal("nonexistent")
	if ok {
		t.Errorf("expected GetStringSignal(nonexistent) to return ok=false")
	}
}

func testSignalsInt(t *testing.T, s *RepoSignals) {
	s.SetInt("test_int", 42)
	if s.GetInt("test_int") != 42 {
		t.Errorf("expected GetInt(test_int) to be 42")
	}
	valInt, ok := s.GetIntSignal("test_int")
	if !ok || valInt != 42 {
		t.Errorf("expected GetIntSignal(test_int) to be 42, true")
	}
	_, ok = s.GetIntSignal("nonexistent")
	if ok {
		t.Errorf("expected GetIntSignal(nonexistent) to return ok=false")
	}
}

func testSignalsRegions(t *testing.T, s *RepoSignals) {
	s.SetRegion("us-east-1")
	s.SetRegion("us-west-2")
	if s.GetRegionCount() != 2 {
		t.Errorf("expected GetRegionCount to be 2, got %d", s.GetRegionCount())
	}
}

func TestRepoSignalsConcurrency(t *testing.T) {
	s := &RepoSignals{
		Files:           make(map[string]bool),
		FileContent:     make(map[string]string),
		BoolSignals:     make(map[string]bool),
		StringSignals:   make(map[string]string),
		IntSignals:      make(map[string]int),
		DetectedRegions: make(map[string]bool),
	}

	done := make(chan bool)
	go func() {
		for i := 0; i < 1000; i++ {
			s.SetBool("concurrency", true)
			_ = s.GetBool("concurrency")
		}
		done <- true
	}()

	go func() {
		for i := 0; i < 1000; i++ {
			s.SetInt("concurrency", i)
			_ = s.GetInt("concurrency")
		}
		done <- true
	}()

	<-done
	<-done
}
