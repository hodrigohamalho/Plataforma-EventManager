var EventDb = require("../EventDb");

var config = 
{ 
    influxip : "localhost",
    database : "test002"
} 

eventDb = new EventDb(config);


var start = 1517505896487000000;
var end =   1517507178375000000;

var  promise = eventDb.findByInterval(start, end);

promise
.then((events) => { 
    var size = events.length; 
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
})
.catch((e) => {
    console.log("error = ",e)
});