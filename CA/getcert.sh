#!/bin/sh

SID="spiffe\:\/\/foo\.com\/root\/$1"

sed s/XXchangeXX/${SID}/ templ.cnf > $1.cnf

openssl req -new -config ./$1.cnf -newkey rsa:2048 -nodes -keyout $1.key -out $1.csr
#openssl req -new -config ./$1.cnf -newkey dsa:dsaparams -nodes -keyout $1.key -out $1.csr

openssl ca -config ./$1.cnf -days 3650 -extensions req_ext -in $1.csr -out $1.crt
