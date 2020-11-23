#source ./infra/.env
. ./infra/.env # DEZU: This is only because Ubuntu is a special boi
GO=/usr/local/go/bin/go #DEZU: Added because the script hate on poor GO
COMPOSE_PARALLEL_LIMIT=8 #DEZU: Added because of error, apparently some library has a thread limit?
docker-compose -f ./infra/docker-compose.yaml up --build -d
echo "Waiting for multichain to boot..."
docker logs --details infra_qtum_1 >> qtumdockerlog.txt
sleep 30
docker logs --details infra_qtum_1 >> qtumdockerlog.txt
echo "Done waiting, running tests on Qtum..." # DEZU: Added for debugging
$GO test -v ./chain/qtum/...
docker logs --details infra_qtum_1 >> qtumdockerlog.txt
echo "Done with Qtum, running all tests..." # DEZU: Added for debugging
$GO test -v ./...
echo "Done testing, closing..." # DEZU: Added for debugging
docker-compose -f ./infra/docker-compose.yaml down
echo "Done!"
