echo ==============================================================
echo Generate client X.509 certificate

subject_display_name="deploy-service-component-2"

cert_directory="client-$subject_display_name"

echo ====================
echo subject_display-name: $subject_display_name
echo ---
echo cert_directory: $cert_directory
echo ---
echo ====================

openssl rand -base64 -out "./$cert_directory/pass.txt" 14
openssl genrsa -out "./$cert_directory/key.pem" 4096
openssl req -new -key "./$cert_directory/key.pem" -out "./$cert_directory/req.pem" -outform PEM -nodes -config openssl.cnf
openssl ca -config openssl.cnf -in "./$cert_directory/req.pem" -out "./$cert_directory/cert.cer" -notext -extensions client_extensions
openssl pkcs12 -export -out "./$cert_directory/cert.p12" -in "./$cert_directory/cert.cer" -inkey "./$cert_directory/key.pem" -passout file:"./$cert_directory/pass.txt"
openssl x509 -fingerprint -in "./$cert_directory/cert.cer" -out "./$cert_directory/fingerprint.txt"

echo ==============================================================
echo Done.