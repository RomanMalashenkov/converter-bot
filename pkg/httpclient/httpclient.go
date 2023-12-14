package httpclient

import (
	"crypto/tls"
	"net/http"
)

// CloneTransport создает копию HTTP транспорта с определенной конфигурацией
func CloneTransport() *http.Transport {
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.TLSClientConfig = &tls.Config{
		InsecureSkipVerify: true,
		// Другие параметры безопасности, если необходимо
	}
	return transport
}
