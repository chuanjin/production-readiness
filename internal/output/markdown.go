package output

import (
	"fmt"
	"sort"
	"strings"

	"github.com/chuanjin/production-readiness/internal/engine"
	"github.com/chuanjin/production-readiness/internal/scanner"
)

// Markdown generates a human-readable report
func Markdown(summary engine.Summary, findings []engine.Finding, signals *scanner.RepoSignals) string {
	var b strings.Builder

	b.WriteString("# Production Readiness Report\n\n")
	b.WriteString(fmt.Sprintf("**Overall Score: %d / 100**\n\n", summary.Score))
	b.WriteString(fmt.Sprintf("- ✅ Passed: %d rules\n", summary.Passed))
	b.WriteString(fmt.Sprintf("- ❌ Triggered: %d rules\n", summary.Triggered))
	b.WriteString(fmt.Sprintf("- 📊 Total: %d rules\n\n", summary.Total))

	writeSection := func(title string, emoji string, findings []engine.Finding) {
		if len(findings) == 0 {
			return
		}
		b.WriteString(fmt.Sprintf("## %s %s\n\n", emoji, title))
		for i := range findings {
			f := &findings[i]
			b.WriteString(fmt.Sprintf("### %s\n\n", f.Rule.Title))
			b.WriteString(f.Rule.Description + "\n\n")

			if len(f.Rule.Why) > 0 {
				b.WriteString("**Why it matters:**\n")
				for _, w := range f.Rule.Why {
					b.WriteString("- " + w + "\n")
				}
				b.WriteString("\n")
			}
		}
	}

	// Group findings by severity
	var high, medium, low []engine.Finding
	for i := range findings {
		f := &findings[i]
		if !f.Triggered {
			continue
		}
		switch f.Rule.Severity {
		case "high":
			high = append(high, *f)
		case "medium":
			medium = append(medium, *f)
		case "low":
			low = append(low, *f)
		}
	}

	writeSection("High Risk", "🔴", high)
	writeSection("Medium Risk", "🟠", medium)
	writeSection("Low Risk", "🟡", low)

	// Add signals status section
	b.WriteString("---\n\n")
	b.WriteString("## 📊 Detected Signals\n\n")
	b.WriteString("These signals were detected during the repository scan:\n\n")

	// Boolean signals
	if len(signals.BoolSignals) > 0 {
		b.WriteString("### Boolean Signals\n\n")
		b.WriteString("| Signal | Status |\n")
		b.WriteString("|--------|--------|\n")

		// Sort for consistent output
		var keys []string
		for k := range signals.BoolSignals {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for _, key := range keys {
			value := signals.BoolSignals[key]
			status := "❌"
			if value {
				status = "✅"
			}
			b.WriteString(fmt.Sprintf("| `%s` | %s |\n", key, status))
		}
		b.WriteString("\n")
	}

	// String signals
	if len(signals.StringSignals) > 0 {
		b.WriteString("### String Signals\n\n")
		b.WriteString("| Signal | Value |\n")
		b.WriteString("|--------|-------|\n")

		var keys []string
		for k := range signals.StringSignals {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for _, key := range keys {
			value := signals.StringSignals[key]
			b.WriteString(fmt.Sprintf("| `%s` | `%s` |\n", key, value))
		}
		b.WriteString("\n")
	}

	// Integer signals
	if len(signals.IntSignals) > 0 {
		b.WriteString("### Integer Signals\n\n")
		b.WriteString("| Signal | Value |\n")
		b.WriteString("|--------|-------|\n")

		var keys []string
		for k := range signals.IntSignals {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for _, key := range keys {
			value := signals.IntSignals[key]
			b.WriteString(fmt.Sprintf("| `%s` | %d |\n", key, value))
		}
		b.WriteString("\n")
	}

	// File statistics
	b.WriteString("### Repository Statistics\n\n")
	b.WriteString(fmt.Sprintf("- **Files scanned:** %d\n", len(signals.Files)))
	b.WriteString(fmt.Sprintf("- **Files with content:** %d\n\n", len(signals.FileContent)))

	return b.String()
}

// MarkdownSummary generates a short summary report (without signals)
func MarkdownSummary(summary engine.Summary, findings []engine.Finding) string {
	var b strings.Builder

	b.WriteString("# Production Readiness Summary\n\n")
	b.WriteString(fmt.Sprintf("**Score: %d / 100**\n\n", summary.Score))

	// Count by severity
	var highCount, mediumCount, lowCount int

	for i := range findings {
		f := &findings[i]
		if !f.Triggered {
			continue
		}
		switch f.Rule.Severity {
		case "high":
			highCount++
		case "medium":
			mediumCount++
		case "low":
			lowCount++
		}
	}

	b.WriteString("## Issues Found\n\n")
	if highCount > 0 {
		b.WriteString(fmt.Sprintf("- 🔴 **High:** %d issues\n", highCount))
	}
	if mediumCount > 0 {
		b.WriteString(fmt.Sprintf("- 🟠 **Medium:** %d issues\n", mediumCount))
	}
	if lowCount > 0 {
		b.WriteString(fmt.Sprintf("- 🟡 **Low:** %d issues\n", lowCount))
	}
	if highCount == 0 && mediumCount == 0 && lowCount == 0 {
		b.WriteString("✅ No issues found!\n")
	}
	b.WriteString("\n")

	return b.String()
}
