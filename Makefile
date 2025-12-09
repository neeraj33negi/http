run:
	go run cmd/tcplistener/main.go
runtee:
	go run cmd/tcplistener/main.go | tee /tmp/tcp.txt
