(function() {
  'use strict';

  angular.module('foxcommUiApp')
      .config(['$stateProvider', AccumulationRulesConfig])
      .controller('AccumulationRulesCtrl', [
        'foxcommAccumulationRule',
        'bootstrap.dialog',
        AccumulationRulesCtrl]
  );

  function AccumulationRulesConfig($stateProvider) {
    $stateProvider.state('admin.rules_engine.accumulation_rules', {
      url: '/accumulation-rules',
      templateUrl: 'admin/rules_engine/accumulation_rules.html',
      controller: 'AccumulationRulesCtrl as accumulationRule',
      data: {
        title: 'Rules Engine'
      }
    });
  }

  function AccumulationRulesCtrl(foxcommAccumulationRule, bsDialog) {
    var accumulationRule = this;

    accumulationRule.newRule = {};

    accumulationRule.isValidActionRule = function(action) {
      return action == 'signup' || action == 'checkout';
    };

    var accumulationActionsHandler = function(data){
      accumulationRule.actions = data.site_actions;
      accumulationRule.signupRuleTypes = data.signup_action_rules;
      accumulationRule.checkoutRuleTypes = data.checkout_action_rules;
      accumulationRule.ActionName = 'signup';

      foxcommAccumulationRule.getByActionName(accumulationRule.ActionName, accumulationsByActionHandler);
    };

    var accumulationsByActionHandler = function(data) {
      $('.selectpicker').selectpicker();
      accumulationRule.ruleSet = data;
      accumulationRule.ruleSet.RuleStack = data.RuleStack == null ? [] : data.RuleStack;
    };

    foxcommAccumulationRule.getAccumulationActions(accumulationActionsHandler);

    accumulationRule.getByActionName = function(actionName) {
      accumulationRule.ActionName = actionName;
      foxcommAccumulationRule.getByActionName(actionName, accumulationsByActionHandler);
    };

    accumulationRule.goBack = function(){
      modalInstance.dismiss('cancel');
    };

    accumulationRule.openCreateModal = function() {
      if (accumulationRule.ruleSet.Id == "") {
        foxcommAccumulationRule.create(accumulationRule.ActionName, function(data){
          accumulationRule.ruleSet.Id = data.Id;
        });
      }

      var template = 'admin/rules_engine/' + accumulationRule.ActionName +'_rule_modal.html';

      bsDialog.confirmationDialog('Create Rule', accumulationRule.newRule, template, '', 'Save')
          .then(function(data){
            var ruleType = {};

            switch (accumulationRule.ruleType.split('.')[1]) {

              case "SignupCCRule":
                ruleType = {
                  SignupCCRule: data,
                  RuleType: accumulationRule.ruleType
                };
                break;

              case "LoyaltyPointsCheckoutRule":
                ruleType = {
                  LoyaltyPointsCheckoutRule: data,
                  RuleType: accumulationRule.ruleType
                };
                break;

              case "BeneficiaryCheckoutRule":
                ruleType = {
                  BeneficiaryCheckoutRule: data,
                  RuleType: accumulationRule.ruleType
                };
                break;

              case "OwnPurchaseCheckoutRule":
                ruleType = {
                  OwnPurchaseCheckoutRule: data,
                  RuleType: accumulationRule.ruleType
                };
                break;

            }
            if (accumulationRule.newRule.TreeLevels){
              accumulationRule.newRule.TreeLevels = JSON.parse("[" + accumulationRule.newRule.TreeLevels + "]");
            }
            accumulationRule.ruleSet.RuleStack.push(ruleType);
            accumulationRule.saveRuleSet();
          });

    };

    accumulationRule.openEditModal = function(rule) {
      var template = 'admin/rules_engine/' + accumulationRule.ActionName +'_rule_modal.html';
      accumulationRule.newRule = rule[rule.RuleType.split('.')[1]];
      bsDialog.confirmationDialog('Edit Rule', accumulationRule.newRule, template, '', 'Save')
          .then(function() {
            accumulationRule.newRule.TreeLevels = JSON.parse("[" + accumulationRule.newRule.TreeLevels + "]");
            accumulationRule.saveRuleSet();
          });
    };

    accumulationRule.openDeleteModal = function(rule) {
      var markForDelete = _.indexOf(accumulationRule.ruleSet.RuleStack, rule);
      bsDialog.deleteDialog()
          .then(function(){
            delete accumulationRule.ruleSet.RuleStack[markForDelete];
            accumulationRule.ruleSet.RuleStack = _.compact(accumulationRule.ruleSet.RuleStack);
            accumulationRule.saveRuleSet();
          })
    };

    accumulationRule.saveRuleSet = function(){
      accumulationRule.ruleSet.fromServer = false;
      accumulationRule.ruleSet.ActionName = accumulationRule.ActionName;

      accumulationRule.ruleSet.put();
    };

  }
})();
