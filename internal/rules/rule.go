package rules

type Severity string

const (
	High     Severity = "high"
	Medium   Severity = "medium"
	Low      Severity = "low"
	Positive Severity = "positive"
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
	AnyOf  []Condition `yaml:"any_of"`
	AllOf  []Condition `yaml:"all_of"`
	NoneOf []Condition `yaml:"none_of"`
}

type Condition map[string]interface{}
