(function (){
  'use strict';

angular.module('foxcommUiApp')
    .factory('foxcommPromotionRule', function (Restangular) {

      var endpoint = Restangular.withConfig(function(RestangularConfigurer) {
        RestangularConfigurer.setBaseUrl('/foxcomm/social_analytics/');
      });

      endpoint.getByActionName = function(actionName, promise) {
        endpoint.one('admin/promotion_rules/', actionName).get().then(promise, function() {
          //fail function
          console.log("The request for promotion actions failed.");
        })
      };


      endpoint.getRuleOptions = function(promise) {
        endpoint.one('admin/promotion_rules/').get().then(promise, function() {
          //fail function
          console.log("The request for promotion actions failed.");
        })
      };

      endpoint.create = function(activityName, promise) {
        endpoint.service('admin/promotion_rules/').post({ActivityName: activityName}).then(promise, function() {
          //fail function
          console.log("The request for promotion actions failed.");
        })
      };

      return endpoint;
    });
})();
