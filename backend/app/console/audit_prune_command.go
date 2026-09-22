package console

import (
	"fmt"
	"strconv"

	"github.com/goravel/framework/contracts/console"
	"github.com/goravel/framework/contracts/console/command"

	auditservices "goravel/app/services/audit"
)

type AuditPruneCommand struct{}

func (AuditPruneCommand) Signature() string { return "admin:prune-audit-logs" }

func (AuditPruneCommand) Description() string {
	return "Delete audit logs older than the configured retention window"
}

func (AuditPruneCommand) Extend() command.Extend {
	return command.Extend{
		Category: "admin",
		Flags: []command.Flag{
			&command.StringFlag{Name: "days", Usage: "retain audit logs for this many days", Value: fmt.Sprintf("%d", auditservices.DefaultAuditRetentionDays)},
		},
	}
}

func (AuditPruneCommand) Handle(ctx console.Context) error {
	days, err := strconv.Atoi(ctx.Option("days"))
	if err != nil {
		return fmt.Errorf("invalid --days: %w", err)
	}
	deleted, cutoff, err := auditservices.NewAuditService().Cleanup(days)
	if err != nil {
		return err
	}
	ctx.Success(fmt.Sprintf("Deleted %d audit logs older than %s.", deleted, cutoff.Format("2006-01-02 15:04:05Z07:00")))
	return nil
}
