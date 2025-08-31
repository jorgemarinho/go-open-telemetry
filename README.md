# Desafio OpenTelemetry - Sistema de Temperatura por CEP #

Este repositório foi criado exclusivamente para hospedar o código do desenvolvimento do Desfio de implementação do OpenTelemetry no Sistema de temperatura por CEP da Pós Go Expert, ministrado pela Full Cycle.

## Descrição do Desafio ##

A seguir estão os dados fornecidos na descrição do desafio.
Objetivo

Desenvolver um sistema em Go que receba um CEP, identifica a cidade e retorna o clima atual (temperatura em graus celsius, fahrenheit e kelvin) juntamente com a cidade. Esse sistema deverá implementar OTEL(Open Telemetry) e Zipkin.

Basedo no cenário conhecido "Sistema de temperatura por CEP" denominado Serviço B, será incluso um novo projeto, denominado Serviço A.

## Requisitos - Serviço A (responsável pelo input) ##

    O sistema deve receber um input de 8 dígitos via POST, através do schema: { "cep": "29902555" }

    O sistema deve validar se o input é valido (contem 8 dígitos) e é uma STRING

        Caso seja válido, será encaminhado para o Serviço B via HTTP

        Caso não seja válido, deve retornar:

            Código HTTP: 422

            Mensagem: invalid zipcode

## Requisitos - Serviço B (responsável pela orquestração) ##

    O sistema deve receber um CEP válido de 8 digitos

    O sistema deve realizar a pesquisa do CEP e encontrar o nome da localização, a partir disso, deverá retornar as temperaturas e formata-lás em: Celsius, Fahrenheit, Kelvin juntamente com o nome da localização.

    O sistema deve responder adequadamente nos seguintes cenários:

        Em caso de sucesso:

            Código HTTP: 200

            Response Body: { "city: "São Paulo", "temp_C": 28.5, "temp_F": 28.5, "temp_K": 28.5 }

        Em caso de falha, caso o CEP não seja válido (com formato correto):

            Código HTTP: 422

            Mensagem: invalid zipcode

        Em caso de falha, caso o CEP não seja encontrado:

            Código HTTP: 404

            Mensagem: can not find zipcode

## OTEL + Zipkin ##

Após a implementação dos serviços, adicione a implementação do OTEL + Zipkin:

    Implementar tracing distribuído entre Serviço A - Serviço B

    Utilizar span para medir o tempo de resposta do serviço de busca de CEP e busca de temperatura


## Como usar ##


Execute o script de inicialização para corrigir permissões e subir o ambiente:

```bash
chmod +x entrypoint.sh
./entrypoint.sh
```

Acesse o sistema em seu navegador ou utilizando ferramentas Postman e executar um post com cep como paramentro

http://localhost:8082/busca/cidade


Sistema A rodará na porta 8082
Sistema B rodará na porta 8080
zipkin rodará na porta 9411

Para ver o rastreamento ir no navegador e digital: http://localhost:9411


chmod +x entrypoint.sh
./entrypoint.sh


grafana 
admin
admin


http://localhost:9090
http://prometheus:9090


sudo chmod go-w ./.docker/filebeat.yml

GCP fazer
sudo chown root:root /home/jorge/projetos/go-open-telemetry-git/.docker/filebeat.yml


jorge@DESKTOP-3GP39HE:~/projetos/go-open-telemetry-git$ sudo tail -20 /var/lib/docker/containers/4e678a2b2c14b3d446470b0ce13afcdf5b66c12456bcb6391a148bf9e73a740b/4e678a2b2c14b3d446470b0ce13afcdf5b66c12456bcb6391a148bf9e73a740b-json.log
[sudo] password for jorge: 



docker logs filebeat | tail -40

sudo chmod go-w ./.docker/filebeat.yml

docker-compose restart filebeat


docker inspect service_a --format='{{.LogPath}}'


## Como visualizar logs dos containers Docker ##

1. Liste os containers ativos para obter o ID:
   ```bash
   docker ps
   ```
2. Copie o CONTAINER ID do serviço desejado (exemplo: serviço A ou B).
3. Execute o comando abaixo, substituindo <ID> pelo ID real do container:
   ```bash
   sudo tail -20 /var/lib/docker/containers/<ID>/<ID>-json.log
   ```
   Exemplo para o serviço A:
   ```bash
   sudo tail -20 /var/lib/docker/containers/4593b0c44d16/4593b0c44d16-json.log
   ```
   Exemplo para o serviço B:
   ```bash
   sudo tail -20 /var/lib/docker/containers/505818707b3d/505818707b3d-json.log
   ```