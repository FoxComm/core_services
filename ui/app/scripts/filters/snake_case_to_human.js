'use strict';

// snake_case To Human Filter
// ---------------------
// Converts a snake_case string to a human readable string.
// i.e. variable_name => My Variable Name


angular.module('foxcommUiApp').filter('snakeCaseToHuman', function() {
  return function(input) {
    return input == undefined ? '' : input.charAt(0) + input.substr(1).replace(/_/, ' ');
  }
});
