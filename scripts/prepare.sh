#!/bin/bash
mkdir run
mkdir run/resource
cp bin/gin-app-server.tar.gz run/
cp -r config run/
cp deployments/docker-compose.yaml run/
echo  "#!/bin/bsh
docker load -i gin-app-server.tar.gz
" >> run/load.sh