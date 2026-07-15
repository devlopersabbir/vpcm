package plugin

type Command struct {
	Name        string
	Description string
	Execute     func(args []string) error
}

type Collector interface {
	Name() string
	Collect() (any, error)
}

type Hook interface {
	Event() string
	Trigger(payload any) error
}

type Plugin interface {
	Name() string
	Version() string
	Register() error
	Commands() []Command
	Collectors() []Collector
	Hooks() []Hook
}
