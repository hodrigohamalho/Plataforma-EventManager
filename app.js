// Dependencies
// ===========================================================
var express = require("express");
var bodyParser = require("body-parser");
var resolveType = require("./manager/ResolveType");
var config = require('./config');

// Configure the Express application
// ===========================================================
var app = express();
var PORT = config.PORT;

// Set up the Express application to handle data parsing
app.use(bodyParser.urlencoded({ extended: false }));
app.use(bodyParser.json());

var broker = new resolveType.BrokerController(app);
var sendHandler = new resolveType.SendEventHandler(broker);
var sendEndpoint = new resolveType.SendEventEndpoint(app, sendHandler);

var receiveHandler = new resolveType.ReceiveEventHandler(broker);
var updateSubscribeEndpoint = new resolveType.UpdateSubscribeEndpoint(app, broker);

// Listener
// ===========================================================
app.listen(PORT, function() {
  console.log("App listening on PORT " + PORT);
});