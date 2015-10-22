(function() {
  'use strict';

  angular.module('foxcommUiApp')
      .config(['$stateProvider', BalanceRulesConfig])
      .controller('BalanceRulesCtrl', [
        'foxcommBalanceRule',
        BalanceRulesCtrl]
  );

  function BalanceRulesConfig($stateProvider) {
    $stateProvider.state('admin.rules_engine.balance_rules', {
      url: '/balance-rules',
      templateUrl: 'admin/rules_engine/balance_rules.html',
      controller: 'BalanceRulesCtrl as balanceRule',
      data: {
        title: 'Rules Engine'
      }
    });
  }
  
  function BalanceRulesCtrl(foxcommBalanceRule) {
    var balanceRule = this;

    balanceRule.newRule = {};
    balanceRule.ruleTypes = ["ExpirationRule", "SpendingThresholdRule", "TimeThresholdRule"];

    var balanceCurrencyHandler = function(data){
      $('.selectpicker').selectpicker();
      balanceRule.ruleSet = data;
    };

    foxcommBalanceRule.getByCurrencyName('loyalty_points', balanceCurrencyHandler);


    balanceRule.saveRuleSet = function(){
      balanceRule.ruleSet.put();
    };

    balanceRule.removeRule = function(rule) {
      _.each(rule, function(v, k) { rule[k] = 0;});
      balanceRule.ruleSet.put();
    }
  }
})();
