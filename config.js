var config = {};

const executorHost = process.env.EXECUTOR_HOST || "localhost";
const routerHost = process.env.ROUTER_HOST || "localhost";

config.PORT = 8081;
config.executorUrl = "http://" + executorHost + ":8000/executor";
config.proxyPresentationUrl = "http://" + routerHost + ":8086/event";

module.exports = config;
