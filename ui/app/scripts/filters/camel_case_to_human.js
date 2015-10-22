'use strict';

// camelCase To Human Filter
// ---------------------
// Converts a camelCase string to a human readable string.
// i.e. myVariableName => My Variable Name


angular.module('foxcommUiApp').filter('camelCaseToHuman', function() {
  return function(input) {
    return input.charAt(0).toUpperCase() + input.substr(1).replace(/[A-Z]/g, ' $&');
  }
});
