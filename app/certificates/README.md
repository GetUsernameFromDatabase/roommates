# Certificates

Files here are used for ssl (https)

## Install Openssl

Please visit <https://github.com/openssl/openssl> to get pkg and install.

## Generate RSA private key

```sh
openssl genrsa -out ./server.key 2048
```

## Generate digital certificate

```sh
openssl req -new -x509 -key ./server.key -out ./server.pem -days 365
```
