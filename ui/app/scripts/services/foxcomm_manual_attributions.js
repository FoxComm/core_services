(function() {
  'use strict';

angular.module('foxcommUiApp')
  .factory('foxcommManualAttributions', function (Restangular) {

    var endpoint = Restangular.withConfig(function(RestangularConfigurer) {
      RestangularConfigurer.setBaseUrl('/foxcomm/social_analytics/admin');
      RestangularConfigurer.setRestangularFields({
        id: "Id"
      });
    });

    return endpoint;
  });
})();
