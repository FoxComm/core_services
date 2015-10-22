(function(){
  'use strict';

angular.module('foxcommUiApp')
  .factory('foxcommBackups', function (Restangular) {

    var endpoint = Restangular.withConfig(function(RestangularConfigurer) {
      RestangularConfigurer.setBaseUrl('/foxcomm/backups/');
      RestangularConfigurer.setRestangularFields({id: 'Id'});
    });

    endpoint.build = function (object) {
      return Restangular.restangularizeElement('', object, 'core/stores');
    };

  return endpoint
  });

})();
