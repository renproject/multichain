# README

## Generate a keypair (privatekey + address)
```bash
$ go run keygen.go
DIGIBYTE_PK=L5Hm2XTZ5PZjQyhRR7U8W6qkDGqM7znzvPpQ5STTNstaiUdaoTVy
DIGIBYTE_ADDRESS=DJ2xUhLh1HsznAJL3UnBTGDtztJmSjeegL
```

## Build your docker container
```bash
docker build .
```

## Run the container
```bash
# Regtest
docker run -p 18443:18443 digibyte:latest "DJ2xUhLh1HsznAJL3UnBTGDtztJmSjeegL"
```