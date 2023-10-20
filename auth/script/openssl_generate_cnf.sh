#!/bin/bash

if [ "$#" -ne 7 ]; then
  echo "Error: Please provide 7 parameters: COUNTRY STATE LOCALITY ORG_NAME ORG_UNIT COMMON_NAME ALT_DNS"
  exit 1
fi

COUNTRY="$(echo $1 | tr '[:lower:]' '[:upper:]')"
STATE="$(echo $2 | tr '[:lower:]' '[:upper:]')"
LOCALITY="$(echo $3 | tr '[:lower:]' '[:upper:]')"
ORG_NAME="$4"
ORG_UNIT="$5"
COMMON_NAME="$6"
ALT_DNS="$7"

# Create the ssl.cnf configuration file
cat <<EOL > ssl.cnf
[ req ]
default_bits       = 2048
distinguished_name = req_distinguished_name
req_extensions     = req_ext

[ req_distinguished_name ]
countryName                 = $COUNTRY
countryName_default         = $COUNTRY
stateOrProvinceName         = $STATE
stateOrProvinceName_default = $STATE
localityName                = $LOCALITY
localityName_default        = $LOCALITY
organizationName            = $ORG_NAME
organizationName_default    = $ORG_NAME
organizationalUnitName        = $ORG_UNIT
organizationalUnitName_default = $ORG_UNIT
commonName                  = $COMMON_NAME
commonName_max              = 64
commonName_default          = $COMMON_NAME

[ req_ext ]
subjectAltName = @alt_names

[alt_names]
DNS.1 = $ALT_DNS
EOL

echo "Configuration file created: openssl.cnf"
