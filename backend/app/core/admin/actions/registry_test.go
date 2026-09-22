package actions

import (
	"testing"

	"github.com/goravel/framework/contracts/http"
)

type fakeHandler struct{ kind string }

func (h fakeHandler) Kind() string  { return h.kind }
func (fakeHandler) Payload() string { return "test-payload" }

func (h fakeHandler) Execute(_ http.Context, request Request) (Result, error) {
	return Result{Action: request.Action, Requested: len(request.IDs), Succeeded: len(request.IDs)}, nil
}

func TestNormalizeRequestDeduplicatesAndRejectsInvalidIDs(t *testing.T) {
	request, err := NormalizeRequest(Request{Action: "set-status", IDs: []int64{2, 2, 3}})
	if err != nil {
		t.Fatalf("expected duplicate IDs to normalize, got %v", err)
	}
	if len(request.IDs) != 2 || request.IDs[0] != 2 || request.IDs[1] != 3 {
		t.Fatalf("unexpected normalized IDs: %#v", request.IDs)
	}
	if _, err := NormalizeRequest(Request{Action: "set-status", IDs: nil}); err != ErrEmptyIDs {
		t.Fatalf("expected empty IDs error, got %v", err)
	}
	if _, err := NormalizeRequest(Request{Action: "set-status", IDs: []int64{0}}); err != ErrInvalidID {
		t.Fatalf("expected invalid ID error, got %v", err)
	}
}

func TestNormalizeRequestRejectsTooManyIDs(t *testing.T) {
	ids := make([]int64, MaxIDs+1)
	for index := range ids {
		ids[index] = int64(index + 1)
	}
	if _, err := NormalizeRequest(Request{Action: "set-status", IDs: ids}); err != ErrTooManyIDs {
		t.Fatalf("expected too many IDs error, got %v", err)
	}
}

func TestRegistryRejectsDuplicateAndFindsHandler(t *testing.T) {
	registry := NewRegistry()
	handler := fakeHandler{kind: "user-status"}
	if err := registry.Register(handler); err != nil {
		t.Fatalf("register failed: %v", err)
	}
	if err := registry.Register(handler); err != ErrDuplicateHandler {
		t.Fatalf("expected duplicate handler error, got %v", err)
	}
	if _, err := registry.Find("missing"); err != ErrHandlerNotFound {
		t.Fatalf("expected missing handler error, got %v", err)
	}
}

func TestNormalizeRequestCarriesPayload(t *testing.T) {
	request, err := NormalizeRequest(Request{Action: "set-status", IDs: []int64{2}, Payload: map[string]any{"status": "active"}})
	if err != nil || request.Payload["status"] != "active" {
		t.Fatalf("expected payload to be preserved, request=%+v err=%v", request, err)
	}
}

func TestRegistryRejectsHandlerWithoutPayloadContract(t *testing.T) {
	registry := NewRegistry()
	if err := registry.Register(fakeHandler{kind: ""}); err != ErrHandlerNotFound {
		t.Fatalf("expected invalid handler to be rejected, got %v", err)
	}
}
