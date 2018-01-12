class UpdateSubscribeEndpoint  {

  constructor(app, broker) {
    this.createEndpoint(app, broker);
    this.broker = broker;
  }

  createEndpoint(app, broker) {

      /**
       * Recebe os eventos para serem enviados ao executor,
       * por enquanto está sendo feito um curto circuito e enviando 
       * diretamente para o executor na URL configurada.
       */
      app.put("/updatesubscribe", function(req, res) {
        
          var evento = req.body;
        
          console.log("___ENTER UPDATE EVENT___" + JSON.stringify(evento));

          broker.updateSubscribe(evento);
          
          res.send("OK");
      });

  }

}

module.exports = UpdateSubscribeEndpoint;
