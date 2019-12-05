const fs      = require('fs');
const Promise = require('bluebird');
const rimraf  = require('rimraf');

const ProjectServ = require('../Services/ProjectServ');

const projectFolder = '../project/';

const create = async (name) => {
  return new Promise((resolve, reject) => {
    try {
      fs.mkdirSync(projectFolder + name + '/images', { recursive: true });
    } catch {
      console.log("folder exist");
    }

    if (ProjectServ.isRun()) {
      reject();
    } else {
      ProjectServ.run(name);
      resolve();
    }
  });
};

const remove = (name) => {
  return new Promise((resolve, reject) => {
    rimraf(projectFolder + name, undefined, resolve);
  });
};

module.exports = {
  create: create,
  remove: remove,
};