(function() {
  'use strict';

  angular.module('foxcommUiApp')
    .config(['$stateProvider', RulesEngineConfig]);

  function RulesEngineConfig($stateProvider) {
    $stateProvider.state('admin.rules_engine', {
      url: 'rules_engine',
      templateUrl: 'admin/rules_engine/index.html',
      abstract: true,
      data: {
        title: 'Rules Engine'
      }
    });
  }
})();
