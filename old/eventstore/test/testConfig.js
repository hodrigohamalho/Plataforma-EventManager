

class Config {


    get() {
        var data = 
        {
            influxip : "localhost",
            database : "test013"
        }
        return data;
    }
}

module.exports = Config;