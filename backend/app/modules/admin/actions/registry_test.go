package actions

import (
	"testing"

	"github.com/goravel/framework/contracts/http"
	adminactions "goravel/app/core/admin/actions"
)

func TestRegistryIncludesModuleActionHandlers(t *testing.T) {
	handler, err := Registry().Find("user-status")
	if err != nil {
		t.Fatalf("find user-status handler: %v", err)
	}
	if handler.Payload() != "user-status" {
		t.Fatalf("payload = %q, want user-status", handler.Payload())
	}
}

func TestRegisterRejectsDuplicateActionHandler(t *testing.T) {
	if err := Register(adminactionsTestHandler{kind: "user-status"}); err != adminactions.ErrDuplicateHandler {
		t.Fatalf("register duplicate error = %v, want %v", err, adminactions.ErrDuplicateHandler)
	}
}

type adminactionsTestHandler struct{ kind string }

func (h adminactionsTestHandler) Kind() string  { return h.kind }
func (adminactionsTestHandler) Payload() string { return "test" }
func (adminactionsTestHandler) Execute(_ http.Context, _ adminactions.Request) (adminactions.Result, error) {
	return adminactions.Result{}, nil
}
