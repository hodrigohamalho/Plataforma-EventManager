var EventStore = require("../EventStore");
const Config = require("./testConfig.js");


eventStore = new EventStore(new Config().get());

// will not find
var start = 1517514814; 
var end =   1517514815; 


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