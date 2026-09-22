package console

import "testing"

func TestAuditPruneCommandMetadata(t *testing.T) {
	command := AuditPruneCommand{}
	if command.Signature() != "admin:prune-audit-logs" {
		t.Fatalf("signature = %q", command.Signature())
	}
	if command.Description() == "" || len(command.Extend().Flags) != 1 {
		t.Fatalf("unexpected command metadata: %+v", command.Extend())
	}
}
