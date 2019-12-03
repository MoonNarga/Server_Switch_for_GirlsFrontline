function FindProxyForURL(url, host) {
    if ( dnsDomainIs(host, "intranet.domain.com")||
    dnsDomainIs(host, "intranet.domain.com") ) {
		return 'PROXY 47.98.34.99:8888';
    }
	
	return 'DIRECT';
}