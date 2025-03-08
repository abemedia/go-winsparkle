package winsparkle_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"text/template"
	"time"

	"github.com/abemedia/go-winsparkle"
	_ "github.com/abemedia/go-winsparkle/dll"
)

func TestWinSparkle(t *testing.T) {
	winsparkle.SetAppDetails("Test", "Test", "1.0")
	winsparkle.SetAppcastURL(server(t, "1.0"))

	last := winsparkle.GetLastCheckTime()

	winsparkle.Init()
	defer winsparkle.Cleanup()

	winsparkle.CheckUpdateWithoutUI()
	time.Sleep(time.Second) // Wait for update check.

	if !winsparkle.GetLastCheckTime().After(last) {
		t.Fatal("should check for updates")
	}
}

func TestSetErrorCallback(t *testing.T) {
	winsparkle.SetAppDetails("Test", "Test", "1.0")
	winsparkle.SetAppcastURL("nope")
	winsparkle.Init()
	defer winsparkle.Cleanup()

	ch := make(chan struct{}, 1)
	winsparkle.SetErrorCallback(func() {
		ch <- struct{}{}
	})

	winsparkle.CheckUpdateWithoutUI()

	select {
	case <-ch:
	case <-time.After(time.Second):
		t.Error("should call callback")
	}
}

func TestSetCanShutdownCallback(t *testing.T) {
	winsparkle.SetAppDetails("Test", "Test", "1.0.0")
	winsparkle.SetAppcastURL(server(t, "2.0.0"))
	winsparkle.Init()
	defer winsparkle.Cleanup()

	ch := make(chan struct{}, 1)
	winsparkle.SetCanShutdownCallback(func() bool {
		ch <- struct{}{}
		return true
	})

	winsparkle.CheckUpdateWithUIAndInstall()

	select {
	case <-ch:
	case <-time.After(time.Second):
		t.Error("should call callback")
	}
}

func server(t *testing.T, version string) string {
	t.Helper()
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/xml")
		appcastTmpl.Execute(w, struct {
			Version string
			Host    string
		}{version, r.Host})
	}))
	t.Cleanup(s.Close)
	return s.URL
}

//nolint:lll
var appcastTmpl = template.Must(template.New("appcast").Parse(`<?xml version="1.0" encoding="utf-8"?>
<rss version="2.0" xmlns:sparkle="http://www.andymatuschak.org/xml-namespaces/sparkle">
	<channel>
		<title>WinSparkle Test Appcast</title>
		<description>Most recent updates to WinSparkle Test</description>
		<language>en</language>
		<item>
			<title>Version {{.Version}}</title>
			<description>This is an update.</description>
			<pubDate>Mon, 28 Jan 2013 14:30:00 +0500</pubDate>
			<enclosure sparkle:version="{{.Version}}" url="http://{{.Host}}/install.msi" length="0" type="application/octet-stream"/>
		</item>
	</channel>
</rss>`))
