all : clean ensure build 
clean:
	rm -rf ./bin
ensure: 
	dep ensure
build:
	GOOS=windows go build -o bin/windows_amd64/tailor.exe main.go