var EventStore = require("../EventStore");
const Config = require("./testConfig.js");

try {
    eventStore = new EventStore(new Config().get());

    // *************
    var event00 = {
        name : "jantar",
        payload : 
        { 
            prato : "peixe", 
            preco : 32.70, 
        },
        user : 
        {
            name : "Martin Barry",
            id : "MS - 9987"
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