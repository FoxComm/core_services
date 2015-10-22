(function() {
  'use strict';

  angular.module('foxcommUiApp')
      .config(['$stateProvider', BackupsConfig])
      .controller('BackupsCtrl', ['foxcommBackups', 'logger', '$interval', '$timeout', '$scope', BackupsCtrl]);

  function BackupsConfig($stateProvider) {
    $stateProvider.state('admin.backups', {
      url: 'backups',
      templateUrl: 'admin/backups/index.html',
      controller: 'BackupsCtrl as backups',
      data: {
        title: 'Backups'
      }
    });
  }

  function BackupsCtrl(foxcommCore, logger, $interval, $timeout, $scope) {
    var backups = this;
    backups.stats = {};
    backups.settings = {};
    backups.actions = [
      ["core_database", "Foxcomm core Db"],
      ["feature_databases", "Foxcomm features Db's"],
      ["origin_database", "Origin Db"],
      ["assets", "Assets"],
    ];
    backups.selectedActions = [];

    backups.saveSettings = function() {
      backups.settings.save();
    }

    backups.enqueueBackup = function(){
      foxcommCore.all("jobs").post({"actions": backups.selectedActions}).then(function(){
        logger.logSuccess("Jobs Enqueued", null, null, true)
        backups.fetchStats();
      });
    }

    backups.fetchSettings = function() {
      var settings = foxcommCore.one('settings');
      settings.get().then(function(settings) {
        backups.settings = settings;
      });
    }

    backups.fetchStats = function() {
      var stats = foxcommCore.one('stats');
      stats.get().then(function(stats){
        backups.stats = stats;
      });
    }

    backups.pollStats = function(){
      backups.poll = $interval(function(){
        backups.fetchStats();
      }, 3000);
    }

    backups.fetchSettings();
    backups.fetchStats();
    backups.pollStats();
  }
})();

