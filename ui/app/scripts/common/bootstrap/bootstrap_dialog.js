(function () {
  'use strict';

  angular.module('foxcommUiApp').factory('bootstrap.dialog', ['$modal', '$templateCache', modalDialog]);

  function modalDialog($modal, $templateCache) {
    var service = {
      deleteDialog: deleteDialog,
      confirmationDialog: confirmationDialog
    };

    $templateCache.put('modalDialog.tpl.html',
            '<div>' +
            '    <div class="modal-header">' +
            '        <button type="button" class="close" data-dismiss="modal" aria-hidden="true" data-ng-click="cancel()">&times;</button>' +
            '        <h3>{{title}}</h3>' +
            '    </div>' +
            '    <div class="modal-body">' +
            '        <p>{{message}}</p>' +
            '    </div>' +
            '    <div class="modal-footer">' +
            '        <button class="btn btn-primary" data-ng-click="ok()">{{okText}}</button>' +
            '        <button class="btn btn-info" data-ng-click="cancel()">{{cancelText}}</button>' +
            '    </div>' +
            '</div>');

    return service;

    function deleteDialog(itemName) {
      var title = 'Confirm Delete';
      itemName = itemName || 'item';
      var msg = 'Are you sure you want to delete ' + itemName + '?';
      var templateUrl = "modalDialog.tpl.html";

      return confirmationDialog(title, '', templateUrl, msg);
    }


    function confirmationDialog(title, resource, templateUrl, msg, okText, cancelText, extraOpts) {

      var modalOptions = {
        templateUrl: templateUrl,
        controller: ModalInstance,
        keyboard: true,
        resolve: {
          options: function () {
            return {
              title: title,
              message: msg,
              okText: okText,
              resource: resource,
              templateUrl: templateUrl,
              cancelText: cancelText,
              extraOpts: extraOpts
            };
          }
        }
      };

      return $modal.open(modalOptions).result;
    }
  }

  var ModalInstance = ['$scope', '$modalInstance', 'options',
    function ($scope, $modalInstance, options) {
      $scope.newRule = options.resource || {};
      $scope.title = options.title || 'Title';
      $scope.templateUrl = options.templateUrl || 'modalDialog.tpl.html';
      $scope.message = options.message || '';
      $scope.okText = options.okText || 'OK';
      $scope.cancelText = options.cancelText || 'Cancel';
      $scope.extraOpts = options.extraOpts || {};
      $scope.ok = function () {
        $modalInstance.close($scope.newRule);
      };
      $scope.cancel = function () { $modalInstance.dismiss('cancel'); };
    }];
})();
