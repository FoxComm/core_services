(function() {
'use strict';

  /**
   * @ngdoc function
   * @name foxcommUiApp.controller:MainCtrl
   * @description
   * # MainCtrl
   * Controller of the foxcommUiApp
   */

  angular.module('foxcommUiApp')
    .config(['$stateProvider', MainConfig])
    .controller('MainCtrl', ['$state', '$window', 'Auth', MainCtrl]);

  function MainConfig($stateProvider) {
    $stateProvider.state('admin', {
      url: '/',
      templateUrl: 'admin/admin.html', 
      controller: 'MainCtrl as vm'
    });
  }

  function MainCtrl($state, $window, Auth) {
    var vm = this;
    vm.busyMessage = 'We are loading the page.';
    vm.logout = logout;

    function logout() {
      Auth.logout(success, error);

      function success() {
        $window.location.pathname = $window.location.pathname + "login.html";
      }

      function error(data) {
        $window.alert("failed to logout");
      }
    }
  }
})();
