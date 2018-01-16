var EventHandler = require('../manager/EventHandler');
var ExecutorController = require("../manager/ExecutorController");

class ReceiveEventHandler extends EventHandler {

    constructor(broker) {
        super();
        this.broker = broker;
        this.broker.setReceiveHandler(this);
        this.executor = new ExecutorController();
    }

    receive(evento) {
      this.executor.sendEventToProcess(evento);
    }
}

module.exports = ReceiveEventHandler;
