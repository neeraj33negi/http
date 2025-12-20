run:
	go run cmd/tcplistener/main.go
runtee:
	go run cmd/tcplistener/main.go | tee /tmp/tcp.txt
runudp:
	go run cmd/udpsender/main.go
run_httpserver:
	go run cmd/httpserver/main.go
test:
	go test ./...
test_headers:
	go test ./internal/headers/...
