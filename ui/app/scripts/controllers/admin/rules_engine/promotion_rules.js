(function() {
  'use strict';

  angular.module('foxcommUiApp')
      .config(['$stateProvider', PromotionRulesConfig])
      .controller('PromotionRulesCtrl', [
        'foxcommPromotionRule',
        'logger',
        'bootstrap.dialog',
        PromotionRulesCtrl]
  );

  function PromotionRulesConfig($stateProvider) {
    $stateProvider.state('admin.rules_engine.promotion_rules', {
      url: '/promotion-rules',
      templateUrl: 'admin/rules_engine/promotion_rules.html',
      controller: 'PromotionRulesCtrl as promotionRule',
      data: {
        title: 'Rules Engine'
      }
    });
  }

  function PromotionRulesCtrl(foxcommPromotionRule, logger, bsDialog) {
    var promotionRule = this;


    promotionRule.newRuleType, promotionRule.newPreferenceType = {};
    promotionRule.collapsed = true;

    var promotionsByActionHandler = function(data) {
      promotionRule.ruleSet = data;

      foxcommPromotionRule.getRuleOptions(promotionOptionsHandler);
    };


    var promotionOptionsHandler = function(data) {
      promotionRule.ruleTypes = data.RuleTypes;
      promotionRule.preferenceTypes = data.PreferenceTypes;
      promotionRule.adjustmentCalculators = data.CalculatorTypes;
    };

    foxcommPromotionRule.getByActionName('signup', promotionsByActionHandler);

    promotionRule.createNewRuleType = function() {
      promotionRule.ruleSet.Rules.push(promotionRule.newRuleType);
      save();
    };

    promotionRule.cancelNewRuleType = function() {

    };

    promotionRule.createNewPreferenceType = function() {
      promotionRule.ruleSet.Actions.push(promotionRule.newPreferenceType);

      save();
    };

    var save = function() {
      promotionRule.ruleSet.save().then(function() {
        logger.logSuccess("Successfully updated rule set", null, null, true);
      }, function(error) {
        logger.logError("An error ocurred. Please try again later.", null, null, true);
      });
    };

    promotionRule.cancelNewItem = function(element) {
      element.Type = '';
    };

    promotionRule.openDeleteModal = function(rule, type) {

      var markForDelete = _.indexOf(promotionRule.ruleSet[type], rule);
      bsDialog.deleteDialog(rule.Type)
          .then(function(){
            delete promotionRule.ruleSet[type][markForDelete];
            promotionRule.ruleSet[type] = _.compact(promotionRule.ruleSet[type]);
            save();
          })
    };

    promotionRule.openEditModal = function(action) {
      var template = 'admin/rules_engine/promotion_action_modal.html';
      bsDialog.confirmationDialog('Edit Action', action, template, '', 'Save')
          .then(function() {
            save();
          });
    };


  }
})();
