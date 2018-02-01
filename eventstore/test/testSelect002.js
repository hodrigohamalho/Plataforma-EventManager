var EventDb = require("../EventDb");

var config = 
{ 
    influxip : "localhost",
    database : "test003"
} 

eventDb = new EventDb(config);

// will not find
var start = 1517514814000; 
var end =   1517514815000; 


var  promise = eventDb.findByInterval(start, end);

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