package gostream

import "embed"

//go:embed migrations/*.sql
var MigrationFS embed.FS
