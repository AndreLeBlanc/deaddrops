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

cleansrv:
	rm -rf deadropfiles/
	rm -rf server/deadropfiles
