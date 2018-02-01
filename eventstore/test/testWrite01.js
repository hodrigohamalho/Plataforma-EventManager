var EventDb = require("../EventDb");

var config = 
{ 
    influxip : "localhost",
    database : "testXYZ"
} 

eventDb = new EventDb(config);

var evento = 
{
    name : "evento00",
    payload : -74
}


var  promise = eventDb.save(evento);

promise
.then((instance_id) => { 
    console.log("instance id = ", instance_id);
})
.catch((e) => {
    console.log("error = ",e)
});