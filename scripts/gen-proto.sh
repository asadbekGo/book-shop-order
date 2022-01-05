#!/bin/bash

CURRENT_DIR=$(pwd)


protoc -I /usr/local/include \
       -I $GOPATH/src/github.com/gogo/protobuf/gogoproto \
       -I $CURRENT_DIR/protos/order_service/ \
        --gofast_out=plugins=grpc:$CURRENT_DIR/genproto/order_service/ \
        $CURRENT_DIR/protos/order_service/*.proto;

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