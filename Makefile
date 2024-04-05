all: windows linux

windows:
	GOOS=windows GOARCH=amd64 go build -ldflags "-H 'windowsgui' -w -s" -trimpath -o windows_x64_payload.exe payload/main.go

linux:
	GOOS=linux GOARCH=amd64 go build -ldflags "-w -s" -trimpath -o linux_x64_payload payload/main.go

clean:
	rm windows_x64_payload.exe
	rm linux_x64_payload