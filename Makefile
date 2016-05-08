build:
	go build

test:
	go test -v

clean:
	rm -rf *#
	rm -rf *~

clean_server_files:
	rm deadropfiles/ -rf
