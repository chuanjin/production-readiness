package output

import (
	"fmt"
	"strings"

	"github.com/chuanjin/production-readiness/internal/engine"
)

func Markdown(summary engine.Summary) string {
	var b strings.Builder

	b.WriteString("# Production Readiness Scan Result\n\n")
	b.WriteString(fmt.Sprintf("**Overall Score:** %d/100\n\n", summary.Score))

	// Issues Found
	if summary.High+summary.Medium+summary.Low > 0 {
		b.WriteString("## Issues Found\n\n")
		for _, f := range summary.Findings {
			if f.Triggered && f.Supported {
				sev := strings.ToUpper(string(f.Rule.Severity))
				b.WriteString(fmt.Sprintf("### [%s] %s\n", sev, f.Rule.Title))
				b.WriteString(f.Rule.Description + "\n\n")
			}
		}
	} else {
		b.WriteString("No readiness risks detected.\n\n")
	}

	// Positive Signals
	if summary.Positive > 0 {
		b.WriteString("## Positive Signals\n\n")
		for _, f := range summary.Findings {
			if f.Triggered && f.Rule.Severity == "positive" {
				b.WriteString(fmt.Sprintf("### %s\n", f.Rule.Title))
				b.WriteString(f.Rule.Description + "\n\n")
			}
		}
	}

	// Supported but NOT triggered
	if len(summary.NotHit) > 0 {
		b.WriteString("## Checked but Not Triggered (Expected default state)\n\n")
		for _, f := range summary.NotHit {
			b.WriteString(fmt.Sprintf("- %s\n", f.Rule.Title))
		}
		b.WriteString("\n")
	}

	// Unsupported Rules
	if len(summary.Skipped) > 0 {
		b.WriteString("## Skipped (Rule logic not implemented yet)\n\n")
		for _, f := range summary.Skipped {
			sev := strings.ToUpper(string(f.Rule.Severity))
			b.WriteString(fmt.Sprintf("- %s (%s)\n", f.Rule.Title, sev))
		}
		b.WriteString("\n")
	}

	return b.String()
}
