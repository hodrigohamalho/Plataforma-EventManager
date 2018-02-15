var EventStore = require("../EventStore");
const Config = require("./testConfig.js");

eventStore = new EventStore(new Config().get());

var name = "PV - 4587";
var  promise = eventStore.findByUserId(name);

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
                        , ", payload.preco =", payload.preco
                        , ", user name = ", event.user.name
                        , ", user id = ", event.user.id);
        }  
        
    }
})
.catch((e) => {
    console.log("error = ",e)
});