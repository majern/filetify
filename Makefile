init:
	go get https://github.com/bufbuild/buf/cmd/buf
	brew install bufbuild/buf/buf
	go install google.golang.org/protobuf/cmd/protoc-gen-go

gen:
	rm -rf pkg/proto/*
	buf lint
	buf generate
	go mod tidy

clean:
	rm pkg/proto/**/*.go