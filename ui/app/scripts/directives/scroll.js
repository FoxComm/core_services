(function(){
  'use strict';

  var app = angular.module('foxcommUiApp');

  app.directive('scrollWidget', function() {
    return {
      restrict: 'C',
      link: function(scope, element, attrs) {
        $(element).slimscroll({
          height: '350px',
          position: 'right',
          color: '#868686'
        });
      }
    }
  });
})();

