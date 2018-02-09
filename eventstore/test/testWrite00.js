var EventStore = require("../EventStore");
const Config = require("./testConfig.js");

try {
    eventStore = new EventStore(new Config().get());

    // *************
    var event00 = {
        name : "cafedamanha",
        payload : 
        { 
            prato : "suco de melancia", 
            preco : 7.20, 
        },
        user : 
        {
            name : "Maria das Neves 6, a moça nota = 10",
            id : "PV - 4587"
        }    
    }
    
    eventStore.save(event00)
        .then((instance_id) => { 
            console.log("instance id = ", instance_id);
        })
        .catch((e) => {
            console.log("error save = ",e);
        });
}
catch(e) {
    console.log("error test =",e);
}