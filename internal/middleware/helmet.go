package middleware

import (
	"github.com/XDoubleU/essentia/internal/core"
	"github.com/goddtriffin/helmet"
)

func Helmet() core.HandlerFunc {
	helmet := helmet.Default()

	return func(c *core.Context) {
		helmet.ContentSecurityPolicy.Header(c.Writer.ResponseWriter)
		helmet.XContentTypeOptions.Header(c.Writer.ResponseWriter)
		helmet.XDNSPrefetchControl.Header(c.Writer.ResponseWriter)
		helmet.XDownloadOptions.Header(c.Writer.ResponseWriter)
		helmet.ExpectCT.Header(c.Writer.ResponseWriter)
		helmet.FeaturePolicy.Header(c.Writer.ResponseWriter)
		helmet.XFrameOptions.Header(c.Writer.ResponseWriter)
		helmet.XPermittedCrossDomainPolicies.Header(c.Writer.ResponseWriter)
		helmet.XPoweredBy.Header(c.Writer.ResponseWriter)
		helmet.ReferrerPolicy.Header(c.Writer.ResponseWriter)
		helmet.StrictTransportSecurity.Header(c.Writer.ResponseWriter)
		helmet.XXSSProtection.Header(c.Writer.ResponseWriter)

		c.Next()
	}
}
