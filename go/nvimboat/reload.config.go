package nvimboat

import (
	"net/http"
	"time"

	"github.com/EinYakAmNil/Nvimboat/go/engine/reload"
	"github.com/EinYakAmNil/Nvimboat/go/engine/reload/mangapill"
)

var CustomReload = map[string]func(
	r reload.Reloader,
	url string,
	header http.Header,
	cacheTime time.Duration,
	cacheDir string,
){
	"https://mangapill.com": func(r reload.Reloader, u string, h http.Header, ct time.Duration, cd string) {
		r.(*mangapill.MangapillReloader).GetRss(u, h, ct, cd)
	},
}
