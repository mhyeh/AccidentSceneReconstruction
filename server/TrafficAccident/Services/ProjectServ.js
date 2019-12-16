const Promise       = require('bluebird');
const child_process = require('child-process-promise');
const kue           = require('kue');

let running = false;

const run = async name => {
  running = true;
  const cp = child_process.exec(`bash ../Reconstruct.bash ${name}`);
  cp.childProcess.stderr.on('data', data => {
    if (data.trim() == "Done!") {
      // socket io
    } else if (data.trim() == "trackloss") {
      // socket io
    }
  });
  try {
    await cp;
  }
  catch (e) {
    console.log(e);
  }
  console.log("finish!");
  running = false;
};

const isRun = () => running;

module.exports = {
  run:   run,
  isRun: isRun,
};