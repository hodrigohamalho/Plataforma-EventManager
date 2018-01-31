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
   * 
   * @param {*} event contains { name : "event-name", payload : "event-payload" }
   */
  save(event, okcb, errorcb) {
    var url = "http://" + this.config.influxip + ":8086/write";

    var req = unirest("POST", url);

    req.query({
      "db": this.config.database
    });

    var instance_id = uuid();

    var payload = JSON.stringify(event.payload).replace(/"/g, '\\"');

    var line = event.name + " payload=" + "\"" + payload + "\"" + ",instance_id=" + "\"" + instance_id + "\""; 

    //console.log("line = ", line);

    req.send(line);
    
    req.end(function (res) {
      if (res.error) {
        errorcb(res.error);
      }
      else {
        okcb(instance_id);
      }
    });

  }
}

module.exports = EventDb;