packName=musicBot-linux-x86
rm -rf ${packName}*
make build
mkdir ${packName}
mv musiccloud-bot ${packName}
cp -rf ./3rd/tools-linux-x86/*  ${packName}
cp ./config.yml.sample ${packName}/config.yml
cp -rf ./shell/*  ${packName}
time=$(date "+%Y%m%dT%H%M%S")
tar cvzf ${packName}-${time}.tar.gz ${packName}/ && rm -rf ${packName}

