package generator

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"
)

var (
	resourceNamePattern = regexp.MustCompile(`^[a-z][a-z0-9]*(?:-[a-z0-9]+)*$`)
	fieldNamePattern    = regexp.MustCompile(`^[a-z][a-z0-9_]*$`)
)

var supportedFieldTypes = map[string]struct{}{
	"text": {}, "email": {}, "password": {}, "integer": {}, "boolean": {}, "select": {},
}

type Input struct {
	Name       string
	Label      string
	Route      string
	Permission string
	Icon       string
	Actions    []string
	Fields     []string
}

type Spec struct {
	Name          string
	GoName        string
	Label         string
	Route         string
	FrontendRoute string
	Permission    string
	Icon          string
	Actions       []string
	Fields        []FieldSpec
}

type FieldSpec struct {
	Name     string
	GoName   string
	Type     string
	Required bool
	Options  []FieldOption
}

type FieldOption struct {
	Value string
	Label string
}

func ParseField(value string) (FieldSpec, error) {
	parts := strings.Split(value, ":")
	if len(parts) < 2 || len(parts) > 4 || !fieldNamePattern.MatchString(parts[0]) {
		return FieldSpec{}, fmt.Errorf("invalid field %q: expected name:type[:required[:value=Label|value=Label]]", value)
	}
	if _, ok := supportedFieldTypes[parts[1]]; !ok {
		return FieldSpec{}, fmt.Errorf("unsupported field type %q", parts[1])
	}
	field := FieldSpec{Name: parts[0], GoName: pascal(parts[0]), Type: parts[1]}
	if len(parts) == 3 {
		if parts[2] != "required" {
			return FieldSpec{}, fmt.Errorf("invalid field modifier %q", parts[2])
		}
		field.Required = true
	}
	if len(parts) == 4 {
		if parts[2] != "required" || field.Type != "select" {
			return FieldSpec{}, fmt.Errorf("field options require a required select field")
		}
		field.Required = true
		options, err := parseFieldOptions(parts[3])
		if err != nil {
			return FieldSpec{}, fmt.Errorf("invalid field options: %w", err)
		}
		field.Options = options
	}
	return field, nil
}

func parseFieldOptions(value string) ([]FieldOption, error) {
	if value == "" {
		return nil, fmt.Errorf("options cannot be empty")
	}
	seen := make(map[string]struct{})
	options := make([]FieldOption, 0)
	for _, raw := range strings.Split(value, "|") {
		parts := strings.SplitN(raw, "=", 2)
		if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
			return nil, fmt.Errorf("expected value=Label pairs")
		}
		if _, exists := seen[parts[0]]; exists {
			return nil, fmt.Errorf("duplicate option %q", parts[0])
		}
		seen[parts[0]] = struct{}{}
		options = append(options, FieldOption{Value: parts[0], Label: parts[1]})
	}
	return options, nil
}

func Normalize(input Input) (Spec, error) {
	if !resourceNamePattern.MatchString(input.Name) {
		return Spec{}, fmt.Errorf("invalid resource name %q", input.Name)
	}
	if len(input.Fields) == 0 {
		return Spec{}, fmt.Errorf("resource %q requires at least one field", input.Name)
	}

	spec := Spec{
		Name:       input.Name,
		GoName:     pascal(input.Name),
		Label:      input.Label,
		Route:      input.Route,
		Permission: input.Permission,
		Icon:       input.Icon,
	}
	if spec.Label == "" {
		spec.Label = humanize(input.Name)
	}
	if spec.Route == "" {
		spec.Route = "/admin/" + input.Name
	}
	spec.FrontendRoute = strings.TrimPrefix(spec.Route, "/admin")
	if spec.FrontendRoute == "" {
		spec.FrontendRoute = "/"
	}
	if !strings.HasPrefix(spec.FrontendRoute, "/") {
		spec.FrontendRoute = "/" + strings.TrimPrefix(spec.FrontendRoute, "/")
	}
	if spec.Permission == "" {
		spec.Permission = "admin." + input.Name + ".view"
	}
	if spec.Icon == "" {
		spec.Icon = "box"
	}
	actions := input.Actions
	if len(actions) == 0 {
		actions = []string{"view", "create", "update", "delete"}
	}
	permissionSpec, err := NormalizePermission(PermissionInput{Name: input.Name, Actions: actions})
	if err != nil {
		return Spec{}, err
	}
	spec.Actions = permissionSpec.Actions

	seen := make(map[string]struct{}, len(input.Fields))
	for _, raw := range input.Fields {
		field, err := ParseField(raw)
		if err != nil {
			return Spec{}, err
		}
		if _, exists := seen[field.Name]; exists {
			return Spec{}, fmt.Errorf("duplicate field %q", field.Name)
		}
		seen[field.Name] = struct{}{}
		spec.Fields = append(spec.Fields, field)
	}
	return spec, nil
}

func pascal(value string) string {
	var builder strings.Builder
	upper := true
	for _, r := range value {
		if r == '-' || r == '_' {
			upper = true
			continue
		}
		if upper {
			builder.WriteRune(unicode.ToUpper(r))
			upper = false
			continue
		}
		builder.WriteRune(r)
	}
	return builder.String()
}

func humanize(value string) string {
	words := strings.Split(value, "-")
	for index, word := range words {
		if word == "" {
			continue
		}
		words[index] = strings.ToUpper(word[:1]) + word[1:]
	}
	return strings.Join(words, " ")
}
