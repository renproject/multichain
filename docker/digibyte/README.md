# README

## Generate a keypair (privatekey + address)
```bash
$ go run keygen.go
DIGIBYTE_PK=eagPs6RBxmTQyjni3K7vqPNBwjN4o5R8CEwP4eyHavMJMz29MCen
DIGIBYTE_ADDRESS=smtdQvMJRLaWwNaFUjdBtFzUR4evxQJcB9
```

## Build your docker container
```bash
docker build .
```

## Run the container
```bash
# Regtest
docker run -p 18443:18443 digibyte:latest "smtdQvMJRLaWwNaFUjdBtFzUR4evxQJcB9"
```