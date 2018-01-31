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

eventDb.save(evento, 
            function() {
                console.log("gravacao ok");
            },
            function(err) {
                console.log("db 'testXYZ' não existe, erro:", err);
            }
        );
