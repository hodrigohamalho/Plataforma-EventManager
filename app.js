var Client = require('node-rest-client').Client;

var config = require('./config');
var CoreRepository = require("plataforma-sdk/services/CoreRepository");
var Processo = require("plataforma-core/Processo");
var Presentation = require("plataforma-core/Presentation");
var utils = require("plataforma-sdk/utils");

// Dependencies
// ===========================================================
var express = require("express");
var bodyParser = require("body-parser");

// Configure the Express application
// ===========================================================
var app = express();
var PORT = config.PORT;

// Set up the Express application to handle data parsing
app.use(bodyParser.urlencoded({ extended: false }));
app.use(bodyParser.json());

var coreRepository = new CoreRepository();

/**
 * Utilizado para testes no lugar do executor, enquanto não temos a url do 
 * executor definida.
 */
app.post("/testexecutor", function(req, res) {
  
    console.log("___ENTER POST EXECUTOR___" + JSON.stringify(req.body));
  
    res.send("OK");
});
  
var prefixEventSystem = "system.event.";

/**
 * Recebe os eventos para serem enviados ao executor,
 * por enquanto está sendo feito um curto circuito e enviando 
 * diretamente para o executor na URL configurada.
 */
app.post("/event", function(req, res) {

  var evento = req.body;

  console.log("___ENTER POST EVENT___" + JSON.stringify(evento));
  
  var client = new Client();

  // TODO os eventos do tipo reprodução ou reprocessamento deve fazer o curto-circuito.

  // TODO poderia ter sido feito um contain apenas
  var operations = coreRepository.getOperationsByEvent(evento.name);

  var args = { data: evento, headers: { "Content-Type": "application/json" } };

  if (operations.length > 0 || utils.isSystemEvent(evento.name)) { 
    
    let postExecution = client.post(config.executorUrl, args, function (data, response) {
      console.log("Evento "+ evento.name + " enviado para o Executor com sucesso");
    });
    postExecution.on('error', function (err) {
      console.log('request error', err);
    });
  } else {
    console.log("Sem operações para " + evento.name);
  }

  // TODO poderia ter sido feito um contain apenas
  var presentations = coreRepository.getPresentationsByEvent(evento.name);
  
  if (presentations.length > 0) {
    console.log("consulta de apresentacoes:"+ evento.name);
    
    let postExecution = client.post(config.proxyPresentationUrl, args, function (data, response) {
      console.log("Evento enviado para o Presentation com sucesso");
    });
    postExecution.on('error', function (err) {
      console.log('request error', err);
    });

  }

  res.send("OK");
});


// Listener
// ===========================================================
app.listen(PORT, function() {
  console.log("App listening on PORT " + PORT);
});