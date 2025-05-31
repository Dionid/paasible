package paasible

import (
	"fmt"
	"regexp"
	"strings"
)

type EntityStorage struct {
	Origin *ConfigFile `yaml:"-"`

	UI       *UIEntity       `yaml:"ui"`
	Auth     *AuthEntity     `yaml:"auth"`
	Paasible *PaasibleEntity `yaml:"paasible"`

	SSHKeys      map[string]SSHKeyEntity      `yaml:"ssh_keys"`
	Hosts        map[string]HostEntity        `yaml:"hosts"`
	Inventories  map[string]InventoryEntity   `yaml:"inventories"`
	Projects     map[string]ProjectEntity     `yaml:"projects"`
	Performances map[string]PerformanceEntity `yaml:"performances"`
}

func NameToId(name string) string {
	if name == "" {
		return ""
	}

	// Convert to lowercase and replace spaces with underscores
	id := strings.ToLower(name)
	id = strings.ReplaceAll(id, " ", "_")

	// Remove any non-alphanumeric characters except underscores
	id = regexp.MustCompile(`[^a-z0-9_]`).ReplaceAllString(id, "")

	return id
}

func ParseConfigFile(storage *EntityStorage, origin *ConfigFile) error {
	for _, sshKey := range origin.SSHKeys {
		sshKey.Origin = origin
		if sshKey.Id == "" {
			sshKey.Id = NameToId(sshKey.Name)

			if sshKey.Id == "" {
				return fmt.Errorf("SSH key must have a name or ID: %v", sshKey)
			}
		}

		_, exist := storage.SSHKeys[sshKey.Id]
		if exist {
			return fmt.Errorf("SSH key with ID '%s' already exists in file '%s'", sshKey.Id, origin.FilePath)
		}

		storage.SSHKeys[sshKey.Id] = sshKey
	}

	for _, host := range origin.Hosts {
		host.Origin = origin
		if host.Id == "" {
			host.Id = NameToId(host.Name)

			if host.Id == "" {
				return fmt.Errorf("SSH key must have a name or ID: %v", host)
			}
		}

		_, exist := storage.Hosts[host.Id]
		if exist {
			return fmt.Errorf("SSH key with ID '%s' already exists in file '%s'", host.Id, origin.FilePath)
		}

		storage.Hosts[host.Id] = host
	}

	for _, inventory := range origin.Inventories {
		inventory.Origin = origin
		if inventory.Id == "" {
			inventory.Id = NameToId(inventory.Name)

			if inventory.Id == "" {
				return fmt.Errorf("SSH key must have a name or ID: %v", inventory)
			}
		}

		_, exist := storage.Inventories[inventory.Id]
		if exist {
			return fmt.Errorf("SSH key with ID '%s' already exists in file '%s'", inventory.Id, origin.FilePath)
		}

		storage.Inventories[inventory.Id] = inventory
	}

	for _, project := range origin.Projects {
		project.Origin = origin
		if project.Id == "" {
			project.Id = NameToId(project.Name)

			if project.Id == "" {
				return fmt.Errorf("SSH key must have a name or ID: %v", project)
			}
		}

		_, exist := storage.Projects[project.Id]
		if exist {
			return fmt.Errorf("SSH key with ID '%s' already exists in file '%s'", project.Id, origin.FilePath)
		}

		storage.Projects[project.Id] = project
	}

	for _, performance := range origin.Performances {
		performance.Origin = origin
		if performance.Id == "" {
			performance.Id = NameToId(performance.Name)

			if performance.Id == "" {
				return fmt.Errorf("SSH key must have a name or ID: %v", performance)
			}
		}

		_, exist := storage.Performances[performance.Id]
		if exist {
			return fmt.Errorf("SSH key with ID '%s' already exists in file '%s'", performance.Id, origin.FilePath)
		}

		storage.Performances[performance.Id] = performance
	}

	for _, include := range origin.Includes {
		err := ParseConfigFile(storage, include)
		if err != nil {
			return fmt.Errorf("Error parsing included config '%s': %w", include.FilePath, err)
		}
	}

	return nil
}

func EntityStorageFromOrigin(origin *ConfigFile) (*EntityStorage, error) {
	storage := &EntityStorage{
		Origin: origin,

		UI:       origin.UI,
		Auth:     origin.Auth,
		Paasible: origin.Paasible,

		SSHKeys:      make(map[string]SSHKeyEntity),
		Hosts:        make(map[string]HostEntity),
		Inventories:  make(map[string]InventoryEntity),
		Projects:     make(map[string]ProjectEntity),
		Performances: make(map[string]PerformanceEntity),
	}

	storage.UI.Origin = origin
	storage.Auth.Origin = origin
	storage.Paasible.Origin = origin

	err := ParseConfigFile(storage, origin)
	if err != nil {
		return nil, fmt.Errorf("Error parsing config file '%s': %w", origin.FilePath, err)
	}

	return storage, nil
}

type UIEntity struct {
	Origin *ConfigFile `yaml:"-"`

	Id   string `yaml:"id"`
	Port int    `yaml:"port"`
}

type AuthEntity struct {
	Origin *ConfigFile `yaml:"-"`

	Id      string `yaml:"id"`
	User    string `yaml:"user"`
	Machine string `yaml:"machine"`
}

type PaasibleEntity struct {
	Origin *ConfigFile `yaml:"-"`

	Id         string `yaml:"id"`
	CliVersion string `yaml:"cli_version"`
	CliEnvPath string `yaml:"cli_env_path"`

	DataFolderPath string `yaml:"data_folder_path"`
}

type SSHKeyEntity struct {
	Origin *ConfigFile `yaml:"-"`

	Id          string `yaml:"id"`
	Name        string `yaml:"name"`
	Description string `yaml:"description"`

	PrivatePath string `yaml:"private_path"`
	PublicPath  string `yaml:"public_path"`

	Private    string `yaml:"private"`
	Public     string `yaml:"public"`
	Passphrase string `yaml:"passphrase"` // optional passphrase for the private key
}

type HostEntity struct {
	Origin *ConfigFile `yaml:"-"`

	Id          string   `yaml:"id"`
	Name        string   `yaml:"name"`
	Description string   `yaml:"description"`
	Address     string   `yaml:"address"`
	SSHKeys     []string `yaml:"ssh_keys"`
}

type InventoryEntity struct {
	Origin *ConfigFile `yaml:"-"`

	Id          string `yaml:"id"`
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	Path        string `yaml:"path"`
}

type ProjectEntity struct {
	Origin *ConfigFile `yaml:"-"`

	Id               string     `yaml:"id"`
	Name             string     `yaml:"name"`
	Description      string     `yaml:"description"`
	Version          string     `yaml:"version"`
	CliVersion       string     `yaml:"cli_version"`
	Repository       string     `yaml:"repository"`
	RepositoryBranch string     `yaml:"repository_branch"`
	RepositoryPath   string     `yaml:"repository_path"`
	LocalPath        string     `yaml:"local_path"`
	Playbooks        []Playbook `yaml:"playbooks"`
}

type PlaybookEntity struct {
	Origin *ConfigFile `yaml:"-"`

	Id              string               `yaml:"id"`
	ID              string               `yaml:"id"`
	Name            string               `yaml:"name"`
	Description     string               `yaml:"description"`
	Path            string               `yaml:"path"`
	VariablesSchema VariableSchemaEntity `yaml:"variables_schema"`
}

type VariableSchemaEntity struct {
	Origin *ConfigFile `yaml:"-"`

	Id          string            `yaml:"id"`
	Name        string            `yaml:"name"`
	Description string            `yaml:"description"`
	Schema      map[string]string `yaml:"schema"`
}

type PerformanceEntity struct {
	Origin *ConfigFile `yaml:"-"`

	Id          string           `yaml:"id"`
	ID          string           `yaml:"id"`
	Name        string           `yaml:"name"`
	Description string           `yaml:"description"`
	Playbooks   []string         `yaml:"playbooks"`
	Inventories []string         `yaml:"inventories"`
	Variables   []VariableEntity `yaml:"variables"`
	Targets     []Target         `yaml:"targets"`
}

type VariableEntity struct {
	Origin *ConfigFile `yaml:"-"`

	Id   string `yaml:"id"`
	Path string `yaml:"path"`
}

type TargetEntity struct {
	Origin *ConfigFile `yaml:"-"`

	Id     string `yaml:"id"`
	Host   string `yaml:"host"`
	SSHKey string `yaml:"ssh_key"`
	User   string `yaml:"user"`
	Port   int    `yaml:"port"`
}
