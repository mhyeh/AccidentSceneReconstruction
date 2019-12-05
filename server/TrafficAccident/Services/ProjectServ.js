const Promise       = require('bluebird');
const child_process = require('child-process-promise');
const kue           = require('kue');

let running = false;

const run = async (name) => {
  running = true;
  await child_process.exec('../Reconstruct.sh ' + name);
  running = false;
};

const isRun = () => {
  return running;
}

module.exports = {
  run: run,
  isRun: isRun
};