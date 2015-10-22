(function() {
  'use strict';

  angular.module('foxcommUiApp')
    .config(['$stateProvider', ManualAttributionConfig])
    .controller('ManualAttributionCtrl', [
      'foxcommManualAttributions',
      'logger',
      'bootstrap.dialog',
      ManualAttributionCtrl]
    );

  function ManualAttributionConfig($stateProvider) {
    $stateProvider.state('admin.rules_engine.manual_attributions', {
      url: '/manual-attributions',
      templateUrl: 'admin/rules_engine/manual_attributions.html',
      controller: 'ManualAttributionCtrl as ma',
      data: {
        title: 'Rules Engine'
      }
    });
  }

  function ManualAttributionCtrl(foxcommManualAttributions, logger, bsDialog) {
    var ma = this;
    ma.openCreateModal = openCreateModal;
    ma.openEditModal = openEditModal;
    ma.openDeleteModal = openDeleteModal;
    ma.attributions = [];
    ma.getAllAttributions = getAllAttributions;
    ma.getAllEntities = getAllEntities;

    var attributionsHandler = function(attributions) {
      ma.attributions = attributions;
    }

    var entitiesHandler = function(entities) {
      ma.entities = entities;
    }

    ma.getAllAttributions();

    function getAllAttributions() {
      foxcommManualAttributions.all('manual_attributions').getList().then(function(attributions){
        ma.attributions = attributions;
      });
    }

    function getAllEntities(query) {
      return foxcommManualAttributions.all("entities").getList({query: query}).then(function(entities){
        ma.entities = entities;
        return entities;
      });
    }

    function openCreateModal() {
      var template = 'admin/rules_engine/manual_attribution_modal.html';
      ma.newAttribution = {ReferralActivity: "manual"}
      bsDialog.confirmationDialog('Create Manual Attribution', ma, template, '', 'Save', '', {}).then(function(data){
        foxcommManualAttributions.all('manual_attributions').post(ma.newAttribution).then(function(data) {
          ma.getAllAttributions();
          logger.logSuccess("Successfully created manual attribution", null, null, true);
        }, function(data) {
          logger.logError("An error ocurred. Please try again later.", null, null, true);
        });
      });
    }

    function openEditModal(attribution) {
      var template = 'admin/rules_engine/manual_attribution_modal.html';
      bsDialog.confirmationDialog('Edit Manual Attribution', ma, template, '', 'Save', '', {}).then(function(data){
        attribution.save().then(function(data) {
          ma.getAllAttributions();
          logger.logSuccess("Successfully created manual attribution", null, null, true);
        }, function(data) {
          logger.logError("An error ocurred. Please try again later.", null, null, true);
        });
      });
    }

    function openDeleteModal(attribution) {
      ma.attrTobeRemoved = attribution;
      bsDialog.deleteDialog().then(function(data){
        ma.attrTobeRemoved.remove().then(function(data) {
          ma.getAllAttributions();
          logger.logSuccess("Successfully deleted manual attribution", null, null, true); 
        })
      }, function(){
        logger.logError("An error ocurred. Please try again later.", null, null, true);
      })
    }
  }
})();
