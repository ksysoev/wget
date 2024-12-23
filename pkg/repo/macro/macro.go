package macro

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/ksysoev/wsget/pkg/core"
	"github.com/ksysoev/wsget/pkg/core/command"
	"gopkg.in/yaml.v3"
)

type config struct {
	Version string              `yaml:"version"`
	Source  string              `yaml:"source"`
	Macro   map[string][]string `yaml:"macro"`
	Domains []string            `yaml:"domains"`
}

type Repo struct {
	macro   map[string]*command.Templates
	domains []string
}

// New creates a new Repo instance with the specified domains.
// The domains parameter is a slice of strings representing the allowed domains for the macro.
// Returns a pointer to the newly created Repo instance.
func New(domains []string) *Repo {
	return &Repo{
		macro:   make(map[string]*command.Templates),
		domains: domains,
	}
}

// AddCommands adds a new macro with the given name and commands to the Repo instance.
// If a macro with the same name already exists, it returns an error.
// If the rawCommands slice is empty, it returns an error.
// If the rawCommands slice has only one command, it adds the command directly to the macro.
// Otherwise, it creates a new Sequence with the commands and adds it to the macro.
func (m *Repo) AddCommands(name string, rawCommands []string) error {
	if _, ok := m.macro[name]; ok {
		return fmt.Errorf("duplicate macro: %s", name)
	}

	if len(rawCommands) == 0 {
		return fmt.Errorf("empty macro: %s", name)
	}

	macro, err := command.NewMacro(rawCommands)

	if err != nil {
		return err
	}

	m.macro[name] = macro

	return nil
}

// merge merges the given macro into the current macro.
// If a macro with the same name already exists, an error is returned.
func (m *Repo) merge(macro *Repo) error {
	for name, cmd := range macro.macro {
		if _, ok := m.macro[name]; ok {
			return fmt.Errorf("duplicate macro: %s", name)
		}

		m.macro[name] = cmd
	}

	return nil
}

// Get returns the Executer associated with the given name, or an error if the name is not found.
func (m *Repo) Get(name, argString string) (core.Executer, error) {
	if cmd, ok := m.macro[name]; ok {
		args := strings.Fields(argString)
		return cmd.GetExecuter(args)
	}

	return nil, fmt.Errorf("unknown command: %s", name)
}

// GetNames returns a list of all macro names stored in the Repo instance.
// It does not take any parameters.
// It returns a slice of strings containing the names of the macros.
func (m *Repo) GetNames() []string {
	names := make([]string, 0, len(m.macro))

	for name := range m.macro {
		names = append(names, name)
	}

	return names
}

// LoadFromFile loads a macro configuration from a file at the given path.
// It returns a Repo instance and an error if the file cannot be read or parsed.
func LoadFromFile(path string) (r *Repo, err error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("fail to open macro file %s: %w", path, err)
	}

	defer func() {
		if e := file.Close(); err == nil && e != nil {
			err = fmt.Errorf("fail to close macro file %s: %w", path, e)
		}
	}()

	r, _, err = parseConfig(file)

	return r, err
}

// LoadMacroForDomain loads and merges macros for a specific domain from YAML files in a given directory.
// It takes macroDir, a string specifying the directory path, and domain, a string specifying the target domain.
// It returns a pointer to a Repo containing merged macros for the domain, or an error in case of failure.
// Errors may occur if the directory cannot be read, files cannot be parsed, or macros fail to merge.
// Ignores non-YAML files, directories, and files without a matching domain.
func LoadMacroForDomain(macroDir, domain string) (*Repo, error) {
	files, err := os.ReadDir(macroDir)
	if err != nil {
		log.Fatal(err)
	}

	var macro *Repo

	for _, file := range files {
		if file.IsDir() || (!strings.HasSuffix(file.Name(), ".yaml") && !strings.HasSuffix(file.Name(), ".yml")) {
			continue
		}

		fileMacro, err := LoadFromFile(macroDir + "/" + file.Name())

		if err != nil {
			return nil, err
		}

		hasDomain := false

		for _, fileDomain := range fileMacro.domains {
			if strings.HasSuffix(domain, fileDomain) {
				hasDomain = true
				break
			}
		}

		if !hasDomain {
			continue
		}

		if macro == nil {
			macro = fileMacro
		} else {
			err := macro.merge(fileMacro)

			if err != nil {
				return nil, fmt.Errorf("fail to loading macro from file %s, %w ", file.Name(), err)
			}
		}
	}

	return macro, nil
}

func parseConfig(src io.Reader) (*Repo, *config, error) {
	var cfg *config
	decoder := yaml.NewDecoder(src)
	if err := decoder.Decode(&cfg); err != nil {
		return nil, nil, err
	}

	if cfg.Version != "1" {
		return nil, nil, fmt.Errorf("unsupported macro version: %s", cfg.Version)
	}

	repo := New(cfg.Domains)

	for name, rawCommands := range cfg.Macro {
		if err := repo.AddCommands(name, rawCommands); err != nil {
			return nil, nil, fmt.Errorf("fail to add macro: %w", err)
		}
	}

	return repo, cfg, nil
}
