var config = {};

const executorHost = process.env.EXECUTOR_HOST || "localhost";
const routerHost = process.env.ROUTER_HOST || "localhost";

config.PORT = 8081;
config.executorUrl = "http://" + executorHost + ":8085/executor";
config.proxyPresentationUrl = "http://" + routerHost + ":8086/event";
//config.executorUrl = "http://localhost:8081/testexecutor";

module.exports = config;