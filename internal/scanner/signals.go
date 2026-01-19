package scanner

// RepoSignals holds scanned information
type RepoSignals struct {
	Files           map[string]bool   // tracks file existence
	FileContent     map[string]string // scanned file content (code only)
	BoolSignals     map[string]bool
	StringSignals   map[string]string
	IntSignals      map[string]int
	DetectedRegions map[string]bool
}
