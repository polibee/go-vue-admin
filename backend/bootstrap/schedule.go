package bootstrap

import (
	"github.com/goravel/framework/contracts/schedule"

	"goravel/app/facades"
)

// Schedule registers unattended, retention-only maintenance tasks. Destructive
// scopes (selected, filtered, and all) remain interactive and are never
// scheduled by default.
func Schedule() []schedule.Event {
	return []schedule.Event{
		facades.Schedule().Command("admin:prune-audit-logs --days=365").DailyAt("02:30").Name("audit-log-retention").SkipIfStillRunning(),
	}
}
