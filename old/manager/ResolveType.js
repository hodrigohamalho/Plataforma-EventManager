var BrokerController = require('../broker/local/BrokerControllerLoopback');
var SendEventHandler = require('../send/SendEventHandler');
var SendEventEndpoint = require('../send/SendEventEndpoint');
var ExecutorController = require('./ExecutorController');
var ReceiveEventHandler = require('../receive/ReceiveEventHandler');
var UpdateSubscribeEndpoint = require('../update/UpdateSubscribeEndpoint');

var resolveType = {};
resolveType.BrokerController = BrokerController;
resolveType.SendEventHandler = SendEventHandler;
resolveType.SendEventEndpoint = SendEventEndpoint;
resolveType.ExecutorController = ExecutorController;
resolveType.ReceiveEventHandler = ReceiveEventHandler;
resolveType.UpdateSubscribeEndpoint = UpdateSubscribeEndpoint;

module.exports = resolveType;