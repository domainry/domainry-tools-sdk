// Package modulehost declares deployment-owned persistence for Tools modules.
package modulehost

import (
	"context"
	"github.com/domainry/domainry-orm/driver"
	"github.com/domainry/domainry-orm/migration"
	"github.com/domainry/domainry-orm/query"
	"github.com/domainry/domainry-orm/sqlhost"
)

type MigrationRegistrar interface {
	ApplyOwnedMigrations(context.Context, string, []migration.Migration) error
}

// Persistence borrows the host's pool, driver profile and sole migration ledger.
// Opening or closing Tools never opens a pool or creates its own ledger.
type Persistence struct {
	Database   sqlhost.Database
	Renderer   query.Renderer
	Profile    driver.Profile
	Migrations MigrationRegistrar
}
