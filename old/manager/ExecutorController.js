var Client = require('node-rest-client').Client;
var utils = require("plataforma-sdk/utils");
var config = require('../config');

var client = new Client();

class ExecutorController {

    constructor(owner){}

    sendEventToPresentation(evento) {

        var args = { data: evento, headers: { "Content-Type": "application/json" } };

        let postExecution = client.post(config.proxyPresentationUrl, args, function (data, response) {
            console.log("Evento " + evento.name +" enviado para o Presentation com sucesso");
        });
        postExecution.on('error', function (err) {
            console.log('request error', err);
        });

    }

    sendEventToProcess(evento) {

        var args = { data: evento, headers: { "Content-Type": "application/json" } };

        let postExecution = client.put(config.executorUrl, args, function (data, response) {
          console.log("Evento "+ evento.name + " enviado para o Executor com sucesso");
        });
        postExecution.on('error', function (err) {
          console.log('request error', err);
        });
    }

}

module.exports = ExecutorController;