(function () {
  'use strict';

  var controllerId = 'auth.login';

  angular.module('foxcommUiApp')
      .config(['$stateProvider', LoginConfig])
      .controller(controllerId, [
        '$window',
        '$scope',
        'Auth',
        LoginCtrl]);

  function LoginConfig($stateProvider) {
    $stateProvider.state('auth.login', {
      url: '/login.html',
      templateUrl: 'auth/login.html',
      controller: 'auth.login as vm',
      data: {
        title: 'Login'
      }
    });
  }

  function LoginCtrl($window, $scope, Auth) {
    var vm = this;
    vm.title = 'Sign In';
    vm.user = undefined;
    vm.login = login;

    function login() {
      if(!vm.user) return;

      Auth.login(vm.user, success, error);

      function success(data) {
        var paths = $window.location.pathname.split("/");
        paths.pop();
        $window.location.pathname = paths.join("/") + "/";
      }

      function error(data) {
        $window.alert("failed to login");
      }
    }
  }
})();
