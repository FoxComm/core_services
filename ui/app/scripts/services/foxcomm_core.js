(function(){
  'use strict';

angular.module('foxcommUiApp')
  .factory('foxcommCore', function (Restangular) {

    var endpoint = Restangular.withConfig(function(RestangularConfigurer) {
      RestangularConfigurer.setBaseUrl('/foxcomm/core/');
      RestangularConfigurer.setRestangularFields({id: 'Id'});
    });

    endpoint.build = function (object) {
      return Restangular.restangularizeElement('', object, 'core/stores');
    };

  return endpoint
  });

})();
