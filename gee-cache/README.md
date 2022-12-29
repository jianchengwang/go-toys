wget https://github.com/protocolbuffers/protobuf/releases/download/v21.12/protoc-21.12-linux-x86_64.zip

unzip -d /usr/local protoc-21.12-linux-x86_64.zip

go install github.com/golang/protobuf/protoc-gen-go@latest

protoc --go_out=. *.proto