"use strict";

const http = require('../services/utils');

module.exports.moments = (req, res) => {
  let moments = require('../../moments.json');
  res.json(moments);
};

module.exports.users = (req, res) => {
  let users = require('../../users.json');
  res.json(users);
};
