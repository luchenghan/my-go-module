#!/bin/bash

inti(){
  rm -rf .env
  cp  .env.template .env
}

linux(){
  case $1 in 'down')
        linuxDown
      ;;
      'up')
        linuxUp
      ;;
      'pull')
        docker compose pull
      ;;
      '')
        linuxDown
        linuxUp
      ;;
  esac
}

linuxUp(){
  addr=$(hostname -I | awk '{print $1}')
  find . -name '*.env' -exec sed -i "s/IP_ADDR/${addr}/g" {} +
  uuid=$(uuidgen | tr -d '-' | base64 | cut -b 1-22)
  find . -name '*.env' -exec sed -i "s/UUID/${uuid}/g" {} +
  docker compose up --build -d
}

linuxDown(){
  docker container prune -f
  docker volume prune -f
  docker buildx prune -f
  docker compose down
}

mac(){
  case $1 in 'down')
        macDown
      ;;
      'up')
        macUp
      ;;
      'pull')
          docker compose pull
      ;;
      '')
        macDown
        macUp
      ;;
  esac
}

macUp(){
  amd64="amd64" #for linux
  arm64="arm64" # for mac
  find . -name '*.env' -print0 | xargs -0 sed -i "" "s/${amd64}/${arm64}/g"
  addr=$(ipconfig getifaddr en0)
  find . -name '*.env' -print0 | xargs -0 sed -i "" "s/IP_ADDR/${addr}/g"
  uuid=$(uuidgen | tr -d '-' | base64 | cut -b 1-22)
  find . -name '*.env' -print0 | xargs -0 sed -i "" "s/UUID/${uuid}/g"
  docker compose up --build -d
}

macDown(){
  docker container prune -f
  docker volume prune -f
  docker buildx prune -f
  docker compose down
}

checkOS(){
  os=$(uname -s)
  case "${os}" in 'Linux')
       linux $1
      ;;
      'Darwin')
       mac $1
      ;;
  esac
}

inti
checkOS $1
