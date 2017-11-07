docker build -t="caapp" .

docker run --network=fixtures_default --mount type=bind,source=/tmp/fabric-client-kvs_peerOrg1,target=/tmp/fabric-client-kvs_peerOrg1 -i -t -p 3005:3000 caapp
