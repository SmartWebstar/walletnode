// Code generated by goa v3.3.1, DO NOT EDIT.
//
// artworks client HTTP transport
//
// Command:
// $ goa gen github.com/pastelnetwork/walletnode/api/design

package client

import (
	"context"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	artworks "github.com/pastelnetwork/walletnode/api/gen/artworks"
	goahttp "goa.design/goa/v3/http"
	goa "goa.design/goa/v3/pkg"
)

// Client lists the artworks service endpoint HTTP clients.
type Client struct {
	// Register Doer is the HTTP client used to make requests to the register
	// endpoint.
	RegisterDoer goahttp.Doer

	// RegisterJobState Doer is the HTTP client used to make requests to the
	// registerJobState endpoint.
	RegisterJobStateDoer goahttp.Doer

	// RegisterJob Doer is the HTTP client used to make requests to the registerJob
	// endpoint.
	RegisterJobDoer goahttp.Doer

	// RegisterJobs Doer is the HTTP client used to make requests to the
	// registerJobs endpoint.
	RegisterJobsDoer goahttp.Doer

	// UploadImage Doer is the HTTP client used to make requests to the uploadImage
	// endpoint.
	UploadImageDoer goahttp.Doer

	// CORS Doer is the HTTP client used to make requests to the  endpoint.
	CORSDoer goahttp.Doer

	// RestoreResponseBody controls whether the response bodies are reset after
	// decoding so they can be read again.
	RestoreResponseBody bool

	scheme     string
	host       string
	encoder    func(*http.Request) goahttp.Encoder
	decoder    func(*http.Response) goahttp.Decoder
	dialer     goahttp.Dialer
	configurer *ConnConfigurer
}

// ArtworksUploadImageEncoderFunc is the type to encode multipart request for
// the "artworks" service "uploadImage" endpoint.
type ArtworksUploadImageEncoderFunc func(*multipart.Writer, *artworks.UploadImagePayload) error

// NewClient instantiates HTTP clients for all the artworks service servers.
func NewClient(
	scheme string,
	host string,
	doer goahttp.Doer,
	enc func(*http.Request) goahttp.Encoder,
	dec func(*http.Response) goahttp.Decoder,
	restoreBody bool,
	dialer goahttp.Dialer,
	cfn *ConnConfigurer,
) *Client {
	if cfn == nil {
		cfn = &ConnConfigurer{}
	}
	return &Client{
		RegisterDoer:         doer,
		RegisterJobStateDoer: doer,
		RegisterJobDoer:      doer,
		RegisterJobsDoer:     doer,
		UploadImageDoer:      doer,
		CORSDoer:             doer,
		RestoreResponseBody:  restoreBody,
		scheme:               scheme,
		host:                 host,
		decoder:              dec,
		encoder:              enc,
		dialer:               dialer,
		configurer:           cfn,
	}
}

// Register returns an endpoint that makes HTTP requests to the artworks
// service register server.
func (c *Client) Register() goa.Endpoint {
	var (
		encodeRequest  = EncodeRegisterRequest(c.encoder)
		decodeResponse = DecodeRegisterResponse(c.decoder, c.RestoreResponseBody)
	)
	return func(ctx context.Context, v interface{}) (interface{}, error) {
		req, err := c.BuildRegisterRequest(ctx, v)
		if err != nil {
			return nil, err
		}
		err = encodeRequest(req, v)
		if err != nil {
			return nil, err
		}
		resp, err := c.RegisterDoer.Do(req)
		if err != nil {
			return nil, goahttp.ErrRequestError("artworks", "register", err)
		}
		return decodeResponse(resp)
	}
}

// RegisterJobState returns an endpoint that makes HTTP requests to the
// artworks service registerJobState server.
func (c *Client) RegisterJobState() goa.Endpoint {
	var (
		decodeResponse = DecodeRegisterJobStateResponse(c.decoder, c.RestoreResponseBody)
	)
	return func(ctx context.Context, v interface{}) (interface{}, error) {
		req, err := c.BuildRegisterJobStateRequest(ctx, v)
		if err != nil {
			return nil, err
		}
		var cancel context.CancelFunc
		ctx, cancel = context.WithCancel(ctx)
		conn, resp, err := c.dialer.DialContext(ctx, req.URL.String(), req.Header)
		if err != nil {
			if resp != nil {
				return decodeResponse(resp)
			}
			return nil, goahttp.ErrRequestError("artworks", "registerJobState", err)
		}
		if c.configurer.RegisterJobStateFn != nil {
			conn = c.configurer.RegisterJobStateFn(conn, cancel)
		}
		go func() {
			<-ctx.Done()
			conn.WriteControl(
				websocket.CloseMessage,
				websocket.FormatCloseMessage(websocket.CloseNormalClosure, "client closing connection"),
				time.Now().Add(time.Second),
			)
			conn.Close()
		}()
		stream := &RegisterJobStateClientStream{conn: conn}
		return stream, nil
	}
}

// RegisterJob returns an endpoint that makes HTTP requests to the artworks
// service registerJob server.
func (c *Client) RegisterJob() goa.Endpoint {
	var (
		decodeResponse = DecodeRegisterJobResponse(c.decoder, c.RestoreResponseBody)
	)
	return func(ctx context.Context, v interface{}) (interface{}, error) {
		req, err := c.BuildRegisterJobRequest(ctx, v)
		if err != nil {
			return nil, err
		}
		resp, err := c.RegisterJobDoer.Do(req)
		if err != nil {
			return nil, goahttp.ErrRequestError("artworks", "registerJob", err)
		}
		return decodeResponse(resp)
	}
}

// RegisterJobs returns an endpoint that makes HTTP requests to the artworks
// service registerJobs server.
func (c *Client) RegisterJobs() goa.Endpoint {
	var (
		decodeResponse = DecodeRegisterJobsResponse(c.decoder, c.RestoreResponseBody)
	)
	return func(ctx context.Context, v interface{}) (interface{}, error) {
		req, err := c.BuildRegisterJobsRequest(ctx, v)
		if err != nil {
			return nil, err
		}
		resp, err := c.RegisterJobsDoer.Do(req)
		if err != nil {
			return nil, goahttp.ErrRequestError("artworks", "registerJobs", err)
		}
		return decodeResponse(resp)
	}
}

// UploadImage returns an endpoint that makes HTTP requests to the artworks
// service uploadImage server.
func (c *Client) UploadImage(artworksUploadImageEncoderFn ArtworksUploadImageEncoderFunc) goa.Endpoint {
	var (
		encodeRequest  = EncodeUploadImageRequest(NewArtworksUploadImageEncoder(artworksUploadImageEncoderFn))
		decodeResponse = DecodeUploadImageResponse(c.decoder, c.RestoreResponseBody)
	)
	return func(ctx context.Context, v interface{}) (interface{}, error) {
		req, err := c.BuildUploadImageRequest(ctx, v)
		if err != nil {
			return nil, err
		}
		err = encodeRequest(req, v)
		if err != nil {
			return nil, err
		}
		resp, err := c.UploadImageDoer.Do(req)
		if err != nil {
			return nil, goahttp.ErrRequestError("artworks", "uploadImage", err)
		}
		return decodeResponse(resp)
	}
}