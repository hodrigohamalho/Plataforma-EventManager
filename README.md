# Plataforma-EventManager

#### Introdução
O módulo eventmanager é responsável pelo roteamento de mensagens de eventos da plataforma.
Neste caso todos os eventos disparados por qualquer outro módulo é direcionado para o eventmanager, que será responsável por enviar para o evento para o broker, e 
da mesma forma receber todas as mensagens esperadas pelos sistemas cadastrados na plataforma.
As mensagens podem ser recebidas para atender uma tela da camada de presentation ou de processo de negócio.
Um tipo de mensagem especial esperado pelo gerenciador de eventos, são os eventos de sistema, como para: reproduzir, reprocessar e etc..

OBS: Os eventos de sistema são identificados por um nome de prefixo: 'system.event.'

OBS: No momento a versão do eventmanager não está enviando mensagens para o Broker, ele está devolvendo os eventos diretamente para a aplicação.

#### Estrutura do Projeto
No projeto podemos encontrar os seguintes arquivos:
* [app.js]: disponibiliza os serviços de recebimento de evento.
* [config.js]: contém as configurações do eventmanager

#### Requisitos

Para executar o gerenciador com sucesso você precisa instalar as seguintes ferramentas:
* [NodeJS](https://nodejs.org)
* NPM (vem junto com o NodeJS)
* [Docker](https://www.docker.com/)
* Docker compose

### Para instalar ou atualizar as dependências é necessário executar o comando:
npm install


Caso queira executar o servidor sem utilizar o docker, tem um script no projeto Plataforma-SDK, em:
Plataforma-SDK/_scripts/shell/start-eventManager.sh

Se você estiver utilizando o windows, é necessário executar o powersheel no modo terminal.

Caso você opte por usar o docker você pode subir com o seguinte comando:
```sh
$ docker-compose up -d
```
Ao executar esse comando o docker irá subir um container com EventManager inicializado.

Após a subida dos containers você pode enviar deve acessar o eventmamanger pelo endereço:
http://localhost:8081/event

Example:
Url: http://localhost:8081/event
Http Method: POST

    Body: 
        Business Event: {"name":"account.put","processName":"cadastra-conta","payload":{"id":0,"titular":"fernando","saldo":100},"origem":"38913f27-b09c-4240-b13b-b46db0e52591"}
        System Event: {"name":"system.event.reproduction","payload":{"instanciaOriginal":"f6648fdb-6d5a-79d8-5852-0a2d892e9a3c"}}
