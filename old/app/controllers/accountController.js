"use strict";

const AWS = require('aws-sdk');
const moment = require('moment');
const shortid = require('shortid');

const http = require('../services/utils');

function log (msg, obj) {
  console.log('MFE API: ' + msg);
  if (obj) console.log(obj);
}

function dynamodbDocumentClient (opts) {
  if (!opts) opts = {};
  opts.apiVersion = '2012-08-10';
  opts.region = 'us-east-1';
  //http://docs.aws.amazon.com/AWSJavaScriptSDK/latest/AWS/DynamoDB/DocumentClient.html
  return new AWS.DynamoDB.DocumentClient(opts);
}

function fetchTemplateParams (id) {
  return {
    TableName: 'TaskRunner_Templates',
    KeyConditionExpression: 'template_id = :templateid',
    ExpressionAttributeValues: {
      ':templateid': id
    }
  };
}

function fetchUserParams (identityId) {
  return {
    TableName: 'Users',
    IndexName: 'identity_id-index',
    KeyConditionExpression: 'identity_id = :identityid',
    ExpressionAttributeValues: {
      ':identityid': identityId
    }
  };
}

function createUserParams (user) {
  return {
    TableName: 'Users',
    Item: {
      user_id: user.id,
      identity_id: user.identity_id,
      name: user.name,
      email: user.email,
      avatar: user.avatar,
      created: user.created,
      updated: user.updated
    }
  };
}

function saveTaskParams (task) {
  return {
    TableName: 'TaskRunner_Tasks',
    Item: {
      task_id: task.id,
      started: task.started,
      task_status: task.status,
      data: task
    }
  };
}

function saveTemplateParams (name, template) {
  return {
    TableName: 'TaskRunner_Templates',
    Item: {
      template_id: name,
      data: template
    }
  };
}

function fetchItem (awsCreds, queryParams) {
  return new Promise(function (resolve, reject) {
    //http://docs.aws.amazon.com/AWSJavaScriptSDK/latest/AWS/DynamoDB/DocumentClient.html#query-property
    dynamodbDocumentClient(awsCreds).query(queryParams, function (err, res) {
      if (err) return reject(err);

      var item = res && res.Items && res.Items[0];
      if (!item) return reject({ code: 'ResourceNotFoundException' });

      resolve(item);
    });
  });
}

function saveItem (params) {
  return new Promise(function (resolve, reject) {
    //http://docs.aws.amazon.com/AWSJavaScriptSDK/latest/AWS/DynamoDB/DocumentClient.html#put-property
    dynamodbDocumentClient().put(params, function (err, res) {
      if (err) return reject(err);
      resolve(res);
    });
  });
}

function createItem (params, awsCreds) {
  return new Promise(function (resolve, reject) {
    //http://docs.aws.amazon.com/AWSJavaScriptSDK/latest/AWS/DynamoDB/DocumentClient.html#put-property
    dynamodbDocumentClient(awsCreds).put(params, function (err, res) {
      if (err) return reject(err);
      resolve(res);
    });
  });
}

function saveTask (task) {
  return new Promise(function (resolve, reject) {
    saveItem(saveTaskParams(task)).then(function () {
      resolve();
    }).catch(reject);
  });
}

function saveTemplate (name, template) {
  return new Promise(function (resolve, reject) {
    saveItem(saveTemplateParams(name, template)).then(function () {
      resolve();
    }).catch(reject);
  });
}

function createUser (user, awsCreds) {
  return new Promise(function (resolve, reject) {
    createItem(createUserParams(user), awsCreds).then(function (res) {
      resolve(user);
    }).catch(reject);
  });
}

function scanTemplateParams () {
  return {
    TableName: 'TaskRunner_Templates',
    Limit: 100
  };
}

function updateTemplateParams (templateId, nextRun, lastRun) {
  var ue = ['cron_next_run = :cronnextrun'];
  var eav = { ':cronnextrun' : nextRun };
  if (lastRun) {
    ue.push('cron_last_run = :cronlastrun');
    eav[':cronlastrun'] = lastRun;
  }

  return {
    TableName: 'TaskRunner_Templates',
    Key: { template_id : templateId },
    UpdateExpression: 'set ' + ue.join(', '),
    ExpressionAttributeValues: eav
  };
}

function scanItems (params) {
  return new Promise(function (resolve, reject) {
    //http://docs.aws.amazon.com/AWSJavaScriptSDK/latest/AWS/DynamoDB/DocumentClient.html#scan-property
    dynamodbDocumentClient().scan(params, function (err, res) {
      if (err) return reject(err);
      resolve(res);
    });
  });
}

function scanTemplates () {
  return new Promise(function (resolve, reject) {
    scanItems(scanTemplateParams()).then(function (templates) {
      resolve(templates);
    }).catch(reject);
  });
}

function updateItem (params) {
  return new Promise(function (resolve, reject) {
    //http://docs.aws.amazon.com/AWSJavaScriptSDK/latest/AWS/DynamoDB/DocumentClient.html#update-property
    dynamodbDocumentClient().update(params, function (err, res) {
      if (err) return reject(err);
      resolve(res);
    });
  });
}

function updateTemplate (templateId, nextRun, lastRun) {
  return new Promise(function (resolve, reject) {
    var params = updateTemplateParams(templateId, nextRun, lastRun);
    updateItem(params).then(resolve).catch(reject);
  });
}

function genTaskId (id) {
  //ex. Test_20160810200315
  return id + '_' + moment.utc().format('YYYYMMDDHHmmss');
}

function fetchUser (identityId, awsCreds) {
  return new Promise ((resolve, reject) => {
    fetchItem(awsCreds, fetchUserParams(identityId))
      .then(resolve)
      .catch(reject);
  });
}

function awsCredentials (auth) {
  return {
    accessKeyId: auth.accessKeyId,
    secretAccessKey: auth.secretAccessKey,
    sessionToken: auth.sessionToken
  };
}

function genTimestamp () {
  //http://momentjs.com/docs/#/displaying/unix-timestamp-milliseconds/
  return moment().valueOf();
}

function genUserId () {
  // https://github.com/dylang/shortid
  return shortid.generate();
}

function createNewUser (user, awsCreds) {
  return new Promise((resolve, reject) => {
    user.id = genUserId();
    createUser(user, awsCreds).then(resolve).catch(err => {
      if (err.code === 'some error code for duplicate userId') {
        createNewUser(user, awsCreds).then(resolve).catch(reject);
      } else {
        reject(err);
      }
    });
  });
}

function formatNewUser (data, auth) {
  let now = genTimestamp();
  return {
    identity_id: auth.identityId,
    name: data.name,
    email: data.email,
    avatar: data.avatar,
    created: now,
    updated: now
  }
}

module.exports.create = (req, res) => {
  let identityId = req.body.auth.identityId;
  let awsCreds = awsCredentials(req.body.auth);
  let newUser = formatNewUser(req.body.data, req.body.auth);
  fetchUser(identityId, awsCreds).then(user => {
    console.log('/account/create', 'exists', user);
    res.status(http.codes.ok).send(user);
  }).catch(err => {
    if (err.code === 'AccessDeniedException') {
      console.log('/account/create', 'unauthorized', err);
      res.status(http.codes.unauthorized).send({ message: 'unauthorizied' });
    } else if (err.code === 'ResourceNotFoundException') {

      createNewUser(newUser, awsCreds).then(user => {
        console.log('/account/create', 'created', user);
        res.status(http.codes.created).send(user);
      }).catch(err => {
        if (err.code === 'AccessDeniedException') {
          console.log('/account/create', 'unauthorized', err);
          res.status(http.codes.unauthorized).send({ message: 'unauthorizied' });
        } else {
          console.log('/account/create', 'internal_server_error', err);
          res.status(http.codes.internal_server_error).send({ message: 'internal_server_error' });
        }
      });

    } else {
      console.log('/account/create', 'internal_server_error', err);
      res.status(http.codes.internal_server_error).send({ message: 'internal_server_error' });
    }
  });
};

module.exports.profile = (req, res) => {
  let identityId = req.body.data.identityId;
  let awsCreds = awsCredentials(req.body.auth);
  fetchUser(identityId, awsCreds).then(user => {

    console.log('/account/profile', 'found', user);
    res.status(http.codes.ok).send(user);

  }).catch(err => {
    if (err.code === 'AccessDeniedException') {
      console.log('/account/profile', 'unauthorizied', err);
      res.status(http.codes.unauthorized).send({ message: 'unauthorizied' });
    } else if (err.code === 'ResourceNotFoundException') {
      console.log('/account/profile', 'profile_does_not_exist', err);
      res.status(http.codes.not_found).send({ 
        message: 'profile_does_not_exist'
      });
    } else {
      console.log('/account/profile', 'internal_server_error', err);
      res.status(http.codes.internal_server_error).send({ message: 'internal_server_error' });
    }
  });
};
