package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type ProcessingOptionsTestSuite struct{ MainTestSuite }

func (s *ProcessingOptionsTestSuite) getRequest(url string) *http.Request {
	req, _ := http.NewRequest("GET", url, nil)
	return req
}

func (s *ProcessingOptionsTestSuite) TestParseBase64URL() {
	imageURL := "http://images.dev/lorem/ipsum.jpg?param=value"
	req := s.getRequest(fmt.Sprintf("http://example.com/encoded/%s?w=100&h=100", base64.RawURLEncoding.EncodeToString([]byte(imageURL))))
	ctx, err := parsePath(context.Background(), req)

	require.Nil(s.T(), err)
	assert.Equal(s.T(), imageURL, getImageURL(ctx))
}

func (s *ProcessingOptionsTestSuite) TestParseBase64URLWithBase() {
	conf.BaseURL = "http://images.dev/"

	imageURL := "lorem/ipsum.jpg?param=value"
	req := s.getRequest(fmt.Sprintf("http://example.com/encoded/%s?w=100&h=100", base64.RawURLEncoding.EncodeToString([]byte(imageURL))))
	ctx, err := parsePath(context.Background(), req)

	require.Nil(s.T(), err)
	assert.Equal(s.T(), conf.BaseURL + imageURL, getImageURL(ctx))
}

func (s *ProcessingOptionsTestSuite) TestParsePlainURL() {
	imageURL := "http://images.dev/lorem/ipsum.jpg"
	req := s.getRequest("http://example.com/" + imageURL)
	ctx, err := parsePath(context.Background(), req)

	require.Nil(s.T(), err)
	assert.Equal(s.T(), imageURL, getImageURL(ctx))
}


func (s *ProcessingOptionsTestSuite) TestParsePlainURLEscaped() {
	imageURL := "http://images.dev/lorem/ipsum.jpg?param=value"
	req := s.getRequest("http://example.com/" + url.PathEscape(imageURL))
	ctx, err := parsePath(context.Background(), req)

	require.Nil(s.T(), err)
	assert.Equal(s.T(), imageURL, getImageURL(ctx))
}

func (s *ProcessingOptionsTestSuite) TestParsePlainURLWithBase() {
	conf.BaseURL = "http://images.dev/"

	imageURL := "lorem/ipsum.jpg"
	req := s.getRequest("http://example.com/" + imageURL)
	ctx, err := parsePath(context.Background(), req)

	require.Nil(s.T(), err)
	assert.Equal(s.T(), conf.BaseURL + imageURL, getImageURL(ctx))
}

func (s *ProcessingOptionsTestSuite) TestParsePlainURLEscapedWithBase() {
	conf.BaseURL = "http://images.dev/"

	imageURL := "lorem/ipsum.jpg?param=value"
	req := s.getRequest("http://example.com/" + url.PathEscape(imageURL))
	ctx, err := parsePath(context.Background(), req)

	require.Nil(s.T(), err)
	assert.Equal(s.T(), conf.BaseURL + imageURL, getImageURL(ctx))
}

func (s *ProcessingOptionsTestSuite) TestParsePathBasic() {
	req := s.getRequest("http://example.com/unsafe/fill/100/200/noea/1/http://images.dev/lorem/ipsum.jpg?f=png&rs=fill&w=100&h=200&g=noea&el=1")
	ctx, err := parsePath(context.Background(), req)

	require.Nil(s.T(), err)

	po := getProcessingOptions(ctx)
	assert.Equal(s.T(), resizeFill, po.Resize)
	assert.Equal(s.T(), 100, po.Width)
	assert.Equal(s.T(), 200, po.Height)
	assert.Equal(s.T(), gravityNorthEast, po.Gravity.Type)
	assert.True(s.T(), po.Enlarge)
	assert.Equal(s.T(), imageTypePNG, po.Format)
}

func (s *ProcessingOptionsTestSuite) TestParsePathAdvancedFormat() {
	req := s.getRequest("http://example.com/http://images.dev/lorem/ipsum.jpg?format=webp")
	ctx, err := parsePath(context.Background(), req)

	require.Nil(s.T(), err)

	po := getProcessingOptions(ctx)
	assert.Equal(s.T(), imageTypeWEBP, po.Format)
}


func (s *ProcessingOptionsTestSuite) TestParsePathAdvancedResizingType() {
	req := s.getRequest("http://example.com/http://images.dev/lorem/ipsum.jpg?rt=fill")
	ctx, err := parsePath(context.Background(), req)

	require.Nil(s.T(), err)

	po := getProcessingOptions(ctx)
	assert.Equal(s.T(), resizeFill, po.Resize)
}

func (s *ProcessingOptionsTestSuite) TestParsePathAdvancedGravityFocuspoint() {
	req := s.getRequest("http://example.com/unsafe/gravity:fp:0.5:0.75/plain/http://images.dev/lorem/ipsum.jpg?gravity=fp:0.5:0.75")
	ctx, err := parsePath(context.Background(), req)

	require.Nil(s.T(), err)

	po := getProcessingOptions(ctx)
	assert.Equal(s.T(), gravityFocusPoint, po.Gravity.Type)
	assert.Equal(s.T(), 0.5, po.Gravity.X)
	assert.Equal(s.T(), 0.75, po.Gravity.Y)
}

func (s *ProcessingOptionsTestSuite) TestParsePathAdvancedQuality() {
	req := s.getRequest("http://example.com/plain/http://images.dev/lorem/ipsum.jpg?quality=55")
	ctx, err := parsePath(context.Background(), req)

	require.Nil(s.T(), err)

	po := getProcessingOptions(ctx)
	assert.Equal(s.T(), 55, po.Quality)
}

func (s *ProcessingOptionsTestSuite) TestParsePathAdvancedBackground() {
	req := s.getRequest("http://example.com/http://images.dev/lorem/ipsum.jpg?background=128:129:130")
	ctx, err := parsePath(context.Background(), req)

	require.Nil(s.T(), err)

	po := getProcessingOptions(ctx)
	assert.True(s.T(), po.Flatten)
	assert.Equal(s.T(), uint8(128), po.Background.R)
	assert.Equal(s.T(), uint8(129), po.Background.G)
	assert.Equal(s.T(), uint8(130), po.Background.B)
}

func (s *ProcessingOptionsTestSuite) TestParsePathAdvancedBackgroundHex() {
	req := s.getRequest("http://example.com/http://images.dev/lorem/ipsum.jpg?background=ffddee")
	ctx, err := parsePath(context.Background(), req)

	require.Nil(s.T(), err)

	po := getProcessingOptions(ctx)
	assert.True(s.T(), po.Flatten)
	assert.Equal(s.T(), uint8(0xff), po.Background.R)
	assert.Equal(s.T(), uint8(0xdd), po.Background.G)
	assert.Equal(s.T(), uint8(0xee), po.Background.B)
}

func (s *ProcessingOptionsTestSuite) TestParsePathAdvancedBackgroundDisable() {
	req := s.getRequest("http://example.com/http://images.dev/lorem/ipsum.jpg?background=")
	ctx, err := parsePath(context.Background(), req)

	require.Nil(s.T(), err)

	po := getProcessingOptions(ctx)
	assert.False(s.T(), po.Flatten)
}

func (s *ProcessingOptionsTestSuite) TestParsePathAdvancedBlur() {
	req := s.getRequest("http://example.com/plain/http://images.dev/lorem/ipsum.jpg?blur=0.2")
	ctx, err := parsePath(context.Background(), req)

	require.Nil(s.T(), err)

	po := getProcessingOptions(ctx)
	assert.Equal(s.T(), float32(0.2), po.Blur)
}

func (s *ProcessingOptionsTestSuite) TestParsePathAdvancedSharpen() {
	req := s.getRequest("http://example.com/plain/http://images.dev/lorem/ipsum.jpg?sharpen=0.2")
	ctx, err := parsePath(context.Background(), req)

	require.Nil(s.T(), err)

	po := getProcessingOptions(ctx)
	assert.Equal(s.T(), float32(0.2), po.Sharpen)
}
func (s *ProcessingOptionsTestSuite) TestParsePathAdvancedDpr() {
	req := s.getRequest("http://example.com/http://images.dev/lorem/ipsum.jpg?dpr=2")
	ctx, err := parsePath(context.Background(), req)

	require.Nil(s.T(), err)

	po := getProcessingOptions(ctx)
	assert.Equal(s.T(), 2.0, po.Dpr)
}
func (s *ProcessingOptionsTestSuite) TestParsePathAdvancedWatermark() {
	req := s.getRequest("http://example.com/http://images.dev/lorem/ipsum.jpg?watermark=0.5:soea:10:20:0.6")
	ctx, err := parsePath(context.Background(), req)

	require.Nil(s.T(), err)

	po := getProcessingOptions(ctx)
	assert.True(s.T(), po.Watermark.Enabled)
	assert.Equal(s.T(), gravitySouthEast, po.Watermark.Gravity)
	assert.Equal(s.T(), 10, po.Watermark.OffsetX)
	assert.Equal(s.T(), 20, po.Watermark.OffsetY)
	assert.Equal(s.T(), 0.6, po.Watermark.Scale)
}

func (s *ProcessingOptionsTestSuite) TestParsePathAdvancedPreset() {
	conf.Presets["test1"] = urlOptions{
		urlOption{Name: "resizing_type", Args: []string{"fill"}},
	}

	conf.Presets["test2"] = urlOptions{
		urlOption{Name: "blur", Args: []string{"0.2"}},
		urlOption{Name: "quality", Args: []string{"50"}},
	}

	req := s.getRequest("http://example.com/http://images.dev/lorem/ipsum.jpg?preset=test1&preset=test2")
	ctx, err := parsePath(context.Background(), req)

	require.Nil(s.T(), err)

	po := getProcessingOptions(ctx)
	assert.Equal(s.T(), resizeFill, po.Resize)
	assert.Equal(s.T(), float32(0.2), po.Blur)
	assert.Equal(s.T(), 50, po.Quality)
}

func (s *ProcessingOptionsTestSuite) TestParsePathPresetDefault() {
	conf.Presets["default"] = urlOptions{
		urlOption{Name: "resizing_type", Args: []string{"fill"}},
		urlOption{Name: "blur", Args: []string{"0.2"}},
		urlOption{Name: "quality", Args: []string{"50"}},
	}

	req := s.getRequest("http://example.com/http://images.dev/lorem/ipsum.jpg?quality=70")
	ctx, err := parsePath(context.Background(), req)

	require.Nil(s.T(), err)

	po := getProcessingOptions(ctx)
	assert.Equal(s.T(), resizeFill, po.Resize)
	assert.Equal(s.T(), float32(0.2), po.Blur)
	assert.Equal(s.T(), 70, po.Quality)
}

func (s *ProcessingOptionsTestSuite) TestParsePathAdvancedPresetLoopDetection() {
	conf.Presets["test1"] = urlOptions{
		urlOption{Name: "resizing_type", Args: []string{"fill"}},
	}

	conf.Presets["test2"] = urlOptions{
		urlOption{Name: "blur", Args: []string{"0.2"}},
		urlOption{Name: "quality", Args: []string{"50"}},
	}

	req := s.getRequest("http://example.com/http://images.dev/lorem/ipsum.jpg?preset=test1:test2:test1")
	ctx, err := parsePath(context.Background(), req)

	require.Nil(s.T(), err)

	po := getProcessingOptions(ctx)
	require.ElementsMatch(s.T(), po.UsedPresets, []string{"test1", "test2"})
}

func (s *ProcessingOptionsTestSuite) TestParsePathAdvancedCachebuster() {
	req := s.getRequest("http://example.com/http://images.dev/lorem/ipsum.jpg?cachebuster=123")
	ctx, err := parsePath(context.Background(), req)

	require.Nil(s.T(), err)

	po := getProcessingOptions(ctx)
	assert.Equal(s.T(), "123", po.CacheBuster)
}

func (s *ProcessingOptionsTestSuite) TestParsePathWebpDetection() {
	conf.EnableWebpDetection = true

	req := s.getRequest("http://example.com/unsafe/plain/http://images.dev/lorem/ipsum.jpg")
	req.Header.Set("Accept", "image/webp")
	ctx, err := parsePath(context.Background(), req)

	require.Nil(s.T(), err)

	po := getProcessingOptions(ctx)
	assert.Equal(s.T(), true, po.PreferWebP)
	assert.Equal(s.T(), false, po.EnforceWebP)
}

func (s *ProcessingOptionsTestSuite) TestParsePathWebpEnforce() {
	conf.EnforceWebp = true

	req := s.getRequest("http://example.com/unsafe/plain/http://images.dev/lorem/ipsum.jpg")
	req.Header.Set("Accept", "image/webp")
	ctx, err := parsePath(context.Background(), req)

	require.Nil(s.T(), err)

	po := getProcessingOptions(ctx)
	assert.Equal(s.T(), true, po.PreferWebP)
	assert.Equal(s.T(), true, po.EnforceWebP)
}

func (s *ProcessingOptionsTestSuite) TestParsePathWidthHeader() {
	conf.EnableClientHints = true

	req := s.getRequest("http://example.com/http://images.dev/lorem/ipsum.jpg")
	req.Header.Set("Width", "100")
	ctx, err := parsePath(context.Background(), req)

	require.Nil(s.T(), err)

	po := getProcessingOptions(ctx)
	assert.Equal(s.T(), 100, po.Width)
}

func (s *ProcessingOptionsTestSuite) TestParsePathWidthHeaderDisabled() {
	req := s.getRequest("http://example.com/http://images.dev/lorem/ipsum.jpg")
	req.Header.Set("Width", "100")
	ctx, err := parsePath(context.Background(), req)

	require.Nil(s.T(), err)

	po := getProcessingOptions(ctx)
	assert.Equal(s.T(), 0, po.Width)
}

func (s *ProcessingOptionsTestSuite) TestParsePathWidthHeaderRedefine() {
	conf.EnableClientHints = true

	req := s.getRequest("http://example.com/http://images.dev/lorem/ipsum.jpg?width=150")
	req.Header.Set("Width", "100")
	ctx, err := parsePath(context.Background(), req)

	require.Nil(s.T(), err)

	po := getProcessingOptions(ctx)
	assert.Equal(s.T(), 150, po.Width)
}

func (s *ProcessingOptionsTestSuite) TestParsePathViewportWidthHeader() {
	conf.EnableClientHints = true

	req := s.getRequest("http://example.com/http://images.dev/lorem/ipsum.jpg")
	req.Header.Set("Viewport-Width", "100")
	ctx, err := parsePath(context.Background(), req)

	require.Nil(s.T(), err)

	po := getProcessingOptions(ctx)
	assert.Equal(s.T(), 100, po.Width)
}

func (s *ProcessingOptionsTestSuite) TestParsePathViewportWidthHeaderDisabled() {
	req := s.getRequest("http://example.com/http://images.dev/lorem/ipsum.jpg")
	req.Header.Set("Viewport-Width", "100")
	ctx, err := parsePath(context.Background(), req)

	require.Nil(s.T(), err)

	po := getProcessingOptions(ctx)
	assert.Equal(s.T(), 0, po.Width)
}

func (s *ProcessingOptionsTestSuite) TestParsePathViewportWidthHeaderRedefine() {
	conf.EnableClientHints = true

	req := s.getRequest("http://example.com/http://images.dev/lorem/ipsum.jpg?width=150")
	req.Header.Set("Viewport-Width", "100")
	ctx, err := parsePath(context.Background(), req)

	require.Nil(s.T(), err)

	po := getProcessingOptions(ctx)
	assert.Equal(s.T(), 150, po.Width)
}

func (s *ProcessingOptionsTestSuite) TestParsePathDprHeader() {
	conf.EnableClientHints = true

	req := s.getRequest("http://example.com/http://images.dev/lorem/ipsum.jpg")
	req.Header.Set("DPR", "2")
	ctx, err := parsePath(context.Background(), req)

	require.Nil(s.T(), err)

	po := getProcessingOptions(ctx)
	assert.Equal(s.T(), 2.0, po.Dpr)
}

func (s *ProcessingOptionsTestSuite) TestParsePathDprHeaderDisabled() {
	req := s.getRequest("http://example.com/http://images.dev/lorem/ipsum.jpg")
	req.Header.Set("DPR", "2")
	ctx, err := parsePath(context.Background(), req)

	require.Nil(s.T(), err)

	po := getProcessingOptions(ctx)
	assert.Equal(s.T(), 1.0, po.Dpr)
}
/*
func (s *ProcessingOptionsTestSuite) TestParsePathSigned() {
	conf.Keys = []securityKey{securityKey("test-key")}
	conf.Salts = []securityKey{securityKey("test-salt")}
	conf.AllowInsecure = false

	req := s.getRequest("http://example.com/HcvNognEV1bW6f8zRqxNYuOkV0IUf1xloRb57CzbT4g/width:150/plain/http://images.dev/lorem/ipsum.jpg@png")
	_, err := parsePath(context.Background(), req)

	require.Nil(s.T(), err)
}

func (s *ProcessingOptionsTestSuite) TestParsePathSignedInvalid() {
	conf.Keys = []securityKey{securityKey("test-key")}
	conf.Salts = []securityKey{securityKey("test-salt")}
	conf.AllowInsecure = false

	req := s.getRequest("http://example.com/unsafe/width:150/plain/http://images.dev/lorem/ipsum.jpg@png")
	_, err := parsePath(context.Background(), req)

	require.Error(s.T(), err)
	assert.Equal(s.T(), errInvalidSignature.Error(), err.Error())
}
*/
func (s *ProcessingOptionsTestSuite) TestParsePathOnlyPresets() {
	conf.OnlyPresets = true
	conf.Presets["test1"] = urlOptions{
		urlOption{Name: "blur", Args: []string{"0.2"}},
	}
	conf.Presets["test2"] = urlOptions{
		urlOption{Name: "quality", Args: []string{"50"}},
	}

	req := s.getRequest("http://example.com/http://images.dev/lorem/ipsum.jpg?preset=test1&preset=test2")

	ctx, err := parsePath(context.Background(), req)

	require.Nil(s.T(), err)

	po := getProcessingOptions(ctx)
	assert.Equal(s.T(), float32(0.2), po.Blur)
	assert.Equal(s.T(), 50, po.Quality)
}

func (s *ProcessingOptionsTestSuite) TestParseBase64URLOnlyPresets() {
	conf.OnlyPresets = true
	conf.Presets["test1"] = urlOptions{
		urlOption{Name: "blur", Args: []string{"0.2"}},
	}
	conf.Presets["test2"] = urlOptions{
		urlOption{Name: "quality", Args: []string{"50"}},
	}

	imageURL := "http://images.dev/lorem/ipsum.jpg?param=value"
	req := s.getRequest(fmt.Sprintf("http://example.com/%s?preset=test1&preset=test2", base64.RawURLEncoding.EncodeToString([]byte(imageURL))))

	ctx, err := parsePath(context.Background(), req)

	require.Nil(s.T(), err)

	po := getProcessingOptions(ctx)
	assert.Equal(s.T(), float32(0.2), po.Blur)
	assert.Equal(s.T(), 50, po.Quality)
}
func TestProcessingOptions(t *testing.T) {
	suite.Run(t, new(ProcessingOptionsTestSuite))
}
