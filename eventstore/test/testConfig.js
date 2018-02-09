

class Config {


    get() {
        var data = 
        {
            influxip : "localhost",
            database : "test005"
        }
        return data;
    }
}

module.exports = Config;