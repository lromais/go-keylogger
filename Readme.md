# Repo de estudo de golang

## Software em desenvolvimento - keylogger em golang

# Como subir o ElasticSearch

```docker-compose up```


# Como Gerar o Build

`go build novo.go`

# Como executar

`sudo ./novo`


# Como acessar o ElasticSearch 

`http://localhost:5601`


# Como consultar os logs no ElasticSearch

`curl -XGET  http://localhost:9200/index/_search?q=hostname:xxxx`

### ou na console do Dev Tools execute essa query

`GET /index/_search?q=hostname:XXXXXX`