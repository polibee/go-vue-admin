package plugin

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

var (
	ErrInvalidPlugin        = errors.New("invalid plugin")
	ErrDuplicatePlugin      = errors.New("plugin is already registered")
	ErrPluginNotFound       = errors.New("plugin is not registered")
	ErrBuiltinOnly          = errors.New("only builtin plugins are supported")
	ErrDependencyNotFound   = errors.New("plugin dependency is not registered")
	ErrDependencyNotEnabled = errors.New("plugin dependency is not enabled")
	ErrDependencyActive     = errors.New("plugin dependency is active")
)

type PluginRuntime string

const (
	RuntimeBuiltin  PluginRuntime = "builtin"
	RuntimeExternal PluginRuntime = "external"
)

type PluginState string

const (
	PluginDisabled PluginState = "disabled"
	PluginEnabled  PluginState = "enabled"
)

type UICompatibility string

const UICompatibilityShadcnVue UICompatibility = "shadcn-vue"

type PluginDependency struct {
	ID      string `json:"id"`
	Version string `json:"version,omitempty"`
}

type ConfigDeclaration struct {
	Route        string   `json:"route"`
	Label        string   `json:"label"`
	Permission   string   `json:"permission"`
	Schema       string   `json:"schema,omitempty"`
	SecretFields []string `json:"secretFields,omitempty"`
}

type PluginMenu struct {
	ID         string `json:"id"`
	Label      string `json:"label"`
	Route      string `json:"route"`
	Permission string `json:"permission,omitempty"`
}

type PluginManifest struct {
	ID              string             `json:"id"`
	Name            string             `json:"name"`
	Version         string             `json:"version"`
	Runtime         PluginRuntime      `json:"runtime"`
	UICompatibility UICompatibility    `json:"uiCompatibility"`
	Permissions     []string           `json:"permissions,omitempty"`
	Menus           []PluginMenu       `json:"menus,omitempty"`
	Dependencies    []PluginDependency `json:"dependencies,omitempty"`
	Config          *ConfigDeclaration `json:"config,omitempty"`
}

type BuiltinPlugin interface {
	Manifest() PluginManifest
	Enable() error
	Disable() error
}

type PluginInfo struct {
	Manifest PluginManifest
	State    PluginState
}

var pluginIDPattern = regexp.MustCompile(`^[a-z][a-z0-9_-]*$`)

func validateManifest(manifest PluginManifest) error {
	manifest.ID = strings.TrimSpace(manifest.ID)
	manifest.Name = strings.TrimSpace(manifest.Name)
	manifest.Version = strings.TrimSpace(manifest.Version)
	if !pluginIDPattern.MatchString(manifest.ID) || manifest.Name == "" || manifest.Version == "" {
		return fmt.Errorf("%w: id, name, and version are required", ErrInvalidPlugin)
	}
	if manifest.Runtime != RuntimeBuiltin {
		return fmt.Errorf("%w: %s", ErrBuiltinOnly, manifest.Runtime)
	}
	if manifest.UICompatibility != UICompatibilityShadcnVue {
		return fmt.Errorf("%w: unsupported uiCompatibility %q", ErrInvalidPlugin, manifest.UICompatibility)
	}
	seenDependencies := make(map[string]struct{}, len(manifest.Dependencies))
	for _, dependency := range manifest.Dependencies {
		id := strings.TrimSpace(dependency.ID)
		if !pluginIDPattern.MatchString(id) || id == manifest.ID {
			return fmt.Errorf("%w: invalid dependency %q", ErrInvalidPlugin, dependency.ID)
		}
		if _, exists := seenDependencies[id]; exists {
			return fmt.Errorf("%w: duplicate dependency %q", ErrInvalidPlugin, id)
		}
		seenDependencies[id] = struct{}{}
	}
	seenMenus := make(map[string]struct{}, len(manifest.Menus))
	for _, menu := range manifest.Menus {
		id := strings.TrimSpace(menu.ID)
		if id == "" || strings.TrimSpace(menu.Label) == "" || strings.TrimSpace(menu.Route) == "" {
			return fmt.Errorf("%w: menu id, label, and route are required", ErrInvalidPlugin)
		}
		if _, exists := seenMenus[id]; exists {
			return fmt.Errorf("%w: duplicate menu %q", ErrInvalidPlugin, id)
		}
		seenMenus[id] = struct{}{}
	}
	return nil
}

func cloneManifest(manifest PluginManifest) PluginManifest {
	manifest.Permissions = append([]string(nil), manifest.Permissions...)
	manifest.Menus = append([]PluginMenu(nil), manifest.Menus...)
	manifest.Dependencies = append([]PluginDependency(nil), manifest.Dependencies...)
	if manifest.Config != nil {
		config := *manifest.Config
		config.SecretFields = append([]string(nil), config.SecretFields...)
		manifest.Config = &config
	}
	return manifest
}
