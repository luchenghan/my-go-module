#!/bin/bash

OS=$(uname -s)

if [ "$OS" == "Linux" ]; then
    IP=$(hostname -I | awk '{print $1}')
elif [ "$OS" == "Darwin" ]; then
    IP=$(ifconfig | grep 'en0' -A 3 | grep 'inet ' | awk '{print $2}')
else
    echo "Unsupported OS: $OS"
    exit 1
fi

for PORT in 7001 7002 7003 7004 7005 7006; do
    BAK_CONFIG_FILE="./$PORT/cluster.bak"
    CONFIG_FILE="./$PORT/cluster.conf"
    if [ -f "$CONFIG_FILE" ]; then
        cp "$BAK_CONFIG_FILE" "$CONFIG_FILE"
        find . -name '*.conf' -print0 | xargs -0 sed -i "" "s/IP/${IP}/g"
    else
        echo "Configuration file $CONFIG_FILE not found!"
        exit 1
    fi
done

# Start
docker compose -f docker-compose-cluster.yaml up -d
# Create cluster
sleep 1
echo "[Info] Setup redis cluster ..."
docker exec -it redis1 sh -c "echo 'yes' | redis-cli --cluster create $IP:7001 $IP:7002 $IP:7003 $IP:7004 $IP:7005 $IP:7006 --cluster-replicas 1"
# Check cluster
sleep 2
echo "[Info] Check cluster health ..."
docker exec -it redis1 sh -c "redis-cli -c -p 7001 CLUSTER INFO"
