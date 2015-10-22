(function(){
  'use strict';

  angular.module('foxcommUiApp')
  .factory('Auth', ['$window', '$rootScope', '$cookieStore', 'Restangular', function ($window, $rootScope, $cookieStore, Restangular) {
    var currentUser = $cookieStore.get('user') || undefined,
    service = Restangular.withConfig(restangularConfig);

    service.currentUser = currentUser;
    service.login = login;
    service.logout = logout;
    service.setCurrentUser = setCurrentUser;
    service.setRoutingHandlers = setRoutingHandlers;
    return service

    function restangularConfig(RestangularConfigurer) {
      RestangularConfigurer.setBaseUrl('/session');
    }

    function login(user, success, error){
      service.one('login').customPUT(user)
      .then(function(data) {
        service.setCurrentUser(data.plain());
        success(data);
      }, error);
    }

    function logout(success, error){
      service.one('logout').get()
      .then(function(data) {
        service.setCurrentUser(null);
        success(data);
      }, error);
    }

    function setCurrentUser(user) {
      if(!user) {
        $cookieStore.remove('user');
      } else {
        $cookieStore.put('user',user);
      }
      service.currentUser = user;
    }

    function setRoutingHandlers() {
      $rootScope.$on('$stateChangeStart',
                     function (event, toState, toParams, fromState, fromParams) {
                       if ((service.currentUser == null || service.currentUser.Role != "admin") &&
                           !$window.location.pathname.match(/login/)) {
                         $window.location.pathname = $window.location.pathname + "login.html";
                       }
                     }
                    );
    }
  }]);
})();

