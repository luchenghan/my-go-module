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

for PORT in 7004 7005 7006; do
    BAK_CONFIG_FILE="./$PORT/sentinel.bak"
    CONFIG_FILE="./$PORT/sentinel.conf"
    if [ -f "$CONFIG_FILE" ]; then
        cp "$BAK_CONFIG_FILE" "$CONFIG_FILE"
        find . -name '*.conf' -print0 | xargs -0 sed -i "" "s/IP/${IP}/g"
    else
        echo "Configuration file $CONFIG_FILE not found!"
        exit 1
    fi
done

docker compose -f docker-compose-sentinel.yaml up -d