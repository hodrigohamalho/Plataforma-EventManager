class SendEventEndpoint  {

  constructor(app, handler) {
    this.handler = handler;
    this.createEndpoint(app, handler);
  }

  createEndpoint(app, handler) {

    console.log("createEndpoint");
    
    app.put("/sendevent", function(req, res) {
      
        var evento = req.body;
      
        console.log("___ENTER PUT EVENT___" + JSON.stringify(evento));

        handler.receive(evento);
        
        res.send("OK");
    });

  }

}

module.exports = SendEventEndpoint;