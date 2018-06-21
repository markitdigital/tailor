all : clean ensure test build 
clean:
	rm -rf ./bin
test:
	go test --cover
ensure: 
	dep ensure
build:
	GOOS=windows go build -o bin/windows_amd64/tailor.exe cmd/tailor.go