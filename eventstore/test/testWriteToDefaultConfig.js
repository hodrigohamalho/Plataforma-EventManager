var EventStore = require("../EventStore");


try {
    eventStore = new EventStore();

    // *************
    var event00 = {
        name : "assaltoageladeira",
        payload : 
        { 
            prato : "bolo com sorvete", 
            preco : 157.80, 
        },
        user : 
        {
            name : "Joe Satriani",
            id : "GM - 1173"
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

