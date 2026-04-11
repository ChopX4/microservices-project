module github.com/ChopX4/raketka/assembly

go 1.25.0

require (
	github.com/ChopX4/raketka/platform v0.0.0-00010101000000-000000000000
	github.com/ChopX4/raketka/shared v0.0.0-00010101000000-000000000000
	github.com/IBM/sarama v1.47.0
	github.com/caarlos0/env/v11 v11.4.0
	github.com/google/uuid v1.6.0
	github.com/joho/godotenv v1.5.1
	go.uber.org/zap v1.27.1
	google.golang.org/protobuf v1.36.11
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/eapache/go-resiliency v1.7.0 // indirect
	github.com/eapache/queue v1.1.0 // indirect
	github.com/hashicorp/go-uuid v1.0.3 // indirect
	github.com/jcmturner/aescts/v2 v2.0.0 // indirect
	github.com/jcmturner/dnsutils/v2 v2.0.0 // indirect
	github.com/jcmturner/gofork v1.7.6 // indirect
	github.com/jcmturner/gokrb5/v8 v8.4.4 // indirect
	github.com/jcmturner/rpc/v2 v2.0.3 // indirect
	github.com/klauspost/compress v1.18.4 // indirect
	github.com/pierrec/lz4/v4 v4.1.25 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/rcrowley/go-metrics v0.0.0-20250401214520-65e299d6c5c9 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/crypto v0.48.0 // indirect
	golang.org/x/net v0.51.0 // indirect
)

replace github.com/ChopX4/raketka/platform => ../platform

replace github.com/ChopX4/raketka/shared => ../shared
