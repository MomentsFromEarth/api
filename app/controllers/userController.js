"use strict";

const http = require('../services/utils');

module.exports.show = (req, res) => {
  let user = require('../../users.json').filter(u => { 
    return u.id === req.params.id; 
  })[0];
  res.json(user);
};
