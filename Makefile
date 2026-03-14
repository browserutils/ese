all:
	go build -o eseparser cmd/eseparser/*.go


windows:
	GOOS=windows GOARCH=amd64 \
            go build \
	    -o eseparser.exe ./cmd/eseparser/*.go

generate:
	cd parser/ && binparsegen conversion.spec.yaml > ese_gen.go


test:
	go test ./...
