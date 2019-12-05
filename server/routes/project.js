const express = require('express');
const project = require('../TrafficAccident/Repositories/ProjectRepo');


const router = express.Router();

// create new project
router.post('/', async (req, res, next) => {
  try {
    await project.create(req.body.name);
    res.status(200).json({'message': 'success'});
  } catch {
    res.status(400).json({'message': 'watting'});
  }
});

// get ply
router.get('/:name', (req, res) => {
  try {
    res.status(200).sendFile('../project/' + name + '/model.ply');
  } catch {
    res.status(400).json({'message': 'error'});
  }
});

router.delete('/:name', async (req, res) => {
  try {
    await project.remove(req.params.name);
    res.status(200).json({'message': 'success'});
  } catch {
    res.status(400).json({'message': 'error'});
  }
});

module.exports = router;
