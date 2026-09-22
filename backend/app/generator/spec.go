package generator

import (
	"fmt"
	"regexp"
	"strconv"
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
	Name             string
	Label            string
	Route            string
	PageMode         string
	Permission       string
	Icon             string
	Actions          []string
	ActionSpecs      []string
	ActionFieldSpecs []string
	Fields           []string
	Scope            string
	OwnerField       string
	Relations        []string
	FormGroups       []string
	Details          []string
}

type Spec struct {
	Name          string
	GoName        string
	Label         string
	Route         string
	FrontendRoute string
	PageMode      string
	Permission    string
	Icon          string
	Actions       []string
	ActionSpecs   []ActionSpec
	Fields        []FieldSpec
	DataScope     string
	OwnerField    string
	Relations     []RelationSpec
	FormGroups    []FormGroupSpec
	Details       []DetailSectionSpec
}

type ActionSpec struct {
	Name          string
	Label         string
	Kind          string
	Permission    string
	Batch         bool
	Payload       string
	PayloadFields []ActionPayloadFieldSpec
}

type ActionPayloadFieldSpec struct {
	Name     string
	Label    string
	Type     string
	Required bool
	Options  []FieldOption
}

type RelationSpec struct {
	Name, Kind, Resource, Field, ForeignField, LabelField string
	Selectable                                            bool
}

type FormGroupSpec struct {
	Name, Label string
	Columns     int
	Fields      []string
}

type DetailSectionSpec struct {
	Name, Label string
	Fields      []string
}

type FieldSpec struct {
	Name             string
	GoName           string
	Type             string
	Required         bool
	Options          []FieldOption
	Visible          *bool
	Readable         *bool
	Writable         *bool
	Sensitive        bool
	PolicyConfigured bool
}

type FieldOption struct {
	Value string
	Label string
}

func ParseField(value string) (FieldSpec, error) {
	parts := strings.Split(value, ":")
	if len(parts) < 2 || len(parts) > 5 || !fieldNamePattern.MatchString(parts[0]) {
		return FieldSpec{}, fmt.Errorf("invalid field %q: expected name:type[:required[:value=Label|value=Label]]", value)
	}
	if _, ok := supportedFieldTypes[parts[1]]; !ok {
		return FieldSpec{}, fmt.Errorf("unsupported field type %q", parts[1])
	}
	field := FieldSpec{Name: parts[0], GoName: pascal(parts[0]), Type: parts[1]}
	visible, readable, writable := true, true, true
	field.Visible, field.Readable, field.Writable = &visible, &readable, &writable
	modifiers := parts[2:]
	if len(modifiers) > 0 && strings.Contains(modifiers[len(modifiers)-1], "=") {
		if field.Type != "select" {
			return FieldSpec{}, fmt.Errorf("field options require a select field")
		}
		options, err := parseFieldOptions(modifiers[len(modifiers)-1])
		if err != nil {
			return FieldSpec{}, fmt.Errorf("invalid field options: %w", err)
		}
		field.Options = options
		modifiers = modifiers[:len(modifiers)-1]
	}
	seenModifiers := make(map[string]struct{}, len(modifiers))
	for _, modifier := range modifiers {
		if _, exists := seenModifiers[modifier]; exists || modifier == "" {
			return FieldSpec{}, fmt.Errorf("invalid field modifier %q", modifier)
		}
		seenModifiers[modifier] = struct{}{}
		switch modifier {
		case "required":
			field.Required = true
		case "sensitive":
			field.Sensitive = true
			field.PolicyConfigured = true
		case "readonly":
			writable = false
			field.Writable = &writable
			field.PolicyConfigured = true
		case "hidden":
			visible, readable = false, false
			field.Visible, field.Readable = &visible, &readable
			field.PolicyConfigured = true
		default:
			return FieldSpec{}, fmt.Errorf("invalid field modifier %q", modifier)
		}
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
	spec.PageMode = input.PageMode
	if spec.PageMode == "" {
		spec.PageMode = "generic"
	}
	if spec.PageMode != "generic" && spec.PageMode != "custom" {
		return Spec{}, fmt.Errorf("unsupported page mode %q", spec.PageMode)
	}
	spec.DataScope = input.Scope
	if spec.DataScope == "" {
		spec.DataScope = "all"
	}
	if spec.DataScope != "all" && spec.DataScope != "own" {
		return Spec{}, fmt.Errorf("unsupported data scope %q", spec.DataScope)
	}
	spec.OwnerField = input.OwnerField
	if spec.DataScope == "own" && spec.OwnerField == "" {
		return Spec{}, fmt.Errorf("own data scope requires an owner field")
	}
	actions := append([]string(nil), input.Actions...)
	if len(actions) == 0 {
		actions = []string{"view", "create", "update", "delete"}
	}
	customActions := make(map[string]ActionSpec)
	customActionNames := make([]string, 0, len(input.ActionSpecs))
	for _, raw := range input.ActionSpecs {
		action, err := parseActionSpec(raw)
		if err != nil {
			return Spec{}, err
		}
		if _, exists := customActions[action.Name]; exists {
			return Spec{}, fmt.Errorf("duplicate action %q", action.Name)
		}
		customActions[action.Name] = action
		customActionNames = append(customActionNames, action.Name)
	}
	for _, raw := range input.ActionFieldSpecs {
		field, err := parseActionPayloadFieldSpec(raw)
		if err != nil {
			return Spec{}, err
		}
		action, exists := customActions[field.Action]
		if !exists {
			return Spec{}, fmt.Errorf("action payload field %q references undeclared action %q", field.Name, field.Action)
		}
		for _, existing := range action.PayloadFields {
			if existing.Name == field.Name {
				return Spec{}, fmt.Errorf("duplicate payload field %q", field.Name)
			}
		}
		action.PayloadFields = append(action.PayloadFields, ActionPayloadFieldSpec{Name: field.Name, Label: field.Label, Type: field.Type, Required: field.Required, Options: field.Options})
		customActions[field.Action] = action
	}
	permissionSpec, err := NormalizePermission(PermissionInput{Name: input.Name, Actions: actions})
	if err != nil {
		return Spec{}, err
	}
	spec.Actions = permissionSpec.Actions
	for _, actionName := range customActionNames {
		if !fieldNamePattern.MatchString(actionName) {
			return Spec{}, fmt.Errorf("invalid action name %q", actionName)
		}
		if !contains(spec.Actions, actionName) {
			spec.Actions = append(spec.Actions, actionName)
		}
	}
	for _, action := range spec.Actions {
		if custom, exists := customActions[action]; exists {
			spec.ActionSpecs = append(spec.ActionSpecs, custom)
			continue
		}
		spec.ActionSpecs = append(spec.ActionSpecs, ActionSpec{Name: action, Label: humanize(action), Permission: "admin." + spec.Name + "." + action})
	}

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
	if spec.DataScope == "own" {
		_, exists := seen[spec.OwnerField]
		if !exists {
			return Spec{}, fmt.Errorf("owner field %q must be declared", spec.OwnerField)
		}
		for _, field := range spec.Fields {
			if field.Name == spec.OwnerField && field.Type != "integer" {
				return Spec{}, fmt.Errorf("owner field %q must use integer type", spec.OwnerField)
			}
		}
	}
	for _, raw := range input.Relations {
		relation, err := parseRelation(raw)
		if err != nil {
			return Spec{}, err
		}
		if relation.Kind == "belongsTo" {
			if _, exists := seen[relation.Field]; !exists {
				return Spec{}, fmt.Errorf("relation %q field %q must be declared", relation.Name, relation.Field)
			}
		} else if relation.Selectable {
			return Spec{}, fmt.Errorf("hasMany relation %q cannot be selectable", relation.Name)
		}
		spec.Relations = append(spec.Relations, relation)
	}
	for _, raw := range input.FormGroups {
		group, err := parseFormGroup(raw)
		if err != nil || group.Columns < 1 || group.Columns > 4 {
			return Spec{}, fmt.Errorf("invalid form group %q", raw)
		}
		if err := validateSpecFieldReferences(group.Fields, seen); err != nil {
			return Spec{}, err
		}
		spec.FormGroups = append(spec.FormGroups, group)
	}
	for _, raw := range input.Details {
		section, err := parseDetailSection(raw)
		if err != nil {
			return Spec{}, fmt.Errorf("invalid detail section %q", raw)
		}
		if err := validateSpecFieldReferences(section.Fields, seen); err != nil {
			return Spec{}, err
		}
		spec.Details = append(spec.Details, section)
	}
	return spec, nil
}

func parseRelation(value string) (RelationSpec, error) {
	parts := strings.Split(value, ":")
	if len(parts) < 6 || len(parts) > 7 || parts[0] == "" || parts[1] == "" || parts[2] == "" || parts[3] == "" || parts[4] == "" || parts[5] == "" {
		return RelationSpec{}, fmt.Errorf("invalid relation %q", value)
	}
	if parts[1] != "belongsTo" && parts[1] != "hasMany" {
		return RelationSpec{}, fmt.Errorf("invalid relation kind %q", parts[1])
	}
	selectable := len(parts) == 7 && parts[6] == "selectable"
	return RelationSpec{Name: parts[0], Kind: parts[1], Resource: parts[2], Field: parts[3], ForeignField: parts[4], LabelField: parts[5], Selectable: selectable}, nil
}

func parseActionSpec(value string) (ActionSpec, error) {
	parts := strings.Split(value, ":")
	if len(parts) != 6 || parts[0] == "" || parts[1] == "" || parts[2] == "" || parts[3] == "" {
		return ActionSpec{}, fmt.Errorf("invalid action %q: expected name:label:kind:permission:batch:payload", value)
	}
	batch, err := strconv.ParseBool(parts[4])
	if err != nil || (batch && parts[5] == "") || (!batch && parts[5] != "") {
		return ActionSpec{}, fmt.Errorf("invalid action %q: batch and payload contract do not match", value)
	}
	return ActionSpec{Name: parts[0], Label: parts[1], Kind: parts[2], Permission: parts[3], Batch: batch, Payload: parts[5]}, nil
}

type actionPayloadFieldInput struct {
	Action string
	ActionPayloadFieldSpec
}

func parseActionPayloadFieldSpec(value string) (actionPayloadFieldInput, error) {
	parts := strings.SplitN(value, ":", 6)
	if len(parts) < 5 || parts[0] == "" || !fieldNamePattern.MatchString(parts[1]) || parts[2] == "" {
		return actionPayloadFieldInput{}, fmt.Errorf("invalid action field %q: expected action:name:label:type:required[:value=Label|value=Label]", value)
	}
	supported := map[string]bool{"text": true, "number": true, "boolean": true, "select": true}
	if !supported[parts[3]] {
		return actionPayloadFieldInput{}, fmt.Errorf("unsupported action payload field type %q", parts[3])
	}
	required, err := strconv.ParseBool(parts[4])
	if err != nil {
		return actionPayloadFieldInput{}, fmt.Errorf("invalid action field %q: required must be boolean", value)
	}
	var options []FieldOption
	if len(parts) == 6 && parts[5] != "" {
		if parts[3] != "select" {
			return actionPayloadFieldInput{}, fmt.Errorf("action field options require a select field")
		}
		options, err = parseFieldOptions(parts[5])
		if err != nil {
			return actionPayloadFieldInput{}, err
		}
	}
	return actionPayloadFieldInput{Action: parts[0], ActionPayloadFieldSpec: ActionPayloadFieldSpec{Name: parts[1], Label: parts[2], Type: parts[3], Required: required, Options: options}}, nil
}

func contains(values []string, value string) bool {
	for _, item := range values {
		if item == value {
			return true
		}
	}
	return false
}

func parseFormGroup(value string) (FormGroupSpec, error) {
	parts := strings.SplitN(value, ":", 4)
	if len(parts) != 4 || parts[0] == "" || parts[1] == "" {
		return FormGroupSpec{}, fmt.Errorf("invalid form group")
	}
	columns, err := strconv.Atoi(parts[2])
	if err != nil || parts[3] == "" {
		return FormGroupSpec{}, fmt.Errorf("invalid form group")
	}
	return FormGroupSpec{Name: parts[0], Label: parts[1], Columns: columns, Fields: strings.Split(parts[3], "|")}, nil
}

func parseDetailSection(value string) (DetailSectionSpec, error) {
	parts := strings.SplitN(value, ":", 3)
	if len(parts) != 3 || parts[0] == "" || parts[1] == "" || parts[2] == "" {
		return DetailSectionSpec{}, fmt.Errorf("invalid detail section")
	}
	return DetailSectionSpec{Name: parts[0], Label: parts[1], Fields: strings.Split(parts[2], "|")}, nil
}

func validateSpecFieldReferences(names []string, fields map[string]struct{}) error {
	seen := make(map[string]struct{}, len(names))
	for _, name := range names {
		if _, ok := fields[name]; !ok {
			return fmt.Errorf("field %q must be declared", name)
		}
		if _, ok := seen[name]; ok {
			return fmt.Errorf("duplicate field %q", name)
		}
		seen[name] = struct{}{}
	}
	return nil
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
