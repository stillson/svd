#!/bin/sh

openssl pkcs12 -export -in $1.crt -inkey $1.key -out $1.p12
