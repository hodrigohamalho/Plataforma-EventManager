var Client = require('node-rest-client').Client;
var BrokerController = require('../../manager/BrokerController');
var config = require('../../config');

var client = new Client();


class BrokerControllerLoopback extends BrokerController {

    constructor(app) {
        super();
        this.createReceiveEventLoopback(app);
    }

    setReceiveHandler(handler) {
        this.handler = handler;    
    }

    createReceiveEventLoopback(app) {
        
        var broker = this;

        app.put("/event", function(req, res) {
            
            var evento = req.body;
        
            console.log("Evento "+ evento.name + " recebido em loopback.");

            broker.handler.receive(evento);
            
            res.send("OK");
        });
    }

    send(message) {

        var evento = message;
        
        var args = { data: evento, headers: { "Content-Type": "application/json" } };

        var urlReceive = `http://localhost:${config.PORT}/event`;

        let postExecution = client.put(urlReceive, args, function (data, response) {
            console.log(`Evento ${evento.name} enviado em loopback.`);
        });
        postExecution.on('error', function (err) {
            console.log('request error', err);
        });
    }
    
    receive(message) {
        this.handler.receive(message);
    }

    updateSubscribe(message) {
        console.log("Atualizada a lista de eventos que o gerenciador escuta." + JSON.stringify(message));
    }
}

module.exports = BrokerControllerLoopback;
