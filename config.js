var config = {};

const executorHost = process.EXECUTOR_HOST || "localhost";

config.PORT = 8081;
config.executorUrl = "http://" + executorHost + ":8085/executor";
config.proxyPresentationUrl = "http://localhost:8086/event";
//config.executorUrl = "http://localhost:8081/testexecutor";

module.exports = config;