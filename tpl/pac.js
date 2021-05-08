// PAC ref: https://developer.mozilla.org/en-US/docs/Web/HTTP/Proxy_servers_and_tunneling/Proxy_Auto-Configuration_PAC_file
// test PAC: https://app.thorsen.pm/proxyforurl
// test js: https://jshint.com


function FindProxyForURL(url, host) {
    var PROXY = 'SOCKS5 127.0.0.1:1082';
    var DIRECT = 'DIRECT';
    var DEFAULT = PROXY;

var block = [{{.Block}}];

var allow = [{{.Allow}}];

    var match = function (domain, list) {
        for (var i = 0; i < list.length; i++) {
            var e = list[i];
            if (domain == e || domain.endsWith('.' + e)) {
                return true;
            }
        }
        return false;
    };

    if (match(host, block)) {
        return PROXY;
    }
    if (match(host, allow)) {
        return DIRECT;
    }
    return DEFAULT;
}

// endsWith ref: https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/String/endsWith#polyfill
if (!String.prototype.endsWith) {
    String.prototype.endsWith = function(search, this_len) {
        if (this_len === undefined || this_len > this.length) {
            this_len = this.length;
        }
        return this.substring(this_len - search.length, this_len) === search;
    };
}
