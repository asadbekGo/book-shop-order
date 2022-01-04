#!/bin/bash

CURRENT_DIR=$(pwd)

protoc -I /usr/local/include \
       -I $GOPATH/src/github.com/gogo/protobuf/gogoproto \
       -I $CURRENT_DIR/or_protos/ \
        --gofast_out=plugins=grpc:$CURRENT_DIR/genproto/ \
        $CURRENT_DIR/or_protos/*.proto;

if [[ "$OSTYPE" == "darwin"* ]]; then
    sed -i "" -e "s/,omitempty//g" $CURRENT_DIR/genproto/*.go
  else
    sed -i -e "s/,omitempty//g" $CURRENT_DIR/genproto/*.go
fi

# make proto-gen 
# cd scripts 
# ll // -rw-rw-r-- => for + x
# chmod +x gen-proto.sh  
# make proto-gen  