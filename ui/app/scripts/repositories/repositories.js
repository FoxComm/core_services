(function(){
  'use strict';

  var serviceId = 'repositories';
  angular.module('foxcommUiApp')
  .factory(serviceId, ['$injector', repositories]);

  function repositories($injector){
    var manager;
    var service = {
      getRepo: getRepo,
      init: init
    };
    return service;

    function init(mgr) { manager = mgr;}

    function getRepo(repoName) {
      var fullRepoName = 'repository.' + repoName;
      var Repo = $injector.get(fullRepoName);
      return new Repo(manager);
    }
  }
})();