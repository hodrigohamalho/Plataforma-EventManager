var Client = require('node-rest-client').Client;

var config = require('./config');
var CoreRepository = require("plataforma-sdk/services/CoreRepository");
var Processo = require("plataforma-core/Processo");
var Presentation = require("plataforma-core/Presentation");

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
  
/**
 * Recebe os eventos para serem enviados ao executor,
 * por enquanto está sendo feito um curto circuito e enviando 
 * diretamente para o executor na URL configurada.
 */
app.post("/event", function(req, res) {

  console.log("___ENTER POST EVENT___" + JSON.stringify(req.body));
  
  var client = new Client();

  // TODO poderia ter sido feito um contain apenas
  var operations = coreRepository.getOperationsByEvent(req.body.name);

  var args = { data: req.body, headers: { "Content-Type": "application/json" } };

  if (operations.length > 0) { 
    
    var reqExec = client.post(config.executorUrl, args, function (data, response) {
      console.log("Evento enviado para o Executor com sucesso");
    });
    reqExec.on('error', function (err) {
      console.log('request error', err);
    });
  } else {
    console.log("Sem operações para " + req.body.name);
  }

  // TODO poderia ter sido feito um contain apenas
  var presentations = coreRepository.getPresentationsByEvent(req.body.name);
  
  if (presentations.length > 0) {
    
    var reqExec = client.post(config.proxyPresentationUrl, args, function (data, response) {
      console.log("Evento enviado para o Presentation com sucesso");
    });
    reqExec.on('error', function (err) {
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