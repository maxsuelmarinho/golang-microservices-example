package config

import (
	"testing"
	. "github.com/smartystreets/goconvey/convey"
	"gopkg.in/h2non/gock.v1"
	"github.com/spf13/viper"
)

const ServiceName = "accountservice"

func TestHandleRefreshEvent(t *testing.T) {
	viper.Set("config_server_url", "http://config-server:8888")
	viper.Set("profile", "test")
	viper.Set("config_branch", "master")

	defer gock.Off()
	gock.New("http://config-server:8888").
		Get("/accountservice/test/master").
		Reply(200).
		BodyString(`{"name":"accountservice-test", "profiles":["test"], "label":null, "version": null, "propertySources": [{"name": "file:/config-repo/accountservice-test.yml", "source": {"server_port": 8181, "server_name": "Account Service RELOADED"}}]}`)

	Convey("Given a refresh event received, targeting out application", t, func() {
		var body = `{"type":"RefreshRemoteApplicationEvent", "timestamp": 1494514362123, "originService": "config-server:docker:8888", "destinationService": "accountservice:**", "id": "53e61c71-cbae-4b6d-84bb-d0dcc0aeb4dc"}`
		
		Convey("When handled", func() {
			handleRefreshEvent([]byte(body), ServiceName)

			Convey("Then Viper should have been re-populated with values from source", func() {
				So(viper.GetString("server_name"), ShouldEqual, "Account Service RELOADED")
			})
		})
	})
}

func TestHandleRefreshEventForOtherApplication(t *testing.T) {
	gock.Intercept()

	defer gock.Off()

	Convey("Given a refresh event received, targeting another application", t, func() {
		var body = `{"type":"RefreshRemoteApplicationEvent", "timestamp": 1494514362123, "originService": "config-server:docker:8888", "destinationService": "quotesservice:**", "id": "53e61c71-cbae-4b6d-84bb-d0dcc0aeb4dc"}`

		Convey("When parsed", func() {
			handleRefreshEvent([]byte(body), ServiceName)

			Convey("Then no outgoing HTTP requests should have been intercepted", func() {
				So(gock.HasUnmatchedRequest(), ShouldBeFalse)
			})
		})
	})
}