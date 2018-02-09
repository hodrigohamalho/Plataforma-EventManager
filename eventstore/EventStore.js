var unirest = require("unirest");
const uuid = require('uuid-v4');


class EventStore {

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
    this.createDb()
        .catch((e) => {
            console.log("error save = '",e,"'");
            throw e;
        });
  }


  /**
   * @method save 
   * @param {*} event contains 
    event: {
       name : "event-name",
       payload : {},    
       user : {
           id : "user-id",
           name : "user-name"
       }
    }
   * @description saves an event in the 'measurement' "events", 
   * with tag "name" containing 'event.name', tag "instance_id" containing
   * the instance id created, and the value "payload" containing 'event.payload'
   * @returns a Promisse with the instance id created, if success; 
   * or an error, if failed
   * @example
        var EventStore = require("../EventStore");

        var config = 
        { 
            influxip : "localhost",
            database : "test003"
        } 

        eventStore = new EventStore(config);

        var evento = 
        {
            name : "almoco",
            payload : 
            { 
                prato : "churrasco", 
                preco : 38.80, 
            },
            user : 
            {
                name : "Joao da Silva",
                id : "RI - 12874"
            }
        }

        var promise = eventStore.save(evento);

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

    var user_name = this.format(event.user.name);
    var user_id = this.format(event.user.id);

/*     console.log("user name=", user_name);
    console.log("user id=", user_id);
 */

    var line = "events" + "," + 
    "name=" + event.name + "," + 
    "instance_id=" + instance_id + ","  +
    "user_id=" + user_id + "," +
    "user_name=" + user_name + 
    " payload=" + "\"" + payload + "\""; 

//    console.log("line =", line);
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
   * @param start initial timestamp, in miliseconds
   * @param end final initial timestamp, in miliseconds
   * @description retrieves a set of events between 'start' and 'finish', 
   * both excluding   
   * Each element in the set contains:
   * {
   *    timestamp : timestamp-of-the-event
   *    name : event-name
   *    instance_id : instance-that-generated-the-event
   *    payload : JSON specific to the event
   * }
   * 
   * @example
      var EventStore = require("../EventStore");

      var config = 
      { 
          influxip : "localhost",
          database : "test003"
      } 

      eventStore = new EventStore(config);

      // will find
      var start = 1517513761893;
      var end =   1517513864517;


      var  promise = eventStore.findByInterval(start, end);

      promise
      .then((events) => { 
          var size = events.length; 
          if (size == 0) {
              console.log("No event found");
          }
          else {
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
              console.log(events);
          }
      })
      .catch((e) => {
          console.log("error = ",e)
      });
   */
  findByInterval(start, end = new Date().valueOf()) {

    //console.log("start =", start, ", end = ", end);

    var newEnd = end * 1000000;
    var newStart = start * 1000000;

    var select = "SELECT * from \"eventos\" WHERE \"time\" > " + newStart + " AND \"time\" < " + newEnd;
    return this.find(select);
 
  }


  /**
   * @method findEventByInterval
   * @param name name of the event
   * @param start initial timestamp, in miliseconds
   * @param end final initial timestamp, in miliseconds
   * @description retrieves a set of events between 'start' and 'finish', 
   * both excluding   
   * Each element in the set contains:
   * {
   *    timestamp : timestamp-of-the-event
   *    name : event-name
   *    instance_id : instance-that-generated-the-event
   *    payload : JSON specific to the event
   * }
   * 
   * @example
      var EventStore = require("../EventStore");

      var config = 
      { 
          influxip : "localhost",
          database : "test003"
      } 

      eventStore = new EventStore(config);


      var start = 1517513761893;
      var end =   1517514095963;
      var name = "colacao";

      var  promise = eventStore.findByEventInterval(name, start, end);

      promise
      .then((events) => { 
          var size = events.length; 
          if (size == 0) {
              console.log("No event found");
          }
          else {
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
              
          }
      })
      .catch((e) => {
          console.log("error = ",e)
      });
   */  
  findByEventInterval(name, start, end = new Date().valueOf()) {

    var newEnd = end * 1000000;
    var newStart = start * 1000000;
    var select = "SELECT * from \"eventos\" WHERE \"time\" > " + newStart + " AND \"time\" < " + newEnd 
                  + " AND \"name\" = " + "'" + name + "'";

    return this.find(select);
  }

  /**
   * @method find
   * @param {*} select is the SELECT SQL-like influxdb command to be executed
   * @description commands used by other functions are grouped in this function
   */
  find(select) {
    //console.log(select);

    var url = "http://" + this.config.influxip + ":8086/query";

    var req = unirest("POST", url);
    
/*     req.query({
      "pretty": "true"
    });
 */    
    req.headers({
      "Content-Type": "application/x-www-form-urlencoded"
    });

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
          var results = res.body.results[0];
          self.parseValues(results).then((events) => {
            resolve(events);
          });
        
        }
      });
    });       
  }

  /**
   * 
   * @param {*} results
   * @description parsers a structure like:
    {
        "results": [
            {
                "statement_id": 0,
                "series": [
                    {
                        "name": "eventos",
                        "columns": [
                            "time",
                            "instance_id",
                            "name",
                            "payload"
                        ],
                        "values": [
                            [
                                "2018-02-01T19:36:44.926Z",
                                "989caaf1-4bcb-4a65-a142-fec1a9f1bf59",
                                "colacao",
                                "{\"prato\":\"maçã\",\"preco\":1.1}"
                            ],
                            [
                                "2018-02-01T19:37:08.395Z",
                                "9d0ecd35-a83e-4c66-8079-d0fcb188b65d",
                                "almoco",
                                "{\"prato\":\"salada\",\"preco\":18.9}"
                            ]
                        ]
                    }
                ]
            }
        ]
    }

    into:

    [ { timestamp: '2018-02-01T19:36:44.926Z',
        instanceId: '989caaf1-4bcb-4a65-a142-fec1a9f1bf59',
        name: 'colacao',
        payload: { prato: 'maçã', preco: 1.1 } },
      { timestamp: '2018-02-01T19:37:08.395Z',
        instanceId: '9d0ecd35-a83e-4c66-8079-d0fcb188b65d',
        name: 'almoco',
        payload: { prato: 'salada', preco: 18.9 } } ]
      
  */
  parseValues(results) {
    return new Promise((resolve, reject) => { 

      if (results.series) {
          var values = results.series[0].values;
          var size = values.length;          
          var events = [];

          for (var i = 0;  i < size; i++) {
            var event = 
            {
              timestamp : values[i][0],
              instanceId : values[i][1],              
              name : values[i][2],
              payload : JSON.parse(values[i][3])
            }
            events[i] = event;
          }
          resolve(events);
      }
      else {
        resolve([]);
      }
    });
  
  } 


  createDb() {
    var url = "http://" + this.config.influxip + ":8086/query";
    //console.log("url =",url);

    var req = unirest("POST", url);

    req.headers({
        "Cache-Control": "no-cache",
        "Content-Type": "application/x-www-form-urlencoded"
    });
    
    var cmd = "CREATE DATABASE " + this.config.database;
    //console.log("cmd =",cmd);

    req.form({
        "q": cmd
    });

    var self = this;
    return new Promise((resolve, reject) => {
        req.end(function (res) {
          if (res.error) {
            reject(res.error);
          }
          else {
            resolve("db " + self.config.database + " created");
          }
        });
      });    
  }

  
  format(str) {      
    var aux = str;
    aux = aux.replace(/ /g, "\\ ");
    aux = aux.replace(/=/g, "\\=");
    aux = aux.replace(/,/g, "\\,");
    return aux;
  }

}
module.exports = EventStore;