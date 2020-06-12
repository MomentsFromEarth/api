"use strict";

module.exports = (app) => {

  // account
  let accountController = require("../controllers/accountController");
  app.post("/account/create", accountController.create);
  app.post("/account/profile", accountController.profile);

  // moment
  let momentController = require("../controllers/momentController");
  app.post('/moments', momentController.create);
  app.put('/moments/:id', momentController.update);
  app.get("/moments/:id", momentController.show);

  // user
  let userController = require("../controllers/userController");
  app.get("/users/:id", userController.show);

  // search
  let searchController = require("../controllers/searchController");
  app.get("/search/moments", searchController.moments);
  app.get("/search/users", searchController.users);

  // root
  let appController = require("../controllers/appController");
  app.get("/", appController.status);

};
