#!/usr/bin/env bash

dateFormat=$(date '+%Y-%m-%d_%H:%M:%S')
executableName='goSafeChatBackend'
newExecutableName="new_$executableName"
oldExecutableName="old_$executableName"

echo 'Compiling the project locally'
./build.sh
chmod +x $executableName

echo "Uploading the new executable $newExecutableName"
scp -i ~/.ssh/pakpak.pem goSafeChatBackend ubuntu@ec2-35-165-24-103.us-west-2.compute.amazonaws.com:~/$newExecutableName

ssh -i ~/.ssh/pakpak.pem ubuntu@ec2-35-165-24-103.us-west-2.compute.amazonaws.com << EOF
set +e
# Check if the executable exists
if [ -f $executable ]; then
	mkdir -p oldBuilds
	mv $executable oldBuilds/$oldExecutableName
fi
mv $newExecutableName $executable

docker rm -f safechat_docker
docker build -t safechat .
docker run --name="safechat_docker" -d -t safechat

EOF