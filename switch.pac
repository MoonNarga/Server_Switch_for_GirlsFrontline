function FindProxyForURL(url, host) {
    if ( dnsDomainIs(host, "adr.transit.gf.ppgame.com")||
    dnsDomainIs(host, "ios.transit.gf.ppgame.com")) {
		return 'PROXY 175.24.18.94:8888';
    }

	return 'DIRECT';
}