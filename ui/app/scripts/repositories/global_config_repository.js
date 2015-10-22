(function(){
  'use strict';

  var serviceId = 'repository.globalConfig';
  angular.module('foxcommUiApp')
    .factory(serviceId, ['Restangular', function(Restangular) {
      var route = 'core',
          repository = Restangular.service(route);

      repository.build = build;
      return repository;

      function build(object, route) {
        return Restangular.restangularizeElement('', object, route);
      }

    }]);
})();
