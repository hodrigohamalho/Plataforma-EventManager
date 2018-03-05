var EventStore = require("../EventStore");
const Config = require("./testConfig.js");

try {
    eventStore = new EventStore(new Config().get());

    // *************
    var event00 = {
        name : "colacao",
        payload : 
        { 
            prato : "iogurte", 
            preco : 4.30, 
        },
        user : 
        {
            name : "Solange Inês, III",
            id : "VD - 8547"
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
