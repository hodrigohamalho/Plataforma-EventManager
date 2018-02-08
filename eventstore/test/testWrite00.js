var EventStore = require("../EventStore");
const Config = require("./testConfig.js");


eventStore = new EventStore(new Config().get());

// *************
eventStore.save({ name : "cafedamanha", payload : { prato : "suco", preco : 6.50 } })
.then((instance_id) => { 
    console.log("instance id = ", instance_id);
})
.catch((e) => {
    console.log("error = ",e)
});
/* 

// *************
eventStore.save({ name : "colacao", payload : { prato : "bananas", preco : 5.20 } })
.then((instance_id) => { 
    console.log("instance id = ", instance_id);
})
.catch((e) => {
    console.log("error = ",e)
});

// *************
eventStore.save({ name : "almoco", payload : { prato : "salada", preco : 17.30 } })
.then((instance_id) => { 
    console.log("instance id = ", instance_id);
})
.catch((e) => {
    console.log("error = ",e)
});

// *************
eventStore.save({ name : "lanche", payload : { prato : "misto quente", preco : 9.80 } })
.then((instance_id) => { 
    console.log("instance id = ", instance_id);
})
.catch((e) => {
    console.log("error = ",e)
});

// *************
eventStore.save({ name : "jantar", payload : { prato : "peixe", preco : 32.70 } })
.then((instance_id) => { 
    console.log("instance id = ", instance_id);
})
.catch((e) => {
    console.log("error = ",e)
});

// *************
eventStore.save({ name : "ceia", payload : { prato : "iogurte", preco : 2.90 } })
.then((instance_id) => { 
    console.log("instance id = ", instance_id);
})
.catch((e) => {
    console.log("error = ",e)
});

// *************
eventStore.save({ name : "cafedamanha", payload : { prato : "salada de frutas", preco : 8.40 } })
.then((instance_id) => { 
    console.log("instance id = ", instance_id);
})
.catch((e) => {
    console.log("error = ",e)
});

// *************
eventStore.save({ name : "colacao", payload : { prato : "biscoitos", preco : 5.60 } })
.then((instance_id) => { 
    console.log("instance id = ", instance_id);
})
.catch((e) => {
    console.log("error = ",e)
});

// *************
eventStore.save({ name : "almoco", payload : { prato : "churrasco", preco : 25.10 } })
.then((instance_id) => { 
    console.log("instance id = ", instance_id);
})
.catch((e) => {
    console.log("error = ",e)
});

// *************
eventStore.save({ name : "almoco", payload : { prato : "churrasco", preco : 25.10 } })
.then((instance_id) => { 
    console.log("instance id = ", instance_id);
})
.catch((e) => {
    console.log("error = ",e)
});

// *************
eventStore.save({ name : "lanche", payload : { prato : "bolo", preco : 4.50 } })
.then((instance_id) => { 
    console.log("instance id = ", instance_id);
})
.catch((e) => {
    console.log("error = ",e)
});

// *************
eventStore.save({ name : "jantar", payload : { prato : "frango xadrez", preco : 24.70 } })
.then((instance_id) => { 
    console.log("instance id = ", instance_id);
})
.catch((e) => {
    console.log("error = ",e)
});

// *************
eventStore.save({ name : "ceia", payload : { prato : "cerveja!!", preco : 45.80 } })
.then((instance_id) => { 
    console.log("instance id = ", instance_id);
})
.catch((e) => {
    console.log("error = ",e)
}); */