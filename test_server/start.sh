#!/usr/bin/env bash

export PMSG_HOST=0.0.0.0
export PMSG_PORT=5050
export PMSG_SSL_CERT=../pmsg.crt
export PMSG_SSL_KEY=../pmsg.key

cd ../test_client
go build
cd ../test_server
go build
./test_server