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
        
        if (this.executor.hasProcessWaitEvent(evento)) {
            this.executor.sendEventToProcess(evento);
        } else {
            console.log("Sem operações para " + evento.name);
        }

        if (this.executor.hasPresentationWaitEvent(evento)) {
            this.executor.sendEventToPresentation(evento);
        } else {
            console.log("Sem apresentação para " + evento.name);
        }

    }
}

module.exports = ReceiveEventHandler;