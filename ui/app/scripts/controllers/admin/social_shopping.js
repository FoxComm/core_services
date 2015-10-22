(function() {
  'use strict';

  angular.module('foxcommUiApp')
    .config(['$stateProvider', SocialShoppingConfig])
    .controller('SocialShoppingCtrl', ['$scope', 'foxcommSocialShopping', 'repository.globalConfig', 'logger', SocialShoppingCtrl]);

  function SocialShoppingConfig($stateProvider) {
    $stateProvider.state('admin.social_shopping', {
      url: 'social_shopping',
      templateUrl: 'admin/social_shopping/index.html',
      controller: 'SocialShoppingCtrl',
      data: {
        title: 'Social Shopping'
      }
    });
  }

  function SocialShoppingCtrl($scope, foxcommSocialShopping, GlobalConfig, logger) {
    $scope.socialShoppingEnable = true;
    $scope.fcSocialShopping = foxcommSocialShopping.one('admin');
    $scope.fcSocialShopping.get().then(function(data) {
      $scope.socialShoppingPrefs = data;
      console.log($scope.socialShoppingPrefs);
    });

    $scope.updatePreference = function(preference) {
      if (preference.put === undefined) {
        GlobalConfig.build(preference, 'social_shopping/admin')
      }
      preference.put().then(function(data) {
        logger.logSuccess("Successfully updated", null, null, true);
      }, function(error) {
        logger.logError("An error ocurred. Please try again later.", null, null, true);
      });
    }
  }
})();
