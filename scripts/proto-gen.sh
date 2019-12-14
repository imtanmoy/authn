protoc -I $GOPATH/src --go_out=$GOPATH/src $GOPATH/src/github.com/imtanmoy/authy/api/protos/organization.proto
protoc -I $GOPATH/src --go_out=plugins=grpc:$GOPATH/src $GOPATH/src/github.com/imtanmoy/authy/api/protos/organization.proto
