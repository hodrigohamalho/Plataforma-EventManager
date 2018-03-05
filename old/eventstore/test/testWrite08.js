var EventStore = require("../EventStore");
const Config = require("./testConfig.js");

try {
    eventStore = new EventStore(new Config().get());

    // *************
    var event00 = {
        name : "almoco",
        payload : 
        { 
            prato : "salada", 
            preco : 17.30, 
        },
        user : 
        {
            name : "Francis June",
            id : "PL - 5874"
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