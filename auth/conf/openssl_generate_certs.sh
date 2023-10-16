#!/bin/bash

# Function to interactively request input
prompt_input() {
    read -rp "$1: " "$2"
}

# Set the full path of the script
script_path=$(realpath "$0")
script_dir=$(dirname "$script_path")

# Set the full path of the script
script_path=$(realpath "$0")
script_dir=$(dirname "$script_path")

# Set file names based on the arguments passed to the script or requested interactively
if [[ "$1" != "--no-interactive" ]]; then
    prompt_input "Enter the domain name" domain
    prompt_input "Enter the environment" environment
else
    domain=$2
    environment=$3
fi

# Check if the file ${domain}.conf exists in the script's folder
if [ ! -f "${script_dir}/${domain}.conf" ]; then
    echo "Error: The file ${domain}.conf is not present in the script's folder." >&2
    exit 1
fi

# Create the domain folder if it doesn't already exist
mkdir -p "${domain}"

# Move into the domain folder
cd "${domain}" || exit

# Create the root key
openssl genrsa -out "${domain}-${environment}-CA-key.pem" 2048

# Generate the certificate request for the root key
openssl req -x509 -new -key "${domain}-${environment}-CA-key.pem" -sha256 -days 1825 -out "${domain}-${environment}-CA.pem"

# Generate the dev certificate key
openssl genrsa -out "${domain}-${environment}.key" 2048

# Generate the dev certificate request
openssl req -new -sha256 -out "${domain}-${environment}.csr" -key "${domain}-${environment}.key" -config "${script_dir}/${domain}.conf"

# Generate the dev certificate using the root key
openssl x509 -req -in "${domain}-${environment}.csr" -CA "${domain}-${environment}-CA.pem" -CAkey "${domain}-${environment}-CA-key.pem" -CAcreateserial -out "${domain}-${environment}.crt" -days 1825 -sha256 -extfile "${script_dir}/${domain}.conf" -extensions req_ext
