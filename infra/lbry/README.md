# README

## Generate a keypair (privatekey + address)

```bash
$ go run keygen.go
LBRY_PK=cV9WB2d7VTCJkkJy5RrzkQvqtiFJvAXbL8XUte8YRBH9zG2u1LTu
LBRY_ADDRESS=n3qZBVhiMHW6vxizTdpnBRpvHXHqer3Q2x
```

## Build your docker container

```bash
docker build .
```

## Run the container

```bash
# Regtest
docker run -p 29245:29245 lbrycrd:latest "n3qZBVhiMHW6vxizTdpnBRpvHXHqer3Q2x"
```
