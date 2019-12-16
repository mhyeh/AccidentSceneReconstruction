const fs      = require('fs');
const Promise = require('bluebird');
const rimraf  = require('rimraf');
let formidable = require('formidable')

const ProjectServ = require('../Services/ProjectServ');

const projectFolder = '../project';

const processFormData = data => new Promise((resolve, reject) => {
  let form = new formidable.IncomingForm();
  form.encoding       = 'utf-8';
  form.keepExtensions = true;
  form.multiples      = true;

  form.parse(data, (err, fields, files) => {
    if (err) {
      console.log(err);
      reject('file error');
      return;
    }

    resolve({name: fields['name'], file: files['calibration'].path});
  });
});


const create = async req => new Promise(async (resolve, reject) => {
  let data 
  try {
    data = await processFormData(req);
  } catch {
    reject('file error');
    return;
  }

  try {
    fs.mkdirSync(`${projectFolder}/${data.name}/images`, { recursive: true });
  } catch {
    console.log("folder exist");
  }

  fs.rename(data.file, `${projectFolder}/${data.name}/camera_parameter.yml`, err => {
    if (err) {
      console.log(err);
      reject('file error');
      return;
    }
    
    if (ProjectServ.isRun()) {
      reject();
    } else {
      ProjectServ.run(data.name);
      resolve();
    }
  });
});

const remove = name => new Promise((resolve, reject) => {
  rimraf(`${projectFolder}/${name}`, undefined, resolve);
});

module.exports = {
  create: create,
  remove: remove,
};