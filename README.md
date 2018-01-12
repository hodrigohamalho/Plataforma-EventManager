# Plataforma-EventManager

#### Introdução
O módulo eventmanager é responsável pelo roteamento de mensagens de eventos da plataforma.
Neste caso todos os eventos disparados por qualquer outro módulo é direcionado para o eventmanager, 
que será responsável por enviar para o evento para o broker, e da mesma forma receber todas as mensagens esperadas pelos 
sistemas cadastrados na plataforma.
As mensagens podem ser recebidas para atender uma tela da camada de presentation ou de processo de negócio.
Um tipo de mensagem especial esperado pelo gerenciador de eventos, são os eventos de sistema, como para: reproduzir, reprocessar e etc..

OBS: Os eventos de sistema são identificados por um nome de prefixo: 'system.event.'

OBS: No momento a versão do eventmanager não está enviando mensagens para o Broker, ele está devolvendo os eventos diretamente para a aplicação.
     Utilizando uma implementação de BokerControllerLoopback.

#### Estrutura do Projeto
No projeto podemos encontrar os seguintes arquivos:
* [app.js]: inicializa o serviço web da aplicação.
* [config.js]: contém as configurações do eventmanager

Pacotes:
* [broker]:
    - Neste pacote se encontra a implementação do broker de mensagem, ou seja, a comunicação do broker é realizada nesse componente.
      [local]: BrokerControllerLoopback - Implementação para a mensagem retornar para o próprio sistema.
      Obs: Para alterar a implementação a ser utilizada deve ser alterado no ResolveType
* [manager]:
    - Neste pacote ficam as classes de definição geral e gestão dos eventos. O ExecutorController faz a comunicação e envio das mensagens para o executor.
* [send]:
    - Neste pacote ficam as classes de envio de mensagens, para o broker, dos eventos recebidos do processapp e presentationapp. são ´endpoint´ e ´handler´.
* [receive]:
    - Neste pacote ficam as classes de recebimento de mensagens do broker, que serão enviadas para o executor, ´handler´.


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
http://localhost:8081/

Para enviar evento o caminho seria: http://localhost:8081/sendevent

### Atualização de eventos subescritos:

O eventmanager deve se subscrever para receber os eventos dos sistemas registrados na plataforma, ou seja, deve receber todas as mensagens
que são de interesse dos processapps e presentations. Para tanto, quando o eventmanager a cada novo deploy deve receber da plataforma
a lista de novos eventos para monitorar.

Url para receber atualizações de eventos para o eventmanager escutar: http://localhost:8081/updatesubscribe


### Example:

// Evento de negócio
Url: http://localhost:8081/sendevent
Http Method: PUT

    Body: 
        {"name":"account.put","processName":"cadastra-conta","payload":{"id":0,"titular":"fernando","saldo":100},"origem":"38913f27-b09c-4240-b13b-b46db0e52591"}

// Evento de reprodução
Url: http://localhost:8081/sendevent
Http Method: PUT

    Body: 
        {"name":"system.event.reproduction","payload":{"instanciaOriginal":"f6648fdb-6d5a-79d8-5852-0a2d892e9a3c"}}
        
// Evento de atualização de subscrição
Url: http://localhost:8081/updatesubscribe
Http Method: PUT

    Body: 
        ["event_1","event_2"]
