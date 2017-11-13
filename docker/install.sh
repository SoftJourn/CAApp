docker build -t="caapp" .

docker run --name=caapp --network=fixtures_default --mount type=bind,source=/tmp/fabric-client-kvs_peerOrg1,target=/tmp/fabric-client-kvs_peerOrg1 -i -t -p 3000:3000 caapp
