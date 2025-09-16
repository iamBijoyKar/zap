package tasks

type Task struct {
	Name       string   `yaml:"name"`
	Command    []string `yaml:"command"`
	Retries    int      `yaml:"retries"`
	Parallel   bool     `yaml:"parallel"`
	Depends_On []string `yaml:"depends_on"`
}
type Yml struct {
	Tasks []Task `yaml:"tasks"`
}
