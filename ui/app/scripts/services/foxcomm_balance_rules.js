(function(){
  'use strict';

angular.module('foxcommUiApp')
  .factory('foxcommBalanceRule', function (Restangular) {

    var endpoint = Restangular.withConfig(function(RestangularConfigurer) {
      RestangularConfigurer.setBaseUrl('/foxcomm/social_analytics/');
    });

    endpoint.getByCurrencyName = function(currencyName, promise) {
      endpoint.one('admin/balance_rules/', currencyName).get().then(promise, function() {
        //fail function
        console.log("The request for balance actions failed.");
      })
    };

    endpoint.create = function(currencyName, promise) {
      endpoint.service('admin/balance_rules/').post({CurrencyName: currencyName}).then(promise, function() {
        //fail function
        console.log("The request for balance actions failed.");
      })
    };

    return endpoint;
  });
})();
