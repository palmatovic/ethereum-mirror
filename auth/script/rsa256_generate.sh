#!/bin/bash

# Check if the input parameter is provided
if [ -z "$1" ]; then
echo "You must specify a base name for the keys."
exit 1
fi

# Base name for the keys, with spaces replaced by underscores and converted to lowercase
input=$(echo "$1" | tr ' ' '_' | tr '[:upper:]' '[:lower:]')

# Name of the file for the private key
private_key_file="${input}_private_key.pem"

# Name of the file for the public key
public_key_file="${input}_public_key.pem"

# Generate a 2048-bit RSA private key
openssl genpkey -algorithm RSA -out "$private_key_file" -aes256

# Generate the corresponding public key
openssl rsa -pubout -in "$private_key_file" -out "$public_key_file"

echo "RSA private key generated successfully: $private_key_file"
echo "Corresponding public key generated successfully: $public_key_file"
