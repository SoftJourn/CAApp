# CAApp
Part of https://github.com/SoftJourn/sj_coin_fabric_POC project. 
Basicaly designed to generate Hyperledger fabric user certificates and public/private keys.

Instalation 
Use CAApp/docker/install.sh to run application in docker container

Contains several parts:
- Web application (uses LDAP user to create and deploy identity file and keys for NodeJS fabric client application)
URL: 
  http://localhost:3000/login

- Face login API
URLS: 
  http://localhost:3000/api/register
  http://localhost:3000/api/login

- Certificate API
URL:
  http://localhost:3000/api/certificate
