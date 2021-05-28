// PAC ref: https://developer.mozilla.org/en-US/docs/Web/HTTP/Proxy_servers_and_tunneling/Proxy_Auto-Configuration_PAC_file
// template: https://github.com/clowwindy/gfwlist2pac
// test PAC: https://app.thorsen.pm/proxyforurl
// test js: https://jshint.com

function FindProxyForURL(url, host) {
    var proxy = "SOCKS5 127.0.0.1:6153;HTTP 127.0.0.1:6152;DIRECT;";  // 'PROXY' or 'SOCKS5' or 'HTTPS' or 'HTTP'
    var direct = 'DIRECT;';

    // var domains = {
    //     "google.com": 1,
    //     "youtube.com": 1,
    // };
    var domains = { {{.Block}}    };

    var hasOwnProperty = Object.hasOwnProperty;

    var suffix;
    var pos = host.lastIndexOf('.');
    pos = host.lastIndexOf('.', pos - 1);
    while(1) {
        if (pos == -1) {
            if (hasOwnProperty.call(domains, host)) {
                return proxy;
            } else {
                return direct;
            }
        }
        suffix = host.substring(pos + 1);
        if (hasOwnProperty.call(domains, suffix)) {
            return proxy;
        }
        pos = host.lastIndexOf('.', pos - 1);
    }
}
