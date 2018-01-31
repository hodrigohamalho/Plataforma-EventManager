var EventDb = require("../EventDb");

var config = 
{ 
    influxip : "localhost",
    database : "test000"
} 

eventDb = new EventDb(config);

var evento = 
{
    name : "evento04",
    payload : { val1 : "valor 01", val2 : 3687.74, val3 : "valor 103" }
}

eventDb.save(evento, 
            function(id) {
                console.log("gravacao ok, id =", id)
            },
            function(err) {
                console.log("erro:", err);
            }
        );
