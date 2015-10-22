(function(){
  'use strict';

angular.module('foxcommUiApp')
    .factory('foxcommAccumulationRule', function (Restangular) {

      var endpoint = Restangular.withConfig(function(RestangularConfigurer) {
        RestangularConfigurer.setBaseUrl('/foxcomm/social_analytics/');
      });

      endpoint.getAccumulationActions = function(promise) {
        endpoint.one('admin/accumulation_rules').get().then(promise, function() {
          //fail function
          console.log("The request for accumulation actions failed.");
        });
      };

      endpoint.getByActionName = function(actionName, promise) {
        endpoint.one('admin/accumulation_rules/', actionName).get().then(promise, function() {
          //fail function
          console.log("The request for accumulation actions failed.");
        })
      };

      endpoint.create = function(activityName, promise) {

        endpoint.service('admin/accumulation_rules/').post({ActivityName: activityName}).then(promise, function() {
          //fail function
          console.log("The request for accumulation actions failed.");
        })
      };

      return endpoint;
    });
})();
