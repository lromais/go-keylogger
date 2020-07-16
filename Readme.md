# Repo Software de captura keylooger com elasticsearch

## Repo para estudos da linguagem golang e elasticsearch

## Software em desenvolvimento - keylogger em golang

Esse software é desenvolvido em Golang com a finalidade de fazer a captura das teclas digitadas em um sistema operacional linux e evniar para ElasticSearch.
Lembrando que  esse repo se fere a um estudo e não é um software para se usar em ambiente real.
Para que possa ser executado com perfeição será necessário:

* Ter um sistema linux (preferencia Debian like)
* Possuir acesso ao sudores
* Ter o docker-compose instalado 

## Funcionamento 

Deve se subir em primeiro lugar o elasticsearch, neste repositório possui um docker-compose.yml já pronto para que isso possa ser executado sem maiores dificuldades.
Apartir que o software em golang for gerado o binario deve se executar o keylogger (binario) com sudo

# Passo a Passo

## Como subir o ElasticSearch

```docker-compose up```

**Não será necessário realizar go get pois no Repo já possui go mod**

## Como Gerar o Build

`go build keyloggger.go`

## Como executar

`sudo ./keylogger`


# Como acessar o Kibana 

`http://localhost:5601`


### Para efetuar a query utilizando a console do Dev Tools execute:

`GET /index/_search?q=hostname:XXXXXX`

## Para consultar os logs no ElasticSearch via shell

`curl -XGET  http://localhost:9200/index/_search?q=hostname:xxxx`

