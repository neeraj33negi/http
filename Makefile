run:
	go run cmd/tcplistener/main.go
runtee:
	go run cmd/tcplistener/main.go | tee /tmp/tcp.txt
runudp:
	go run cmd/udpsender/main.go
test:
	go test ./...
test_headers:
	go test ./internal/headers/...
