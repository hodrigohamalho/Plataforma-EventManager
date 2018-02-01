var EventDb = require("../EventDb");

var config = 
{ 
    influxip : "localhost",
    database : "test002"
} 

eventDb = new EventDb(config);

var evento = 
{
    name : "assalto",
    payload : 
    { 
        prato : "pudim", 
        preco : 0.50, 
    }
}

var  promise = eventDb.save(evento);

promise
.then((instance_id) => { 
    console.log("instance id = ", instance_id);
})
.catch((e) => {
    console.log("error = ",e)
});