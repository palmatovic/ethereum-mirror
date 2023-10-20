#!/bin/bash

CA_KEY="CA-key.pem"
CA_CERT="CA.pem"
SERVER_KEY="server.key"
SERVER_CSR="server.csr"
SERVER_CERT="server.crt"

# Controlla che vengano forniti almeno 1 parametro (il file di configurazione)
if [ "$#" -ne 1 ]; then
  echo "Usage: $0 <config_file>"
  exit 1
fi

CONFIG_FILE="$1"

# Creazione della chiave CA senza password
openssl genrsa -out "$CA_KEY" 2048

# Creazione del certificato CA
openssl req -x509 -new -nodes -key "$CA_KEY" -sha256 -days 1825 -out "$CA_CERT"

# Creazione della chiave del server
openssl genrsa -out "$SERVER_KEY" 2048

# Creazione del certificato di richiesta del server
openssl req -new -sha256 -out "$SERVER_CSR" -key "$SERVER_KEY" -config "$CONFIG_FILE"

# Creazione del certificato del server
openssl x509 -req -in "$SERVER_CSR" -CA "$CA_CERT" -CAkey "$CA_KEY" -CAcreateserial -out "$SERVER_CERT" -days 1825 -sha256 -extfile "$CONFIG_FILE" -extensions req_ext

echo "Server certificate generated: $SERVER_CERT"
