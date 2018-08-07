### Event Manager

Event Manager é o componente de plataforma responsável por rotear e distribuir os eventos de plataforma. Ele é o coração da mecânica de execução.

### Build

O build da aplicação é feito através de um arquivo Makefile, para buildar a aplicação execute o seguinte comando:

```sh
$ make
```

Após executar o make será criada uma pasta dist e o executável da aplicação event_manager.

### Deploy

O processo de deploy do event_manager na plataforma é feito através do installer, os componentes em Go são compilados e comitados dentro do installer então para atualizar a versão do event_manager para atualizar a versão do event_manager na plataforma utilize o seguinte comando:

```sh
$ mv dist/event_manager ~/installed_plataforma/Plataforma-Installer/Dockerfiles
$ plataforma --upgrade event_manager
```

### API

#### Enviar um evento na plataforma
```http
PUT /sendevent HTTP/1.1
Host: localhost:8081
Content-Type: application/json
Cache-Control: no-cache

{
    "name": "nome.do.event",
    "payload":{
      "key": "value"
    }
}
```



### Organização do código

1. actions
    * São as principais ações do serviço, por exemplo, fazer o split de eventos em comandos, finalizar uma instância de processo dentro outras funcionalidades;
2. api
    * É a declaração da API do event manager;
3. bus
    * É onde está implementado a integração com o rabbitmq, este pacote contém as ações básicas de operação do broker;
4. domain
    * Contém as definições de domíniom, por exemplo, definição da entidade Evento;
5. eventstore
    * Implementa a integração com o InfluxDb;
6. flow
    * Este pacote tem o router de eventos dentro do event_manager para aonde cada eventos será roteado para ser melhor tratado;
7. handlers
    * São o handlers de eventos cada tipo de evento tem um handler associado e todos eles ficam nesse pacote;
    * Neste pacote também contém o subpacote de middlewares de eventos;
8. infra
    * É um pacote com rotinas que lidam com clone de objetos, variáveis de ambiente e fábrica de objetos
9.  processor
    * É a implementação do router de eventos, um router poder conter middlwares e rotas;
10. sdk
    * São alguns clients da plataforma que são utilizados pelo event_manager;
11. util
    * É um pacote basico com algumas funções utilitárias, por exemplo, sha1 encoder;
12. vendor
    * É um pacote do Go onde ficam todas as bibliotecas de terceiros, os arquivos deste pacote jamais devem ser alterados diretamente;