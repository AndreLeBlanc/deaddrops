build:
	go build

test:
	go test ./...

testv:
	go test ./... -v

fmt:
	go fmt ./...

clean:
	rm -rf *~
	rm -rf api/*~
	rm -rf server/*~

clean_server_files:
	rm deadropfiles/ -rf
