var EventStore = require("../EventStore");
const Config = require("./testConfig.js");


eventStore = new EventStore(new Config().get());
var evento = 
{
    name : "evento00",
    payload : -74
}


var  promise = eventStore.save(evento);

promise
.then((instance_id) => { 
    console.log("instance id = ", instance_id);
})
.catch((e) => {
    console.log("error = ",e)
});