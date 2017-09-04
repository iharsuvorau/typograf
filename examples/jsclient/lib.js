(function (exports) {
  var lib = exports.lib || {};

  // sends a request with optional data (PUT, POST) and exectutes the provided function on success
  lib.request = function (verb, url, data, successFunc, errFunc) {
    if (!errFunc) {
      errFunc = function(response) {
        console.log('error: ', verb, url, response);
      };
    }

    if (!successFunc) {
      successFunc = function(response) {
        console.log('success: ', verb, url, response);
      };
    }

    var httpRequest = new XMLHttpRequest();
    httpRequest.open(verb, url);
    httpRequest.setRequestHeader("Content-Type", "application/json;charset=UTF-8");
    httpRequest.send(JSON.stringify(data));
    httpRequest.onreadystatechange = function () {
      var response;
      if (httpRequest.readyState === XMLHttpRequest.DONE) {
        if (httpRequest.status === 200) {
          response = JSON.parse(httpRequest.responseText);
          successFunc(response);
        } else {
          if (httpRequest.responseText.length > 0) {
            response = JSON.parse(httpRequest.responseText);
            errFunc(response);
          }
        }
      }
    };
  };

  // updates the lib property
  exports.lib = lib;
  
})(window);
