var EventHandler = require('../manager/EventHandler');

class SendEventHandler extends EventHandler {
    
    constructor(broker) {
        super();
        this.broker = broker;
    }

    receive(evento) {
        console.log("TT:"+this.broker.constructor.name);
        this.broker.send(evento);
    }
}

module.exports = SendEventHandler;