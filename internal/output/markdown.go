package output

import (
	"fmt"
	"strings"

	"github.com/chuanjin/production-readiness/internal/engine"
)

// Markdown generates a human-readable report
func Markdown(summary engine.Summary, findings []engine.Finding) string {
	var b strings.Builder

	b.WriteString("# Production Readiness Report\n\n")
	b.WriteString(fmt.Sprintf("Overall Score: **%d / 100**\n\n", summary.Score))

	writeSection := func(title string, rules []engine.Finding) {
		if len(rules) == 0 {
			return
		}
		b.WriteString(fmt.Sprintf("## %s\n\n", title))
		for _, f := range rules {
			b.WriteString(fmt.Sprintf("### %s\n", f.Rule.Title))
			b.WriteString(f.Rule.Description + "\n\n")
			for _, w := range f.Rule.Why {
				b.WriteString("- " + w + "\n")
			}
			b.WriteString("\n")
		}
	}

	// Group findings by severity
	var high, medium, low []engine.Finding
	for _, f := range findings {
		if !f.Triggered {
			continue
		}
		switch f.Rule.Severity {
		case "high":
			high = append(high, f)
		case "medium":
			medium = append(medium, f)
		case "low":
			low = append(low, f)
		}
	}

	writeSection("🔴 High Risk", high)
	writeSection("🟠 Medium Risk", medium)
	writeSection("🟡 Low Risk", low)

	return b.String()
}
