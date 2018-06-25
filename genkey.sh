#!/usr/bin/env bash
openssl genrsa -out $1.key 2048
openssl ecparam -genkey -name secp384r1 -out $1.key
openssl req -new -x509 -sha256 -key $1.key -out $1.crt -days 3650
