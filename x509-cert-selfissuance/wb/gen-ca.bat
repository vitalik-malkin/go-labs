echo off
cls
echo ==============================================================
echo Generate CA X.509 certificate
::
openssl rand -base64 -out ca/pass.txt 14
openssl genrsa -out ca/key.pem 4096
openssl req -config openssl.cnf -days 1095 -new -x509 -key ca/key.pem -out ca/cert.cer
openssl pkcs12 -export -out ca/cert.p12 -in ca/cert.cer -inkey ca/key.pem -passout file:ca/pass.txt
::
echo ==============================================================
echo Done.
::
pause