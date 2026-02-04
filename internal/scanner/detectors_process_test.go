package scanner

import (
	"testing"
)

func TestDetectManualSteps(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		relPath  string
		expected bool
	}{
		{
			name: "Manual steps found in README",
			content: `
To deploy this application:
1. Manually copy the config file
2. Run command "go run main.go"
3. Execute script "init.sh"
`,
			relPath:  "README.md",
			expected: true,
		},
		{
			name: "Not enough indicators",
			content: `
Simple deployment guide.
Just run command "make deploy"
`,
			relPath:  "README.md",
			expected: false,
		},
		{
			name: "Ignored file type",
			content: `
// Manually copy
// Run command
// Execute script
`,
			relPath:  "main.go",
			expected: false,
		},
		{
			name: "Documentation with no manual steps",
			content: `
This project uses automated deployment via GitHub Actions.
`,
			relPath:  "docs/deploy.md",
			expected: false,
		},
		{
			name: "Manual steps in txt file",
			content: `
Deployment instructions:
1. manual step required
2. run command manually
3. execute script
`,
			relPath:  "INSTALL.txt",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			signals := &RepoSignals{
				BoolSignals: make(map[string]bool),
			}
			detectManualSteps(tt.content, tt.relPath, signals)

			if got := signals.GetBool("manual_steps_documented"); got != tt.expected {
				t.Errorf("detectManualSteps() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestDetectMigrationTool(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected bool
	}{
		{
			name:     "Golang-migrate detected",
			content:  `import "github.com/golang-migrate/migrate"`,
			expected: true,
		},
		{
			name:     "Alembic detected",
			content:  `from alembic import context`,
			expected: true,
		},
		{
			name:     "Flyway detected",
			content:  `Flyway flyway = Flyway.configure()`,
			expected: true,
		},
		{
			name:     "No migration tool",
			content:  `SELECT * FROM users;`,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			signals := &RepoSignals{
				BoolSignals: make(map[string]bool),
			}
			detectMigrationTool(tt.content, "migration.go", signals)

			if signals.GetBool("migration_tool_detected") != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, signals.GetBool("migration_tool_detected"))
			}
		})
	}
}

func TestDetectBackwardCompatibleMigration(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected bool
	}{
		{
			name:     "Backward compatible mentioned",
			content:  `# This migration is backward compatible`,
			expected: true,
		},
		{
			name:     "Zero-downtime deployment",
			content:  `Supports zero-downtime deployment`,
			expected: true,
		},
		{
			name:     "Nullable column",
			content:  `ALTER TABLE users ADD COLUMN email VARCHAR(255) NULL DEFAULT '';`,
			expected: true,
		},
		{
			name:     "No compatibility info",
			content:  `ALTER TABLE users DROP COLUMN email;`,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			signals := &RepoSignals{
				BoolSignals: make(map[string]bool),
			}
			detectBackwardCompatibleMigration(tt.content, "migration.sql", signals)

			if signals.GetBool("backward_compatible_migration_hint") != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, signals.GetBool("backward_compatible_migration_hint"))
			}
		})
	}
}

func TestDetectMigrationValidation(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected bool
	}{
		{
			name:     "Dry-run flag",
			content:  `migrate up --dry-run`,
			expected: true,
		},
		{
			name:     "Validation test",
			content:  `def test_migration_validation():`,
			expected: true,
		},
		{
			name:     "Rollback test",
			content:  `it('should rollback successfully', () => {})`,
			expected: true,
		},
		{
			name:     "No validation",
			content:  `migrate up`,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			signals := &RepoSignals{
				BoolSignals: make(map[string]bool),
			}
			detectMigrationValidation(tt.content, "test.sh", signals)

			if signals.GetBool("migration_validation_step") != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, signals.GetBool("migration_validation_step"))
			}
		})
	}
}

func TestDetectUnsafeMigration(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected bool
	}{
		{
			name:     "Drop table detected",
			content:  `DROP TABLE users;`,
			expected: true,
		},
		{
			name:     "Rename column detected",
			content:  `ALTER TABLE users RENAME COLUMN name TO full_name;`,
			expected: true,
		},
		{
			name:     "Alter column type detected",
			content:  `ALTER TABLE users ALTER COLUMN age TYPE bigint;`,
			expected: true,
		},
		{
			name:     "Safe migration",
			content:  `ALTER TABLE users ADD COLUMN email TEXT;`,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			signals := &RepoSignals{
				BoolSignals: make(map[string]bool),
			}
			detectUnsafeMigration(tt.content, "migration.sql", signals)

			if signals.GetBool("unsafe_migration_detected") != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, signals.GetBool("unsafe_migration_detected"))
			}
		})
	}
}

func TestDetectGracefulShutdown(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected bool
	}{
		{
			name:     "Go signal notify",
			content:  `signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)`,
			expected: true,
		},
		{
			name:     "Node.js sigterm listener",
			content:  `process.on('SIGTERM', () => { server.close() })`,
			expected: true,
		},
		{
			name:     "Spring graceful shutdown",
			content:  `spring.lifecycle.timeout-per-shutdown-phase=20s`,
			expected: true,
		},
		{
			name:     "No signal handling",
			content:  `func main() { fmt.Println("hello") }`,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			signals := &RepoSignals{
				BoolSignals: make(map[string]bool),
			}
			detectGracefulShutdown(tt.content, "main.go", signals)

			if signals.GetBool("graceful_shutdown_detected") != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, signals.GetBool("graceful_shutdown_detected"))
			}
		})
	}
}
