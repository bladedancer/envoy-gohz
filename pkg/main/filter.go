package main

import (
	"context"
	"fmt"
	"strconv"

	"github.com/envoyproxy/envoy/contrib/golang/common/go/api"
)

var UpdateUpstreamBody = "upstream response body updated by the simple plugin"

type filter struct {
	api.PassThroughStreamFilter

	callbacks api.FilterCallbackHandler
	path      string
	config    *config
}

func (f *filter) sendLocalReplyInternal(body string, status int) api.StatusType {
	headers := make(map[string][]string)
	f.callbacks.DecoderFilterCallbacks().SendLocalReply(status, body, headers, -1, "response-from-go")
	return api.LocalReply
}

func (f *filter) DecodeHeaders(header api.RequestHeaderMap, endStream bool) api.StatusType {
	queue, err := f.config.hazelcastClient.GetQueue(context.Background(), "someQueue")
	if err != nil {
		return f.sendLocalReplyInternal(err.Error(), 500)
	}

	// Add something to the queue
	queue.Add(context.Background(), f.path)

	f.path, _ = header.Get(":path")
	if f.path == "/localreply_by_config" {
		body := fmt.Sprintf("%s, path: %s\r\n", f.config.echoBody, f.path)
		return f.sendLocalReplyInternal(body, 200)
	}
	return api.Continue
}

func (f *filter) DecodeData(buffer api.BufferInstance, endStream bool) api.StatusType {
	return api.Continue
}

func (f *filter) DecodeTrailers(trailers api.RequestTrailerMap) api.StatusType {
	return api.Continue
}

func (f *filter) EncodeHeaders(header api.ResponseHeaderMap, endStream bool) api.StatusType {
	if f.path == "/update_upstream_response" {
		header.Set("Content-Length", strconv.Itoa(len(UpdateUpstreamBody)))
	}
	header.Set("Rsp-Header-From-Go", "bar-test")
	return api.Continue
}

func (f *filter) EncodeData(buffer api.BufferInstance, endStream bool) api.StatusType {
	if f.path == "/update_upstream_response" {
		if endStream {
			buffer.SetString(UpdateUpstreamBody)
		} else {
			// TODO implement buffer->Drain, buffer.SetString means buffer->Drain(buffer.Len())
			buffer.SetString("")
		}
	}
	return api.Continue
}

func (f *filter) EncodeTrailers(trailers api.ResponseTrailerMap) api.StatusType {
	return api.Continue
}

func (f *filter) OnDestroy(reason api.DestroyReason) {
}
