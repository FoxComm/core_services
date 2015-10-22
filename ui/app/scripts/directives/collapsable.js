(function(){
  'use strict';

  var app = angular.module('foxcommUiApp');

  app.factory('CollapsableItemArray', function(){
    var collapsableItems = [];
    return {
      getCollapsableItems: function() {
        return collapsableItems;
      },
      setCollapsableItem: function(collapsableItem) {
        collapsableItems.push(collapsableItem);
      }
    };
  });

  app.controller('CollapsableItemCtrl', ['$scope', 'CollapsableItemArray', function($scope, CollapsableItemArray){
    // The list to be shown when CollapsableItemLink(anchor) is clicked
    $scope.collapsableItemGroup = undefined;

    // Function to add the collapsableItem in the CollapsableItemArray
    $scope.addCollapsableItem = function(collapsableItem) {
      CollapsableItemArray.setCollapsableItem(collapsableItem);
    };

    $scope.collapsableItemOnClick = function() {
      if (!$scope.collapsableActive) {
        angular.forEach(CollapsableItemArray.getCollapsableItems(), function(collapsableItem) {
          if (collapsableItem != $scope) {
            collapsableItem.hideCollapsableGroup();
          }
        });

        $scope.showCollapsableGroup();
      }
    };

    $scope.showCollapsableGroup = function() {
      $scope.collapsableActive = true;
      $scope.collapsableItemGroup.removeClass('is-hidden').addClass('is-shown');
      $scope.safeApply($scope);
    };

    $scope.hideCollapsableGroup = function() {
      $scope.collapsableActive = false;
      $scope.collapsableItemGroup.removeClass('is-shown').addClass('is-hidden');
      $scope.safeApply($scope);
    };

    // Needed function to upload the scope
    $scope.safeApply = function(scope, fn) {
      if(!(scope.$$phase || scope.$root.$$phase)){
        scope.$apply(fn);
      }
    };
  }]);

  app.directive('collapsableItem', function() {
    return {
      restrict: 'E',
      transclude: true,
      replace: true,
      template: '<li class="collapsable-item" ng-transclude ng-class="{collapsableActive: collapsableActive, collapsableClosed: !collapsableActive}"></li>',
      scope: {
        collapsableActive: "@"
      },
      controller: 'CollapsableItemCtrl',
      link: function($scope, elem, attrs) {
        $scope.addCollapsableItem($scope);
        $scope.collapsableItemGroup = $(elem).children('.collapsable-group:first');

        if ($scope.collapsableActive) {
          $scope.showCollapsableGroup();
        }

        $(elem).children('.collapsable-item-link:first').on('click', $scope.collapsableItemOnClick);
      }
    };
  });
})();
