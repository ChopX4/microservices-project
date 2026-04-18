module github.com/ChopX4/raketka/iam

go 1.25.0

require (
	github.com/ChopX4/raketka/platform v0.0.0-00010101000000-000000000000
	github.com/ChopX4/raketka/shared v0.0.0-00010101000000-000000000000
	github.com/caarlos0/env/v11 v11.4.0
	github.com/gomodule/redigo v1.9.3
	github.com/google/uuid v1.6.0
	github.com/jackc/pgx/v5 v5.9.1
	github.com/joho/godotenv v1.5.1
	github.com/pressly/goose v2.7.0+incompatible
	go.uber.org/zap v1.27.1
	golang.org/x/crypto v0.48.0
	google.golang.org/grpc v1.79.3
)

require (
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/net v0.51.0 // indirect
	golang.org/x/sync v0.19.0 // indirect
	golang.org/x/sys v0.41.0 // indirect
	golang.org/x/text v0.34.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20251202230838-ff82c1b0f217 // indirect
	google.golang.org/protobuf v1.36.11 // indirect
)

replace github.com/ChopX4/raketka/platform => ../platform

replace github.com/ChopX4/raketka/shared => ../shared
