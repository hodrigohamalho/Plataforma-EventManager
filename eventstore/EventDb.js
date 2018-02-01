var unirest = require("unirest");
const uuid = require('uuid-v4');


class EventDb {

  /**
   * 
   * @param {*} config contains 
   * { 
   *     influxip : "some-ip",
   *     database : "some-database"
   * }
   */
  constructor(config) {
    this.config = config;

  }

  /**
   * @method save 
   * @param {*} event contains { name : "event-name", payload : "event-payload" }
   * @description saves an event in the 'measurement' "eventos", 
   * with tag "event" containing 'event.name', tag "instance_id" containing
   * the instance id created, and the value "payload" containing 'event.payload'
   * @returns a Promisse with the instance id created, if success; 
   * or an error, if failed
   * @example
        var EventDb = require("../EventDb");

        var config = 
        { 
            influxip : "localhost",
            database : "test000"
        } 

        eventDb = new EventDb(config);

        var evento = 
        {
            name : "almoco",
            payload : 
            { 
                prato : "churrasco", 
                preco : 38.80, 
            }
        }

        var promise = eventDb.save(evento);

        promise
        .then((instance_id) => { 
            console.log("instance id = ", instance_id);
        })
        .catch((e) => {
            console.log("error = ",e)
        });  
   */
  save(event) {
    var url = "http://" + this.config.influxip + ":8086/write";

    var req = unirest("POST", url);

    req.query({
      "db": this.config.database,
      "precision": "ms"
    });

    var instance_id = uuid();

    var payload = JSON.stringify(event.payload).replace(/"/g, '\\"');

    var line = "eventos" + "," + 
    "event=" + event.name + "," + 
    //"instance_id=" + "\"" + instance_id + "\"" + 
    "instance_id=" + instance_id +  
    " payload=" + "\"" + payload + "\""; 
    //" payload=" + payload; 

    console.log("line = ", line);

    req.send(line);
    
    return new Promise((resolve, reject) => {
      req.end(function (res) {
        if (res.error) {
          reject(res.error);
        }
        else {
          resolve(instance_id);
        }
      });
    });
  }

 

  /**
   * @method findByInterval
   * @param start initial timestamp
   * @param end final initial timestamp
   * @description retrieves a set of events between 'start' and 'finish', 
   * both excluding   
   * Cada elemento do conjunto contém:
   * {
   *    timestamp : timestamp-of-the-event
   *    name : event-name
   *    instance_id : instance-that-generated-the-event
   *    payload : JSON specific to the event
   * }
   * 
   * @example
        var EventDb = require("../EventDb");

        var config = 
        { 
            influxip : "localhost",
            database : "test002"
        } 

        eventDb = new EventDb(config);


        var start = 1517505841045000000;
        var end =   1517505958417000000;

        var  promise = eventDb.findByInterval(start, end);

        promise
        .then((events) => { 
            var size = events.length; 
            for (var i = 0;  i < size; i++) {
                var event = events[i];
                var payload = event.payload;
                console.log("evento", i, ":"
                            , "ts =", event.timestamp
                            , ", name =", event.name
                            , ", instanceId =", event.instanceId
                            , ", payload.prato =", payload.prato
                            , ", payload.preco =", payload.preco);
            }  
        })
        .catch((e) => {
            console.log("error = ",e)
        });
   */
  findByInterval(start, end) {

    var url = "http://" + this.config.influxip + ":8086/query";

    var req = unirest("POST", url);
    
    req.query({
      "pretty": "true"
    });
    
    req.headers({
      "Content-Type": "application/x-www-form-urlencoded"
    });

    var select = "SELECT * from \"eventos\" WHERE \"time\" > " + start + " and \"time\" < " + end;

    req.form({
      "db": this.config.database,
      "q": select
    });
    
    var self = this;

    return new Promise((resolve, reject) => {
      req.end(function (res) {
        if (res.error) {
          reject(res.error);
        }
        else {
          var values = res.body.results[0].series[0].values;
          //var events = self.parseValues(values);
          self.parseValues(values).then((events) => {
            resolve(events);
          });
        }
      });
    });        
  }

  /**
   * 
   * @param {*} values 
   * @description parsers a structure like:
    [
        [
            "2018-02-01T14:09:10.68Z",
            "almoco",
            "\"5705553d-dd2b-41f6-9dd8-5b83ed132964\"",
            "{\"prato\":\"salada\",\"preco\":19.5}"
        ],
        [
            "2018-02-01T14:09:59.127Z",
            "lanche",
            "\"44910883-9363-4522-bc66-8a875df7f783\"",
            "{\"prato\":\"joelho\",\"preco\":7.3}"
        ]
    ]   

    into:

    [ { timestamp: '2018-02-01T17:24:56.487Z',
        name: 'almoco',
        instanceId: 'ab77b461-109d-4086-b2ff-d1c3817e9e09',
        payload: { prato: 'salada', preco: 13.87 } },
      { timestamp: '2018-02-01T17:25:23.403Z',
        name: 'lanche',
        instanceId: '4172adfb-7d64-46e3-a18c-7b131e837065',
        payload: { prato: 'sanduiche', preco: 9.5 } } ]
    
  */
  parseValues(values) {
    return new Promise((resolve, reject) => { 
      var size = values.length;          
      var events = [];

      for (var i = 0;  i < size; i++) {
        var event = 
        {
          timestamp : values[i][0],
          name : values[i][1],
          instanceId : values[i][2],
          payload : JSON.parse(values[i][3])
        }
        events[i] = event;
      }
      resolve(events);
    });
  } 


}
module.exports = EventDb;