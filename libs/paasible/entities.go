package paasible

import (
	"fmt"
	"regexp"
	"strings"
)

type EntityStorage struct {
	Origin *ConfigFile

	UI       *UIEntity
	Auth     *AuthEntity
	Paasible *PaasibleEntity

	SSHKeys      map[string]SSHKeyEntity
	Hosts        map[string]HostEntity
	Inventories  map[string]InventoryEntity
	Projects     map[string]ProjectEntity
	Playbooks    map[string]PlaybookEntity
	Performances map[string]PerformanceEntity
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
	if origin.UI != nil {
		storage.UI = origin.UI
		storage.UI.Origin = origin
	}

	if origin.Auth != nil {
		storage.Auth = origin.Auth
		storage.Auth.Origin = origin
	}

	for _, sshKey := range origin.SshKeys {
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
				return fmt.Errorf("Host must have a name or ID: %v", host)
			}
		}

		_, exist := storage.Hosts[host.Id]
		if exist {
			return fmt.Errorf("Host with ID '%s' already exists in file '%s'", host.Id, origin.FilePath)
		}

		storage.Hosts[host.Id] = host
	}

	for _, inventory := range origin.Inventories {
		inventory.Origin = origin
		if inventory.Id == "" {
			inventory.Id = NameToId(inventory.Name)

			if inventory.Id == "" {
				return fmt.Errorf("Inventory must have a name or ID: %v", inventory)
			}
		}

		_, exist := storage.Inventories[inventory.Id]
		if exist {
			return fmt.Errorf("Inventory with ID '%s' already exists in file '%s'", inventory.Id, origin.FilePath)
		}

		storage.Inventories[inventory.Id] = inventory
	}

	for _, project := range origin.Projects {
		project.Origin = origin
		if project.Id == "" {
			project.Id = NameToId(project.Name)

			if project.Id == "" {
				return fmt.Errorf("Project must have a name or ID: %v", project)
			}
		}

		_, exist := storage.Projects[project.Id]
		if exist {
			return fmt.Errorf("Project with ID '%s' already exists in file '%s'", project.Id, origin.FilePath)
		}

		storage.Projects[project.Id] = project

		for _, playbook := range project.Playbooks {
			playbook.Origin = origin
			if playbook.Id == "" {
				playbook.Id = NameToId(playbook.Name)

				if playbook.Id == "" {
					return fmt.Errorf("Playbook must have a name or ID: %v", playbook)
				}
			}

			playbook.Id = project.Id + "." + playbook.Id
			playbook.ProjectId = project.Id

			_, exist := storage.Playbooks[playbook.Id]
			if exist {
				return fmt.Errorf("Playbook with ID '%s' already exists in file '%s'", playbook.Id, origin.FilePath)
			}

			storage.Playbooks[playbook.Id] = playbook
		}
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
		Playbooks:    make(map[string]PlaybookEntity),
		Performances: make(map[string]PerformanceEntity),
	}

	storage.Paasible.Origin = origin

	err := ParseConfigFile(storage, origin)
	if err != nil {
		return nil, fmt.Errorf("Error parsing config file '%s': %w", origin.FilePath, err)
	}

	return storage, nil
}

type UIEntity struct {
	Origin *ConfigFile `mapstructure:"-"`

	Id   string `mapstructure:"id"`
	Port int    `mapstructure:"port"`
}

type AuthEntity struct {
	Origin *ConfigFile `mapstructure:"-"`

	Id      string `mapstructure:"id"`
	User    string `mapstructure:"user"`
	Machine string `mapstructure:"machine"`
}

type PaasibleEntity struct {
	Origin *ConfigFile `mapstructure:"-"`

	CliVersion string `mapstructure:"cli_version"`
	CliEnvPath string `mapstructure:"cli_env_path"`

	DataFolderPath string `mapstructure:"data_folder_path"`
}

type SSHKeyEntity struct {
	Origin *ConfigFile `mapstructure:"-"`

	Id          string `mapstructure:"id"`
	Name        string `mapstructure:"name"`
	Description string `mapstructure:"description"`

	PrivatePath string `mapstructure:"private_path"`
	PublicPath  string `mapstructure:"public_path"`
}

type HostEntity struct {
	Origin *ConfigFile `mapstructure:"-"`

	Id          string   `mapstructure:"id"`
	Name        string   `mapstructure:"name"`
	Description string   `mapstructure:"description"`
	Address     string   `mapstructure:"address"`
	SSHKeys     []string `mapstructure:"ssh_keys"`
}

type InventoryEntity struct {
	Origin *ConfigFile `mapstructure:"-"`

	Id          string `mapstructure:"id"`
	Name        string `mapstructure:"name"`
	Description string `mapstructure:"description"`
	Path        string `mapstructure:"path"`
}

type ProjectEntity struct {
	Origin *ConfigFile `mapstructure:"-"`

	Id          string           `mapstructure:"id"`
	Name        string           `mapstructure:"name"`
	Description string           `mapstructure:"description"`
	Version     string           `mapstructure:"version"`
	LocalPath   string           `mapstructure:"local_path"`
	Playbooks   []PlaybookEntity `mapstructure:"playbooks"`
}

type PlaybookEntity struct {
	Origin *ConfigFile `mapstructure:"-"`

	Id          string `mapstructure:"id"`
	ID          string `mapstructure:"id"`
	Name        string `mapstructure:"name"`
	Description string `mapstructure:"description"`
	Path        string `mapstructure:"path"`
	ProjectId   string `mapstructure:"project_id"`
}

type VariableSchemaEntity struct {
	Origin *ConfigFile `mapstructure:"-"`

	Id          string            `mapstructure:"id"`
	Name        string            `mapstructure:"name"`
	Description string            `mapstructure:"description"`
	Schema      map[string]string `mapstructure:"schema"`
}

type PerformanceEntityPlaybook struct {
	Project  string `mapstructure:"project"`
	Playbook string `mapstructure:"playbook"`
}

type PerformanceEntity struct {
	Origin *ConfigFile `mapstructure:"-"`

	Id          string                      `mapstructure:"id"`
	ID          string                      `mapstructure:"id"`
	Name        string                      `mapstructure:"name"`
	Description string                      `mapstructure:"description"`
	Playbooks   []PerformanceEntityPlaybook `mapstructure:"playbooks"`
	Inventories []string                    `mapstructure:"inventories"`
	Targets     []TargetEntity              `mapstructure:"targets"`
}

type VariableEntity struct {
	Origin *ConfigFile `mapstructure:"-"`

	Id   string `mapstructure:"id"`
	Path string `mapstructure:"path"`
}

type TargetEntity struct {
	Origin *ConfigFile `mapstructure:"-"`

	Id     string `mapstructure:"id"`
	Host   string `mapstructure:"host"`
	SSHKey string `mapstructure:"ssh_key"`
	User   string `mapstructure:"user"`
	Group  string `mapstructure:"group"`
	Port   int    `mapstructure:"port"`
}
