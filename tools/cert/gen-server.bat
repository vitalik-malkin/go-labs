echo off
cls
echo ==============================================================
echo Generate server X.509 certificate
::
set subject-display-name=localhost01
::
set cert-directory=server-%subject-display-name%
set cert-file-base-name=server-%subject-display-name%
::
echo ====================
echo subject-display-name: %subject-display-name%
echo ---
echo cert-directory: %cert-directory%
echo ---
echo cert-file-base-name: %cert-file-base-name%
echo ---
echo ====================
::
openssl rand -base64 -out %cert-directory%/%cert-file-base-name%-pass.txt 14
openssl genrsa -out %cert-directory%/%cert-file-base-name%-key.pem 4096
openssl req -new -key %cert-directory%/%cert-file-base-name%-key.pem -out %cert-directory%/%cert-file-base-name%-req.pem -outform PEM -nodes -config openssl.cnf
openssl ca -config openssl.cnf -in %cert-directory%/%cert-file-base-name%-req.pem -out %cert-directory%/%cert-file-base-name%-cert.cer -notext -extensions server_extensions
openssl pkcs12 -export -out %cert-directory%/%cert-file-base-name%-cert.p12 -in %cert-directory%/%cert-file-base-name%-cert.cer -inkey %cert-directory%/%cert-file-base-name%-key.pem -passout file:%cert-directory%/%cert-file-base-name%-pass.txt
::
echo ==============================================================
echo Done.
::
pause