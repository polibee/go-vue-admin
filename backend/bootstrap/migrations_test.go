package bootstrap

import "testing"

func TestMigrationsIncludeDepartments(t *testing.T) {

	found := false
	for _, migration := range Migrations() {
		if migration.Signature() == "create_departments_table" {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("departments migration is not registered")
	}
}
