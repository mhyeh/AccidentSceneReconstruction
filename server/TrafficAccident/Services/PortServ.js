let ports = [9000, 9001, 9002, 9003, 9004, 9005, 9006, 9007];

const getPort = () => ports.pop();

const releasePort = port => {
  ports.push(port);
};

module.exports = {
    getPort:     getPort,
    releasePort: releasePort,
};