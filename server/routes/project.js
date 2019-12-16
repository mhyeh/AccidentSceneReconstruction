const express = require('express');
const path = require('path');
const child_process = require('child-process-promise');

const projectRepo = require('../TrafficAccident/Repositories/ProjectRepo');

const router = express.Router();

// create new project
router.post('/', async (req, res) => {
  try {
    await projectRepo.create(req);
    res.status(200).end();
  } catch {
    res.status(400).end();
  }
});

// get ply
router.get('/:name', (req, res) => {
  // try {
  //   const filePath = path.resolve(`../project/${req.params.name}/models/model.ply`);
  //   res.status(200).sendFile(filePath);
  // } catch (e) {
  //   console.log(e);
  //   res.status(400).end();
  // }
  const cp = child_process.exec(`../viewer/src/src ${req.params.name}`);
  cp.childProcess.stdout.on('data', data => {
    console.log(data)
  });
  cp.childProcess.stderr.on('data', data => {
    console.log(data)
  });
  cp;
  res.status(200).end();
});

router.delete('/:name', async (req, res) => {
  try {
    await projectRepo.remove(req.params.name);
    res.status(200).end();
  } catch {
    res.status(400).end();
  }
});

module.exports = router;
