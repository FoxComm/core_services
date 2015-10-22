(function() {
  'use strict';

angular.module('foxcommUiApp')
  .factory('foxcommAttributionRule', function (Restangular) {

    var endpoint = Restangular.withConfig(function(RestangularConfigurer) {
      RestangularConfigurer.setBaseUrl('/foxcomm/social_analytics/');
    });

    endpoint.getAttributionActions = function(promise) {
      endpoint.one('admin/attribution_rules').get().then(promise, function() {
        //fail function
        console.log("The request for attribution actions failed.");
      });
    };

    endpoint.getByActionName = function(actionName, promise) {
      endpoint.one('admin/attribution_rules/', actionName).get().then(promise, function() {
        //fail function
        console.log("The request for attribution actions failed.");
      })
    };

    endpoint.create = function(activityName, promise) {

      endpoint.service('admin/attribution_rules/').post({ActivityName: activityName}).then(promise, function() {
        //fail function
        console.log("The request for attribution actions failed.");
      })
    };

    return endpoint;
  });
})();
