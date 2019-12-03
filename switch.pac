function FindProxyForURL(url, host) {
    if ( dnsDomainIs(host, "adr.transit.gf.ppgame.com")||
    dnsDomainIs(host, "ios.transit.gf.ppgame.com") ) {
		return 'PROXY 47.98.34.99:8888';
    }
	
	return 'DIRECT';
}