module github.com/VauntDev/tqla/example/mysql

go 1.21

replace github.com/VauntDev/tqla => ../../

require (
	github.com/VauntDev/tqla v0.0.0-00010101000000-000000000000
	github.com/go-sql-driver/mysql v1.8.1
)

require filippo.io/edwards25519 v1.1.0 // indirect
