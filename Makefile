DOC = http://localhost:6060/pkg/deadrop/

build:
	go build

dev:
	./deadrop dev

fmt:
	go fmt ./...

test:
	go test ./...

testv:
	go test ./... -v

cov:
	go test ./... -cover

race:
	go test ./... -race

#echo $(DOC)
doc:
	godoc -http=:6060 & \
	google-chrome $(DOC)

clean:
	rm -rf *~
	rm -rf api/*~
	rm -rf server/*~
	rm -rf logfile
	rm -rf server/logfile

cleansrv:
	rm -rf deadropfiles/
	rm -rf server/deadropfiles
