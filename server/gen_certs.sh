#!/bin/bash

# recipe from
# https://stackoverflow.com/questions/64814173/how-do-i-use-sans-with-openssl-instead-of-common-name

openssl genrsa -out rootCA.key 2048
openssl req -new -x509 -days 365 -key rootCA.key -subj "/C=CN/ST=GD/L=SZ/O=Acme, Inc./CN=Acme Root CA" -out rootCA.crt

openssl req -newkey rsa:2048 -nodes -keyout server.key -subj "/C=CN/ST=GD/L=SZ/O=Acme, Inc./CN=*.localhost" -out server.csr
openssl x509 -req -extfile <(printf "subjectAltName=DNS:localhost") -days 365 -in server.csr -CA rootCA.crt -CAkey rootCA.key -CAcreateserial -out server.crt
