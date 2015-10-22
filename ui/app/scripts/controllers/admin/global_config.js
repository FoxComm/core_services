(function() {
  'use strict';

  angular.module('foxcommUiApp')
    .config(['$stateProvider', GlobalConfigConfig])
    .controller('GlobalConfigCtrl', ['$scope', 'foxcommCore', GlobalConfigCtrl]);

  function GlobalConfigConfig($stateProvider) {
    $stateProvider.state('admin.global_config', {
      url: 'global_config',
      templateUrl: 'admin/global_config/index.html',
      controller: 'GlobalConfigCtrl as vm',
      data: {
        title: 'Global Config'
      }
    });
  }

  function GlobalConfigCtrl($scope, foxcommCore) {
    var vm = this;

    vm.store = {};

    var features = foxcommCore.all('features/');
    features.getList().then(function(allFeatures) {
      vm.features = allFeatures;
    });

    var store = foxcommCore.one('stores');
    store.get().then(function(data) {
      vm.store = data;
    });

    $scope.updateFeatureStatus = function(feature, status) {
      feature.Enabled = status;
      feature.put();
    };

    $scope.StatusBtnText = function(status) {
      return status ? "Turn off" : "Turn on";
    };

    vm.save = function() {
      foxcommCore.build(vm.store).put().then(function(data){
        vm.store = data;
      });
    };

    vm.cancel = function() {
      vm.store = {}
    };

  }
})();
