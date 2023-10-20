go-build: go-generate
	go build -o target/kafka cmd/core/main.go

go-generate:
	go generate -x 'mini-kafka-go/pkg/protocol'

