package controllers

import (
	"testing"

	"goravel/app/core/resource"
)

func TestResourceDeleteStrategyKeepsDomainServicesBehindUnifiedRoute(t *testing.T) {
	cases := []struct {
		name string
		want string
	}{
		{name: "users", want: "users-service"},
		{name: "roles", want: "roles-service"},
		{name: "orders", want: "generic-table"},
	}
	for _, testCase := range cases {
		manifest := resource.Manifest{Name: testCase.name}
		if got := resourceDeleteStrategy(manifest); got != testCase.want {
			t.Errorf("resourceDeleteStrategy(%q) = %q, want %q", testCase.name, got, testCase.want)
		}
	}
}
