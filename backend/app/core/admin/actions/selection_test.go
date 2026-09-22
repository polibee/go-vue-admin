package actions

import "testing"

func TestSelectionValidate(t *testing.T) {
	for name, selection := range map[string]Selection{
		"explicit ids":          {Mode: "ids", IDs: []int64{1, 2}},
		"filtered query":        {Mode: "query", Query: map[string]string{"status": "active"}},
		"query with exclusions": {Mode: "query", Query: map[string]string{"search": "alice"}, ExcludeIDs: []int64{2}},
	} {
		t.Run(name, func(t *testing.T) {
			if err := selection.Validate(); err != nil {
				t.Fatalf("valid selection rejected: %v", err)
			}
		})
	}
	for name, selection := range map[string]Selection{
		"missing mode":         {IDs: []int64{1}},
		"ids without ids":      {Mode: "ids"},
		"ids mixed with query": {Mode: "ids", IDs: []int64{1}, Query: map[string]string{"status": "active"}},
		"empty query":          {Mode: "query"},
	} {
		t.Run(name, func(t *testing.T) {
			if err := selection.Validate(); err == nil {
				t.Fatal("invalid selection accepted")
			}
		})
	}
}
