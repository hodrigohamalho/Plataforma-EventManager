var EventStore = require("../EventStore");
const Config = require("./testConfig.js");

try {
    eventStore = new EventStore(new Config().get());

    // *************
    var event00 = {
        name : "ceia",
        payload : 
        { 
            prato : "iogurte", 
            preco : 32.70, 
        },
        user : 
        {
            name : "Jean Luc-Ponty",
            id : "TE - 0057"
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

