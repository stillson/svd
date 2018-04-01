#!/usr/bin/env bash 
#
# This script generates a new dummy CA certificate and key for use in the
# SPIRE development environment. Note that it will place the generated certificate
# and key in the configuration directory, replacing any existing dummy certificates.
#

###
# Constants
###
OPENSSL=/Users/stillson/openssl/openssl-1.0.2o/apps/openssl
CAKEY=ROOT
INTKEY=INT
SID="spiffe\:\/\/foo\.com"

###
# set up directory
###
mkdir -p castuff
mkdir -p castuff/certsdb
touch castuff/index.txt
echo "unique_subject = no" > castuff/index.txt.attr
echo "00000001" > castuff/serial

###
# set up configuration files
###
sed s/XXchangeXX/${SID}/ ca-templ.cnf > tmp.cnf
sed s/XXchangeXX/${SID}/ int-templ.cnf > int-tmp.cnf

###
# openssl calls
###

# generate root
$OPENSSL genrsa -out ${CAKEY}.key 2048
#openssl dsaparam -out dsaparams 2048
#$OPENSSL ecparam -name secp521r1 -genkey -noout -out ${CAKEY}.key
$OPENSSL req -new -x509 -key ${CAKEY}.key -out ${CAKEY}.crt -days 1825 -subj "/C=US/ST=/L=/O=SPIFFE/OU=/CN=/"  -config ./tmp.cnf -extensions 'req_ext'

#generate intermediates
openssl req -new -config ./int-tmp.cnf -newkey rsa:2048 -nodes -keyout ${INTKEY}.key -out ${INTKEY}.csr
#openssl req -new -config ./int-tmp.cnf -newkey dsa:dsaparams -nodes -keyout ${INTKEY}.key -out ${INTKEY}.csr
openssl ca -config ./int-tmp.cnf -days 3650 -extensions req_ext -in ${INTKEY}.csr -out ${INTKEY}.crt

rm tmp.cnf
rm int-tmp.cnf
