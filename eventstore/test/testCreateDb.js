var EventStore = require("../EventStore");
const Config = require("./testConfig.js");


try {

    eventStore = new EventStore(new Config().get());

    eventStore.createDb()
        .then((created) => { 
            console.log(created);
        })
        .catch((e) => {
            console.log("error save = '",e,"'")
        });
    }
catch(e) {
    console.log("error test =",e);
}