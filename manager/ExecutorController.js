var Client = require('node-rest-client').Client;
var CoreRepository = require("plataforma-sdk/services/CoreRepository");
var utils = require("plataforma-sdk/utils");
var config = require('../config');

var client = new Client();
var coreRepository = new CoreRepository();

class ExecutorController {

    constructor(owner){}

    hasPresentationWaitEvent(evento) {
        var presentations = coreRepository.getPresentationsByEvent(evento.name);
        return presentations.length > 0;
    }

    sendEventToPresentation(evento) {

        var args = { data: evento, headers: { "Content-Type": "application/json" } };

        let postExecution = client.post(config.proxyPresentationUrl, args, function (data, response) {
            console.log("Evento " + evento.name +" enviado para o Presentation com sucesso");
        });
        postExecution.on('error', function (err) {
            console.log('request error', err);
        });

    }

    hasProcessWaitEvent(evento) {
        var operations = coreRepository.getOperationsByEvent(evento.name);
        return operations.length > 0 || utils.isSystemEvent(evento.name);
    }
        
    sendEventToProcess(evento) {
        
        var args = { data: evento, headers: { "Content-Type": "application/json" } };
        
        let postExecution = client.post(config.executorUrl, args, function (data, response) {
          console.log("Evento "+ evento.name + " enviado para o Executor com sucesso");
        });
        postExecution.on('error', function (err) {
          console.log('request error', err);
        });        
    }

}

module.exports = ExecutorController;