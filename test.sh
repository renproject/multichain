source ./infra/.env
GO=/usr/local/go/bin/go #DEZU: Added because the script hate on poor GO
COMPOSE_PARALLEL_LIMIT=10 #DEZU: Added because of error, apparently some library has a thread limit of 10?
docker-compose -f ./infra/docker-compose.yaml up --build -d
echo "Waiting for multichain to boot..."
sleep 30
echo "Done waiting, running tests..." # DEZU: Added for debugging
$GO version # DEZU: Added for de lulz
$GO test -v ./...
echo "Done testing, closing..." # DEZU: Added for debugging
docker-compose -f ./infra/docker-compose.yaml down
echo "Done!"
