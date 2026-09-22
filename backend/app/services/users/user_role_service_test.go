package userservices

import (
	"testing"

	"github.com/goravel/framework/database/orm"

	"goravel/app/models"
)

func TestRoleSetChanged(t *testing.T) {
	role := models.Role{Model: orm.Model{ID: 1}}
	if roleSetChanged([]models.Role{role}, []models.Role{{Model: orm.Model{ID: 1}}}) {
		t.Fatal("expected identical role sets to be unchanged")
	}
	if !roleSetChanged([]models.Role{role}, []models.Role{{Model: orm.Model{ID: 2}}}) {
		t.Fatal("expected different role sets to be changed")
	}
}
