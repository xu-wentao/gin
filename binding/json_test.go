// Copyright 2019 Gin Core Team. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package binding

import (
	stdjson "encoding/json"
	"io"
	"net/http/httptest"
	"testing"

	jsoncodec "github.com/gin-gonic/gin/codec/json"
	"github.com/gin-gonic/gin/render"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJSONBindingBindBody(t *testing.T) {
	var s struct {
		Foo string `json:"foo"`
	}
	err := jsonBinding{}.BindBody([]byte(`{"foo": "FOO"}`), &s)
	require.NoError(t, err)
	assert.Equal(t, "FOO", s.Foo)
}

func TestJSONBindingBindBodyMap(t *testing.T) {
	s := make(map[string]string)
	err := jsonBinding{}.BindBody([]byte(`{"foo": "FOO","hello":"world"}`), &s)
	require.NoError(t, err)
	assert.Len(t, s, 2)
	assert.Equal(t, "FOO", s["foo"])
	assert.Equal(t, "world", s["hello"])
}

func TestCustomJSONCodec(t *testing.T) {
	previous := jsoncodec.API
	custom := &trackingJSONAPI{}
	jsoncodec.API = custom
	t.Cleanup(func() { jsoncodec.API = previous })

	var data struct {
		Foo string `json:"foo"`
	}
	require.NoError(t, jsonBinding{}.BindBody([]byte(`{"foo":"bar"}`), &data))
	assert.Equal(t, "bar", data.Foo)
	assert.Equal(t, 1, custom.decoderCalls)

	w := httptest.NewRecorder()
	require.NoError(t, (render.PureJSON{Data: data}).Render(w))
	assert.JSONEq(t, `{"foo":"bar"}`, w.Body.String())
	assert.Equal(t, 1, custom.encoderCalls)
}

type trackingJSONAPI struct {
	encoderCalls int
	decoderCalls int
}

func (j *trackingJSONAPI) Marshal(v any) ([]byte, error) {
	return stdjson.Marshal(v)
}

func (j *trackingJSONAPI) Unmarshal(data []byte, v any) error {
	return stdjson.Unmarshal(data, v)
}

func (*trackingJSONAPI) MarshalIndent(v any, prefix, indent string) ([]byte, error) {
	return stdjson.MarshalIndent(v, prefix, indent)
}

func (j *trackingJSONAPI) NewEncoder(writer io.Writer) jsoncodec.Encoder {
	j.encoderCalls++
	return stdjson.NewEncoder(writer)
}

func (j *trackingJSONAPI) NewDecoder(reader io.Reader) jsoncodec.Decoder {
	j.decoderCalls++
	return stdjson.NewDecoder(reader)
}
