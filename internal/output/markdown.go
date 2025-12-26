package output

import (
	"fmt"
	"strings"

	"github.com/chuanjin/production-readiness/internal/engine"
)

func Markdown(summary engine.Summary, findings []engine.Finding) string {
	var b strings.Builder

	b.WriteString("# Production Readiness Report\n\n")
	b.WriteString(fmt.Sprintf("Overall Score: **%d / 100**\n\n", summary.Score))

	for _, f := range findings {
		if !f.Triggered {
			continue
		}
		b.WriteString(fmt.Sprintf("## %s\n", f.Rule.Title))
		b.WriteString(f.Rule.Description + "\n\n")
		for _, w := range f.Rule.Why {
			b.WriteString("- " + w + "\n")
		}
		b.WriteString("\n")
	}

	return b.String()
}
