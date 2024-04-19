package vgin

import (
	"context"
	"mime/multipart"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"github.com/gorilla/schema"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
)

const (
	ParamsJsonTag      = "vjson"
	ParamsMultiFormTag = "vform"
	ParamsQueryTag     = "vquery"
	ParamsPathTag      = "vpath"
	ParamsHeaderTag    = "vheader"
)

type paramsInHandler struct {
	into Handler
}

func (ph *paramsInHandler) HandleFunc(ctx context.Context, c *gin.Context) HandleResponse {
	if err := ParseMapParams(c, ph.into); err != nil {
		return &Ret{
			Code:    http.StatusInternalServerError,
			Data:    nil,
			Err:     err,
			Message: "Parse params error",
		}
	}

	return ph.into.HandleFunc(ctx, c)
}

func ParamsIn(handler Handler) Handler {
	return &paramsInHandler{handler}
}

func ParseMapParams(c *gin.Context, into Handler) (err error) {
	switch c.ContentType() {
	case gin.MIMEJSON:
		var raw []byte
		raw, err = c.GetRawData()
		if err != nil {
			return err
		}
		err = parseJson(raw, into)
	case gin.MIMEMultipartPOSTForm:
		var form *multipart.Form
		form, err = c.MultipartForm()
		if err != nil {
			return err
		}
		err = parseMultiForm(form.Value, into)
	}
	if err != nil {
		return errors.Wrap(err, "parse contentType data")
	}

	allwaysParse := []func(*gin.Context, Handler) error{
		parseQuery,
		parsePath,
		parseHeader,
	}

	for _, parser := range allwaysParse {
		if err := parser(c, into); err != nil {
			return err
		}
	}

	return nil
}

func parseHeader(c *gin.Context, intoHandler Handler) error {
	if len(c.Request.Header) == 0 {
		return nil
	}

	headerMap := c.Request.Header
	if err := schemaDecoder(ParamsHeaderTag).Decode(intoHandler, headerMap); err != nil {
		return errors.Wrap(err, "parseHeader")
	}

	return nil
}

func parsePath(c *gin.Context, intoHandler Handler) error {
	if len(c.Params) == 0 {
		return nil
	}

	tmp := make(map[string][]string)
	for _, p := range c.Params {
		tmp[p.Key] = []string{p.Value}
	}

	if err := schemaDecoder(ParamsPathTag).Decode(intoHandler, tmp); err != nil {
		return errors.Wrap(err, "parsePath")
	}

	return nil
}

func parseQuery(c *gin.Context, intoHandler Handler) error {
	queryMap := c.Request.URL.Query()
	if len(queryMap) == 0 {
		return nil
	}

	if err := schemaDecoder(ParamsQueryTag).Decode(intoHandler, queryMap); err != nil {
		return errors.Wrap(err, "parseQuert")
	}

	return nil
}

func schemaDecoder(tag string) *schema.Decoder {
	decoder := schema.NewDecoder()
	decoder.SetAliasTag(tag)
	decoder.IgnoreUnknownKeys(true)
	return decoder
}

func mapstructureDecoder(tag string, result Handler) (*mapstructure.Decoder, error) {
	config := &mapstructure.DecoderConfig{
		Result:  result,
		TagName: tag,
	}

	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return nil, err
	}

	return decoder, nil
}

func parseJson(data []byte, intoHandler Handler) error {
	tmp := make(map[string]any)
	if err := json.Unmarshal(data, &tmp); err != nil {
		return errors.Wrap(err, "parseJson")
	}
	decoder, err := mapstructureDecoder(ParamsJsonTag, intoHandler)
	if err != nil {
		return errors.Wrap(err, "parseJson")
	}

	if err := decoder.Decode(tmp); err != nil {
		errors.Wrap(err, "parseJson")
	}

	return nil
}

func parseMultiForm(data map[string][]string, intoHandler Handler) error {
	return schemaDecoder(ParamsMultiFormTag).Decode(intoHandler, data)
}
