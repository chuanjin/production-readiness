package rules

type Severity string

const (
	High   Severity = "high"
	Medium Severity = "medium"
	Low    Severity = "low"
)

type Rule struct {
	ID          string   `yaml:"id"`
	Severity    Severity `yaml:"severity"`
	Category    string   `yaml:"category"`
	Title       string   `yaml:"title"`
	Description string   `yaml:"description"`
	Why         []string `yaml:"why_it_matters"`
	Confidence  string   `yaml:"confidence"`
	Detect      Detect   `yaml:"detect"`
}

type Detect struct {
	AllOf  []map[string]interface{} `yaml:"all_of"`
	AnyOf  []map[string]interface{} `yaml:"any_of"`
	NoneOf []map[string]interface{} `yaml:"none_of"`
}
