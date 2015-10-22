(function(){
  'use strict';

  angular.module('foxcommUiApp')
      .factory('foxcommReports', function (Restangular) {


        var endpoint = Restangular.withConfig(function(RestangularConfigurer) {
          RestangularConfigurer.setBaseUrl('/foxcomm/social_analytics/');
        });


        endpoint.getAllReports = function(promise) {
          endpoint.one('admin/reports').get().then(promise, function() {
            //fail function
            console.log("The request for reports failed.");
          })
        };

        return endpoint;
      });
})();
