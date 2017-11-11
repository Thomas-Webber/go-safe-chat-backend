#!/usr/bin/env bash

PORT=8090
dateFormat=$(date '+%Y-%m-%d_%H:%M:%S')
executableName='goSafeChatBackend'
newExecutableName="new_$executableName"
oldExecutableName="old_$executableName"

export GOPATH=$(pwd)

echo 'Compiling the project locally'
./build.sh
chmod +x $executableName

echo "Uploading the new executable $newExecutableName"
mv $executableName $newExecutableName
scp -i ~/.ssh/pakpak.pem $newExecutableName Dockerfile ubuntu@ec2-35-165-24-103.us-west-2.compute.amazonaws.com:~
mv $newExecutableName $executableName

ssh -i ~/.ssh/pakpak.pem ubuntu@ec2-35-165-24-103.us-west-2.compute.amazonaws.com << EOF
set +e
# Check if the executable exists
if [ -f $executable ]; then
	mkdir -p oldBuilds
	mv $executable oldBuilds/$oldExecutableName
fi
mv $newExecutableName $executableName

sudo docker rm -f safechat_docker
sudo docker build -t safechat .
sudo docker run --name="safechat_docker" -d -p $PORT:$PORT -t safechat
sudo docker logs safechat_docker

EOF