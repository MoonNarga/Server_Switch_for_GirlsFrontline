function FindProxyForURL(url, host) {
    if ( dnsDomainIs(host, "adr.transit.gf.ppgame.com")) {
		return 'PROXY 175.24.18.94:8888';
    }

    if ( dnsDomainIs(host, "ios.transit.gf.ppgame.com")) {
		return 'PROXY 175.24.18.94:8889';
    }

	return 'DIRECT';
}