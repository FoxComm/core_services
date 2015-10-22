'use strict';

/**
 * @ngdoc service
 * @name ngFoxcommSocialApp.foxcommSocial
 * @description
 * # foxcommSocial
 * Factory in the ngFoxcommSocialApp.
 */
angular.module('foxcommUiApp')
  .factory('foxcommSocialShopping', function (Restangular) {
    return Restangular.withConfig(function(RestangularConfigurer) {
      RestangularConfigurer.setBaseUrl('/foxcomm/social_shopping/');
    });
  });