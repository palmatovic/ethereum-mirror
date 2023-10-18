#!/bin/bash

# Ensure all required arguments are provided
if [ "$#" -ne 6 ]; then
  echo "Usage: $0 <COUNTRY> <STATE> <LOCALITY> <ORG_NAME> <ORG_UNIT> <COMMON_NAME>"
  exit 1
fi

# Assign the arguments to variables
COUNTRY="$1"
STATE="$2"
LOCALITY="$3"
ORG_NAME="$4"
ORG_UNIT="$5"
COMMON_NAME="$6"

# Create a lowercase version of ORG_NAME with spaces replaced by underscores
ORG_NAME_LOWER=$(echo "$ORG_NAME" | tr ' ' '_' | tr '[:upper:]' '[:lower:]')

# Generate the certificate configuration file name
CONFIG_FILE_NAME="${ORG_NAME_LOWER}_openssl.cnf"

# Create the certificate configuration file
cat << EOF > "$CONFIG_FILE_NAME"
[ req ]
default_bits       = 2048
distinguished_name = req_distinguished_name
req_extensions     = req_ext

[ req_distinguished_name ]
countryName                 = $COUNTRY
stateOrProvinceName         = $STATE
localityName                = $LOCALITY
organizationName            = $ORG_NAME
organizationalUnitName      = $ORG_UNIT
commonName                  = $COMMON_NAME

[ req_ext ]
subjectAltName = @alt_names

[alt_names]
DNS.1 = $COMMON_NAME
EOF

echo "Configuration file '$CONFIG_FILE_NAME' generated successfully."
