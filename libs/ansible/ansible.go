package ansible

// Task represents a single task in an Ansible playbook
type Task struct {
	Name string            `yaml:"name" json:"name"`
	File map[string]string `yaml:"file,omitempty" json:"file,omitempty"`
	Loop []string          `yaml:"loop,omitempty" json:"loop,omitempty"`
}

// AnsiblePlaybook represents the structure of an Ansible playbook
type Play struct {
	Name   string            `yaml:"name" json:"name"`
	Hosts  string            `yaml:"hosts" json:"hosts"`
	Become bool              `yaml:"become" json:"become"`
	Vars   map[string]string `yaml:"vars,omitempty" json:"vars,omitempty"`
	Tasks  []Task            `yaml:"tasks" json:"tasks"`
}

// Playbook represents a list of playbooks (since YAML starts with ---)
type Playbook []Play
