

class Config {


    get() {
        var data = 
        {
            influxip : "localhost",
            database : "test002"
        }
        return data;
    }
}

module.exports = Config;