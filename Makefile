all: check cassini

cassini: protoc
	GO111MODULE=on go build .

protoc:
	protoc \
		-I api/ \
		-I $(GOPATH)/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis/ \
		--go_out=plugins=grpc:api \
		api/*.proto
	protoc \
		-I api/ \
		-I $(GOPATH)/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis/ \
		--grpc-gateway_out=logtostderr=true:api \
		api/*.proto

check: test

test:
	GO111MODULE=on go test -v github.com/embarkstudios/cassini/...
