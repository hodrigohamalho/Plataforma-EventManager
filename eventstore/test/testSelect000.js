var EventStore = require("../EventStore");
const Config = require("./testConfig.js");


eventStore = new EventStore(new Config().get());

// will find
var start = 1517513761893;
//var end =   1517513864517000000;
var end =   new Date().valueOf();


var  promise = eventStore.findByInterval(start, end);

promise
.then((events) => { 
    var size = events.length; 
    if (size == 0) {
        console.log("No event found");
    }
    else {
        for (var i = 0;  i < size; i++) {
            var event = events[i];
            var payload = event.payload;
            console.log("evento", i, ":"
                        , "ts =", event.timestamp
                        , ", name =", event.name
                        , ", instanceId =", event.instanceId
                        , ", payload.prato =", payload.prato
                        , ", payload.preco =", payload.preco);
        }  
    }
})
.catch((e) => {
    console.log("error = ",e)
});