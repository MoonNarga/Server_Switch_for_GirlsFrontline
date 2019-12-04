function FindProxyForURL(url, host) {
    if ( dnsDomainIs(host, "adr.transit.gf.ppgame.com") ||
    dnsDomainIs(host, "ios.transit.gf.ppgame.com")) {
		return 'PROXY 216.24.186.233:10081';
    }

	return 'DIRECT';
}
