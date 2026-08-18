// Package migrations nhúng các file .sql vào binary để chạy migration lúc khởi động.
package migrations

import "embed"

//go:embed *.sql
var FS embed.FS
