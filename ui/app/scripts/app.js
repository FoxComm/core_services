(function(){
  'use strict';

  /**
   * @ngdoc overview
   * @name foxcommUiApp
   * @description
   * # foxcommUiApp
   *
   * Main module of the application.
   */
  var app = angular
    .module('foxcommUiApp', [
      'ngAnimate',
      'ngCookies',
      'ngResource',
      'ngSanitize',
      'ngTouch',
      'ui.router',
      'restangular',
      'ui.bootstrap',
      'checklist-model',
    ]);

    app.config(function($stateProvider, $urlRouterProvider, RestangularProvider) {
      $urlRouterProvider.otherwise('/');
      RestangularProvider.setBaseUrl('/foxcomm');
    });

    app.run(["Auth", function(Auth) {
      Auth.setRoutingHandlers();
    }]);
})();
