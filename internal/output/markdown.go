package output

import (
	"fmt"
	"strings"

	"github.com/chuanjin/production-readiness/internal/engine"
)

// Markdown generates a human-readable report from Summary
func Markdown(summary engine.Summary) string {
	var b strings.Builder

	b.WriteString("# Production Readiness Report\n\n")
	b.WriteString(fmt.Sprintf("Overall Score: **%d / 100**\n\n", summary.Score))

	for _, f := range summary.Findings {
		// Skip rules that are unsupported
		if !f.Supported {
			b.WriteString(fmt.Sprintf("### ⚪ Skipped: %s\n", f.Rule.Title))
			b.WriteString(f.Rule.Description + "\n\n")
			continue
		}

		if f.Triggered {
			var emoji string
			switch f.Rule.Severity {
			case "high":
				emoji = "🔴 High Risk"
			case "medium":
				emoji = "🟠 Medium Risk"
			case "low":
				emoji = "🟡 Low Risk"
			case "positive":
				emoji = "🟢 Good Signal"
			default:
				emoji = "⚪ Unknown"
			}

			b.WriteString(fmt.Sprintf("## %s — %s\n", emoji, f.Rule.Title))
			b.WriteString(f.Rule.Description + "\n\n")
			for _, w := range f.Rule.Why {
				b.WriteString("- " + w + "\n")
			}
			b.WriteString("\n")
		} else {
			// Supported but not triggered
			b.WriteString(fmt.Sprintf("### 🟢 Passed: %s\n", f.Rule.Title))
			b.WriteString(f.Rule.Description + "\n\n")
		}
	}

	return b.String()
}
