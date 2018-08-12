var crypto = require('crypto');
var rp = require('request-promise');

var APP_KEY = '5b6bfaf6a6dd527199fce0c1';
var APP_SECRECT = 'f3f0fad23d0803d01619ad06b9cd1469c9ac61a37d794d76cafd5247c272fe38259a07784dfde3cad9e6a6a55340ed17';
var ADDRESS = '1Ny1Y6cxqqq9gJm2euqQh1q6oo1jWhUu7u';

this.upload = function (data) {
  var sign = crypto.createHash('sha1').update(APP_KEY + data + APP_SECRECT).digest('hex');

  var options = {
    method: 'POST',
    uri: 'http://chromeapi.genyuanlian.com:9005/api/upload/org',
    body: {
        data   : data,
        app_key: APP_KEY,
        sign   : sign
    },
    json: true // Automatically stringifies the body to JSON
  };

  return rp(options);
};

module.exports = this;