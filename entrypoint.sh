#!/bin/bash
set -e

# Corrige permissão do arquivo de configuração do Filebeat
if [ -f ./.docker/filebeat.yml ]; then
  sudo chown root:root ./.docker/filebeat.yml
fi

docker-compose up -d --build
