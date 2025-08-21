#!/bin/bash
set -e

# Corrige permissão do arquivo de configuração do Filebeat
if [ -f ./.docker/filebeat.yml ]; then
  sudo chown root:root ./.docker/filebeat.yml
fi

# Corrige permissão da pasta de logs para o Otel Collector
if [ -d ./.docker/tmp ]; then
  sudo chown -R 1000:1000 ./.docker/tmp
  sudo chmod -R 777 ./.docker/tmp
fi

docker-compose up -d --build
