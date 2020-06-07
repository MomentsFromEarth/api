"use strict";

const AWS = require('aws-sdk');
const moment = require('moment');
const shortid = require('shortid');

const http = require('../services/utils');

function awsCredentials (auth) {
  let creds = { accessKeyId: auth.accessKeyId, secretAccessKey: auth.secretAccessKey };
  if (auth.sessionToken) creds.sessionToken = auth.sessionToken;
  return creds;
}

function genTimestamp () {
  //http://momentjs.com/docs/#/displaying/unix-timestamp-milliseconds/
  return moment().valueOf();
}

function genMomentId () {
  // https://github.com/dylang/shortid
  return shortid.generate();
}

function sqsClient (opts) {
  if (!opts) opts = {};
  opts.apiVersion = '2012-11-05';
  opts.region = 'us-east-1';
  //http://docs.aws.amazon.com/AWSJavaScriptSDK/latest/AWS/SQS.html
  return new AWS.SQS(opts);
}

function dynamodbDocumentClient (opts) {
  if (!opts) opts = {};
  opts.apiVersion = '2012-08-10';
  opts.region = 'us-east-1';
  //http://docs.aws.amazon.com/AWSJavaScriptSDK/latest/AWS/DynamoDB/DocumentClient.html
  return new AWS.DynamoDB.DocumentClient(opts);
}

function putItem (params, awsCreds) {
  return new Promise(function (resolve, reject) {
    //http://docs.aws.amazon.com/AWSJavaScriptSDK/latest/AWS/DynamoDB/DocumentClient.html#put-property
    dynamodbDocumentClient(awsCreds).put(params, function (err, res) {
      if (err) return reject(err);
      resolve(res);
    });
  });
}

function createMoment (moment, awsCreds) {
  return new Promise(function (resolve, reject) {
    putItem(putMomentParams(moment), awsCreds).then(function (res) {
      resolve(moment);
    }).catch(reject);
  });
}

function createNewMoment (moment, awsCreds) {
  return new Promise((resolve, reject) => {
    moment.moment_id = genMomentId();
    createMoment(moment, awsCreds).then(resolve).catch(err => {
      if (err.code === 'some error code for duplicate moment_id') {
        createNewMoment(moment, awsCreds).then(resolve).catch(reject);
      } else {
        reject(err);
      }
    });
  });
}

function putMomentParams (moment) {
  return {
    TableName: 'Moments',
    Item: {
      moment_id: moment.moment_id,
      title: moment.title,
      description: moment.description,
      original_filename: moment.original_filename,
      type: moment.type,
      size: moment.size,
      queue_id: moment.queue_id,
      status: moment.status,
      creator: moment.creator,
      created: moment.created,
      updated: moment.updated,
      host_id: moment.host_id || null
    }
  };
}

function formatNewMoment (data) {
  let now = genTimestamp();
  return {
    title: data.title,
    description: data.description,
    original_filename: data.original_filename,
    type: data.type,
    size: data.size,
    queue_id: data.queue_id,
    status: 'queued',
    creator: data.creator,
    created: now,
    updated: now,
    host_id: null
  }
}

function queueParams (moment) {
  return {
    MessageBody: JSON.stringify(moment),
    QueueUrl: 'https://sqs.us-east-1.amazonaws.com/776913033148/moments.fifo',
    MessageGroupId: 'mfe-api',
    MessageDeduplicationId: moment.queue_id
  };
}

function addMomentToProcessingQueue (moment, awsCreds) {
  return new Promise((resolve, reject) => {
    sqsClient(awsCreds).sendMessage(queueParams(moment), (err, res) => {
      if (err) return reject(err);
      resolve(res.MessageId);
    });
  });
}

module.exports.create = (req, res) => {
  let awsCreds = awsCredentials(req.body.auth);
  let newMoment = formatNewMoment(req.body.data);
  createNewMoment(newMoment, awsCreds).then(moment => {
    console.log('/moments', 'created', moment);

    addMomentToProcessingQueue(moment, awsCreds).then((msg) => {
      console.log('/moments', 'added to queue for processing', msg);
      res.status(http.codes.created).send(moment);
    }).catch(err => {
      console.log('/moments', 'error adding moment to processing queue', err);
      res.status(http.codes.internal_server_error).send({ message: 'internal_server_error' });
    });

  }).catch(err => {
    if (err.code === 'AccessDeniedException') {
      console.log('/moments', 'unauthorized', err);
      res.status(http.codes.unauthorized).send({ message: 'unauthorizied' });
    } else {
      console.log('/moments', 'internal_server_error', err);
      res.status(http.codes.internal_server_error).send({ message: 'internal_server_error' });
    }
  });
};

module.exports.update = (req, res) => {
  let awsCreds = awsCredentials(req.body.auth);
  let momentToUpdate = req.body.data;
  let path = '/moments/' + momentToUpdate.moment_id;
  momentToUpdate.updated = genTimestamp();
  putItem(putMomentParams(momentToUpdate), awsCreds).then(moment => {
    console.log(path, moment);
    res.status(http.codes.ok).send(moment);
  }).catch(err => {
    if (err.code === 'AccessDeniedException') {
      console.log(path, 'unauthorized', err);
      res.status(http.codes.unauthorized).send({ message: 'unauthorizied' });
    } else {
      console.log(path, 'internal_server_error', err);
      res.status(http.codes.internal_server_error).send({ message: 'internal_server_error' });
    }
  });
};

module.exports.show = (req, res) => {
  let moment = require('../../moments.json').filter(m => { 
    return m.id === req.params.id; 
  })[0];
  res.json(moment);
};
