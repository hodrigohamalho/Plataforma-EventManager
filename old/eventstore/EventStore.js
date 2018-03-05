var unirest = require("unirest");
const uuid = require('uuid-v4');


class EventStore {

  /**
   * @description constructor
   * @param {*} config contains
   */
  constructor(config = { influxip : "localhost", database : "event_manager"} ) {
    this.config = config;
    this.fieldsToSelect = " \"time\", \"instance_id\", \"name\", \"payload\", \"user_id\", \"user_name\"";
    this.createDb()
        .catch((e) => {
            console.log("error save = '",e,"'");
            throw e;
        });
  }


  /**
   * @method save
   * @param {*} event contains
   * @description saves an event in the 'measurement' "events",
   * with tag "name" containing 'event.name', tag "instance_id" containing
   * the instance id created, and the value "payload" containing 'event.payload'
   */
  save(event) {
    return new Promise((resolve, reject) => {
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
        var line = "events" + "," +
        "name=" + event.name + "," +
        "instance_id=" + instance_id + ","  +
        "user_id=" + user_id + "," +
        "user_name=" + user_name +
        " payload=" + "\"" + payload + "\"";


        req.send(line);


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
   * @description Retrieves a set of instances that have a `user.name`
   * that matches exactly the one provided
   * @param {*} name
   * @returns a Promise
   */
  findByUserName(name) {
    var select = "SELECT " + this.fieldsToSelect + " from \"events\" WHERE \"user_name\" = '" + name + "'";
    return this.find(select);
  }

  /**
   * @description Retrieves a set of instances that have a `user.id`
   * that matches exactly the one provided
   * @param {*} id
   */
  findByUserId(id) {
    var select = "SELECT " + this.fieldsToSelect + " from \"events\" WHERE \"user_id\" = '" + id + "'";
    return this.find(select);
  }

  /**
   * @method findByInterval
   * @param start initial timestamp, in miliseconds
   * @param end final initial timestamp, in miliseconds
   * @description retrieves a set of events between 'start' and 'finish',
   * both excluding
   */
  findByInterval(start, end = new Date().valueOf()) {

    //console.log("start =", start, ", end = ", end);

    var newEnd = end * 1000000;
    var newStart = start * 1000000;

    var select = "SELECT " + this.fieldsToSelect + " from \"events\" WHERE \"time\" > " + newStart + " AND \"time\" < " + newEnd;
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
   */
  findByEventInterval(name, start, end = new Date().valueOf()) {

    var newEnd = end * 1000000;
    var newStart = start * 1000000;
    var select = "SELECT " + this.fieldsToSelect + " from \"events\" WHERE \"time\" > " + newStart + " AND \"time\" < " + newEnd
                  + " AND \"name\" = " + "'" + name + "'";

    return this.find(select);
  }

  /**
   * @method find
   * @param {*} select is the SELECT SQL-like influxdb command to be executed
   * @description commands used by other functions are grouped in this function
   */
  find(select) {
    var url = "http://" + this.config.influxip + ":8086/query";

    var req = unirest("POST", url);

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
              payload : JSON.parse(values[i][3]),
              user : {
                  name : values[i][4],
                  id : values[i][5]
              }
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
    var req = unirest("POST", url);
    req.headers({
        "Cache-Control": "no-cache",
        "Content-Type": "application/x-www-form-urlencoded"
    });

    var cmd = "CREATE DATABASE " + this.config.database;

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