// Heavily inspired by: https://github.com/kubernetes/kubernetes/blob/1284c99ec9eedeb95d8048f0b1ceb7a4fc5a45ca/staging/src/k8s.io/client-go/discovery/cached/disk/round_tripper.go#L37
package disk

import (
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gregjones/httpcache"
	"github.com/gregjones/httpcache/diskcache"
	"github.com/peterbourgon/diskv"
)

type cacheRoundTripper struct {
	rt  *httpcache.Transport
	log logrus.FieldLogger
}

// NewCacheRoundTripper creates a roundtripper that reads the ETag on
// response headers and send the If-None-Match header on subsequent
// corresponding requests.
func NewCacheRoundTripper(cacheDir string, rt http.RoundTripper, log logrus.FieldLogger) http.RoundTripper {
	d := diskv.New(diskv.Options{
		PathPerm: os.FileMode(0750),
		FilePerm: os.FileMode(0660),
		BasePath: cacheDir,
		TempDir:  filepath.Join(cacheDir, ".diskv-temp"),
	})
	t := httpcache.NewTransport(diskcache.NewWithDiskv(d))
	t.Transport = rt

	return &cacheRoundTripper{rt: t, log: log}
}

func (rt *cacheRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return rt.rt.RoundTrip(req)
}

func (rt *cacheRoundTripper) CancelRequest(req *http.Request) {
	type canceler interface {
		CancelRequest(*http.Request)
	}
	if cr, ok := rt.rt.Transport.(canceler); ok {
		cr.CancelRequest(req)
	} else {
		rt.log.Errorf("CancelRequest not implemented by %T", rt.rt.Transport)
	}
}

func (rt *cacheRoundTripper) WrappedRoundTripper() http.RoundTripper { return rt.rt.Transport }
