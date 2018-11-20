package vistar

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseArray(t *testing.T) {
	params := map[string]string{
		"p1": "v1",
		"p2": "   First,  Second, Third Value , ,, Fourth  ",
		"p3": ",,, ,, ,   ,,, ",
	}

	conf := &adConfig{}
	def := []string{"a", "b"}
	res := conf.parseArray(params, "unknown", def)
	assert.Equal(t, res, def)

	res = conf.parseArray(params, "p2", def)
	assert.Equal(t, res, []string{"First", "Second", "Third Value", "Fourth"})

	res = conf.parseArray(params, "p3", def)
	assert.Equal(t, res, def)
}

func TestParseInt(t *testing.T) {
	params := map[string]string{
		"p1": "0",
		"p2": "100",
		"p3": "not-number",
	}

	conf := &adConfig{}
	res := conf.parseInt(params, "unknown", 6)
	assert.Equal(t, res, int64(6))

	res = conf.parseInt(params, "p1", 6)
	assert.Equal(t, res, int64(0))

	res = conf.parseInt(params, "p2", 6)
	assert.Equal(t, res, int64(100))

	res = conf.parseInt(params, "p3", 6)
	assert.Equal(t, res, int64(6))
}

func TestParseFloat(t *testing.T) {
	params := map[string]string{
		"p1": "0.0",
		"p2": "100.9",
		"p3": "not-number",
	}

	conf := &adConfig{}
	res := conf.parseFloat(params, "unknown", 6.6)
	assert.Equal(t, res, 6.6)

	res = conf.parseFloat(params, "p1", 6.6)
	assert.Equal(t, res, float64(0))

	res = conf.parseFloat(params, "p2", 6.6)
	assert.Equal(t, res, 100.9)

	res = conf.parseFloat(params, "p3", 6.6)
	assert.Equal(t, res, 6.6)
}

func TestParseBool(t *testing.T) {
	params := map[string]string{
		"p1": "false",
		"p2": "true",
		"p3": "not-number",
	}

	conf := &adConfig{}
	res := conf.parseBool(params, "unknown", true)
	assert.True(t, res)

	res = conf.parseBool(params, "p1", true)
	assert.False(t, res)

	res = conf.parseBool(params, "p2", false)
	assert.True(t, res)

	res = conf.parseBool(params, "p3", false)
	assert.False(t, res)
}

func TestParse(t *testing.T) {
	params := map[string]string{
		"vistar.url":               "staging-url",
		"vistar.api_key":           "api-key",
		"vistar.network_id":        "network-id",
		"vistar.venue_id":          "venue-id",
		"vistar.direct_connection": "true",
		"vistar.latitude":          "45.5",
		"vistar.longitude":         "44.4",
		"vistar.mime_types":        "a,b,c",
		"vistar.width":             "100",
		"vistar.height":            "200",
		"vistar.allow_audio":       "true",
		"vistar.static_duration":   "9",
	}

	conf := &adConfig{}
	conf.parse(params)

	assert.Equal(t, conf.url, "staging-url")
	assert.Equal(t, conf.baseRequest.ApiKey, "api-key")
	assert.Equal(t, conf.baseRequest.NetworkId, "network-id")
	assert.Equal(t, conf.baseRequest.DeviceId, "venue-id")
	assert.Equal(t, conf.baseRequest.VenueId, "venue-id")
	assert.True(t, conf.baseRequest.DirectConnection)
	assert.Equal(t, conf.baseRequest.Latitude, 45.5)
	assert.Equal(t, conf.baseRequest.Longitude, 44.4)
	assert.Equal(t, conf.baseRequest.DisplayTime, int64(0))
	assert.Equal(t, conf.baseRequest.NumberOfScreens, int64(1))
	assert.Len(t, conf.baseRequest.DisplayAreas, 1)
	assert.Len(t, conf.baseRequest.DeviceAttributes, 2)

	assert.Equal(t, conf.baseRequest.DisplayAreas[0].Width, int64(100))
	assert.Equal(t, conf.baseRequest.DisplayAreas[0].Height, int64(200))
	assert.True(t, conf.baseRequest.DisplayAreas[0].AllowAudio)
	assert.Equal(t, conf.baseRequest.DisplayAreas[0].SupportedMedia,
		[]string{"a", "b", "c"})
	assert.Equal(t, conf.baseRequest.DisplayAreas[0].StaticDuration, int64(9))
}

func TestUpdateAdRequest(t *testing.T) {
	params := map[string]string{
		"vistar.url":               "staging-url",
		"vistar.api_key":           "api-key",
		"vistar.network_id":        "network-id",
		"vistar.venue_id":          "venue-id",
		"vistar.direct_connection": "true",
		"vistar.latitude":          "45.5",
		"vistar.longitude":         "44.4",
		"vistar.mime_types":        "a,b,c",
		"vistar.width":             "100",
		"vistar.height":            "200",
		"vistar.allow_audio":       "true",
		"vistar.static_duration":   "9",
	}

	conf := &adConfig{}
	conf.parse(params)

	req := &AdRequest{
		DisplayAreas: []DisplayArea{
			DisplayArea{Id: "d1", Width: 500, Height: 500, AllowAudio: false,
				SupportedMedia: []string{"image"}},
		},
		DeviceAttributes: []DeviceAttribute{
			DeviceAttribute{Name: "attr1", Value: "value1"},
			DeviceAttribute{Name: "attr2", Value: "value2"},
		},
	}

	conf.UpdateAdRequest(req)

	assert.Equal(t, req.ApiKey, "api-key")
	assert.Equal(t, req.NetworkId, "network-id")
	assert.Equal(t, req.DeviceId, "venue-id")
	assert.Equal(t, req.VenueId, "venue-id")
	assert.True(t, req.DirectConnection)
	assert.Equal(t, req.Latitude, 45.5)
	assert.Equal(t, req.Longitude, 44.4)
	assert.NotEqual(t, req.DisplayTime, int64(0))
	assert.Equal(t, req.NumberOfScreens, int64(1))
	assert.Len(t, req.DisplayAreas, 1)
	assert.Len(t, req.DeviceAttributes, 4)

	assert.Equal(t, req.DisplayAreas[0].Width, int64(500))
	assert.Equal(t, req.DisplayAreas[0].Height, int64(500))
	assert.False(t, req.DisplayAreas[0].AllowAudio)
	assert.Equal(t, req.DisplayAreas[0].SupportedMedia, []string{"image"})
	assert.Equal(t, req.DisplayAreas[0].StaticDuration, int64(0))
}