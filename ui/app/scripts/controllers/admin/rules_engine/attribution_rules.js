(function() {
  'use strict';

  angular.module('foxcommUiApp')
    .config(['$stateProvider', AttributionRulesConfig])
    .controller('AttributionRulesCtrl', [
      'foxcommAttributionRule',
      'bootstrap.dialog',
      AttributionRulesCtrl]
    );

  function AttributionRulesConfig($stateProvider) {
    $stateProvider.state('admin.rules_engine.attribution_rules', {
      url: '/attribution-rules',
      templateUrl: 'admin/rules_engine/attribution_rules.html',
      controller: 'AttributionRulesCtrl as attributionRule',
      data: {
        title: 'Rules Engine'
      }
    });
  }

  function AttributionRulesCtrl(foxcommAttributionRule, bsDialog, camelCaseToHuman) {
    var attributionRule = this;

    attributionRule.newRule = {};

    attributionRule.isValidActionRule = function(action) {
      return action == 'signup' || action == 'checkout';
    };

    var attributionActionsHandler = function(data){
      attributionRule.actions = data.site_actions;
      attributionRule.ruleTypes = data.action_rules;
      attributionRule.ActivityName = 'signup';

      foxcommAttributionRule.getByActionName(attributionRule.ActivityName, attributionsByActionHandler);
    };

    var attributionsByActionHandler = function(data) {
      $('.selectpicker').selectpicker();
      attributionRule.ruleSet = data;
      attributionRule.ruleSet.RuleStack = data.RuleStack == null ? [] : data.RuleStack;
    };

    foxcommAttributionRule.getAttributionActions(attributionActionsHandler);

    attributionRule.getByActionName = function(actionName) {
      attributionRule.ActivityName = actionName;
      foxcommAttributionRule.getByActionName(actionName, attributionsByActionHandler);
    };

    attributionRule.goBack = function(){
      modalInstance.dismiss('cancel');
    };

    attributionRule.openCreateModal = function() {
      if (attributionRule.ruleSet.Id == "") {
        foxcommAttributionRule.create(attributionRule.ActivityName, function(data){
          attributionRule.ruleSet.Id = data.Id;
        });
      }

      var templateType = attributionRule.ruleType.split('.')[1];

      attributionRule.newRule.RuleName = templateType;

      templateType = templateType.charAt(0) + templateType.substr(1).replace(/[A-Z]/g, '_$&');

      var template = 'admin/rules_engine/' + templateType.toLowerCase() + '_modal.html';



      bsDialog.confirmationDialog('Create Rule', attributionRule.newRule, template, '', 'Save')
          .then(function(data){
            var ruleType = {};

            switch (attributionRule.ruleType.split('.')[1]) {

              case "ClickAttributionRule":
                ruleType = {
                  ClickAttributionRule: data,
                  RuleType: attributionRule.ruleType
                };
                break;

              case "LinkToParentRule":
                ruleType = {
                  LinkToParentRule: data,
                  RuleType: attributionRule.ruleType
                };
                break;

            }
            if (attributionRule.newRule.WeightTable) {
              attributionRule.newRule.WeightTable = JSON.parse("[" + attributionRule.newRule.WeightTable + "]");
            }
            attributionRule.ruleSet.RuleStack.push(ruleType);
            attributionRule.saveRuleSet();
          });

    };

    attributionRule.openEditModal = function(rule) {

      var templateType = rule.RuleType.split('.')[1];
      templateType = templateType.charAt(0) + templateType.substr(1).replace(/[A-Z]/g, '_$&');

      var template = 'admin/rules_engine/' + templateType.toLowerCase() + '_modal.html';
      attributionRule.newRule = rule[rule.RuleType.split('.')[1]];
      bsDialog.confirmationDialog('Edit Rule', attributionRule.newRule, template, '', 'Save')
          .then(function() {
            attributionRule.newRule.WeightTable = JSON.parse("[" + attributionRule.newRule.WeightTable + "]");
            attributionRule.saveRuleSet();
          });
    };

    attributionRule.openDeleteModal = function(rule) {

      var markForDelete = _.indexOf(attributionRule.ruleSet.RuleStack, rule);
      bsDialog.deleteDialog(rule.RuleName)
          .then(function(){
            delete attributionRule.ruleSet.RuleStack[markForDelete];
            attributionRule.ruleSet.RuleStack = _.compact(attributionRule.ruleSet.RuleStack);
            attributionRule.saveRuleSet();
          })
    };


    attributionRule.saveRuleSet = function(){
      attributionRule.ruleSet.fromServer = false;
      attributionRule.ruleSet.ActivityName = attributionRule.ActivityName;
      attributionRule.ruleSet.put();
    };

  }
})();
