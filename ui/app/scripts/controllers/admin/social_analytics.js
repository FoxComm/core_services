(function() {
  'use strict';

  angular.module('foxcommUiApp')
    .config(['$stateProvider', SocialAnalyticsConfig])
    .controller('SocialAnalyticsCtrl', [SocialAnalyticsCtrl]);

  function SocialAnalyticsConfig($stateProvider) {
    $stateProvider.state('admin.social_analytics', {
      url: 'social_analytics',
      templateUrl: 'admin/social_analytics/index.html',
      controller: 'SocialAnalyticsCtrl as vm',
      data: {
        title: 'Social Analytics'
      }
    });
  }

  function SocialAnalyticsCtrl() {
    var vm = this;

    vm.donkey = 'hello smonken';
  }
})();
