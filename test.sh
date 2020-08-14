source ./docker/docker-compose.env
docker-compose -f ./docker/docker-compose.yaml up --build -d
echo "Waiting for multichain to boot..."
sleep 30
go test -v ./...
docker-compose -f ./docker/docker-compose.yaml down
echo "Done!"