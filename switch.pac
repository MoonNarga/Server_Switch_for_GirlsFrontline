function FindProxyForURL(url, host) {
    if ( dnsDomainIs(host, "gfcn-transit.gw.sunborngame.com")) {
		return 'PROXY 175.24.18.94:8888';
    }

    if ( dnsDomainIs(host, "gfcn-transit.ios.sunborngame.com")) {
		return 'PROXY 175.24.18.94:8889';
    }

	return 'DIRECT';
}