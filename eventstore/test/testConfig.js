

class Config {


    get() {
        var data = 
        {
            influxip : "localhost",
            database : "test003"
        }
        return data;
    }
}

module.exports = Config;