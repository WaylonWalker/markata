// Package renderprotocol exposes a transport-neutral request/response envelope
// for the standalone Go renderer. The stdio adapter uses newline-delimited JSON.
package renderprotocol

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/WaylonWalker/markata/pkg/plugins"
)

const (
	// MaxFrameBytes bounds one newline-delimited request, including JSON framing.
	MaxFrameBytes = 1 << 20
	// MaxSourceBytes matches the local editor's limit for individual documents.
	MaxSourceBytes = 200_000
	MaxRequestIDBytes = 128
)

// Options uses JSON field names independent of Go struct-field matching.
type Options struct {
	Unsafe       bool   `json:"unsafe,omitempty"`
	InlineAssets bool   `json:"inline_assets,omitempty"`
	AssetBase    string `json:"asset_base,omitempty"`
}

// Request is the v1 request envelope shared by future protocol adapters.
type Request struct {
	ID      string  `json:"id"`
	Op      string  `json:"op"`
	Source  string  `json:"source,omitempty"`
	Options Options `json:"options,omitempty"`
}

// ProtocolError is a stable machine-readable error for a single request.
type ProtocolError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// Response is exactly one response to one nonempty input line.
type Response struct {
	ID           string                          `json:"id,omitempty"`
	OK           bool                            `json:"ok"`
	Result       *plugins.StandaloneResult       `json:"result,omitempty"`
	Capabilities *plugins.StandaloneCapabilities `json:"capabilities,omitempty"`
	Error        *ProtocolError                  `json:"error,omitempty"`
}

func failure(id, code, message string) Response {
	return Response{ID: id, OK: false, Error: &ProtocolError{Code: code, Message: message}}
}

// Handle is transport-neutral. It does not retain per-request configuration.
func Handle(req Request) Response {
	if req.ID == "" || len(req.ID) > MaxRequestIDBytes {
		return failure("", "invalid_id", "id must contain 1 to 128 UTF-8 bytes")
	}
	switch req.Op {
	case "capabilities":
		caps := plugins.GetStandaloneCapabilities()
		return Response{ID: req.ID, OK: true, Capabilities: &caps}
	case "render":
		if len(req.Source) > MaxSourceBytes {
			return failure(req.ID, "source_too_large", fmt.Sprintf("source exceeds %d bytes", MaxSourceBytes))
		}
		result, err := plugins.RenderStandalone(req.Source, plugins.StandaloneOptions{
			Unsafe: req.Options.Unsafe, InlineAssets: req.Options.InlineAssets,
			AssetBase: req.Options.AssetBase,
		})
		if err != nil {
			return failure(req.ID, "render_failed", err.Error())
		}
		return Response{ID: req.ID, OK: true, Result: result}
	case "shutdown":
		return Response{ID: req.ID, OK: true}
	default:
		return failure(req.ID, "unsupported_operation", fmt.Sprintf("unsupported operation %q", req.Op))
	}
}

// session retains parsers for one stdio connection. Serve handles requests
// sequentially: StandaloneRenderer is reusable but not concurrency-safe.
type session struct {
	renderers map[Options]*plugins.StandaloneRenderer
}

func (s *session) handle(req Request) Response {
	if req.Op != "render" || req.ID == "" || len(req.ID) > MaxRequestIDBytes ||
		len(req.Source) > MaxSourceBytes {
		return Handle(req)
	}
	if s.renderers == nil {
		s.renderers = make(map[Options]*plugins.StandaloneRenderer)
	}
	renderer := s.renderers[req.Options]
	if renderer == nil {
		// Bound renderer variants when clients choose arbitrary asset bases.
		if len(s.renderers) >= 8 {
			s.renderers = make(map[Options]*plugins.StandaloneRenderer)
		}
		renderer = plugins.NewStandaloneRenderer(plugins.StandaloneOptions{
			Unsafe: req.Options.Unsafe, InlineAssets: req.Options.InlineAssets,
			AssetBase: req.Options.AssetBase,
		})
		s.renderers[req.Options] = renderer
	}
	result, err := renderer.Render(req.Source)
	if err != nil {
		return failure(req.ID, "render_failed", err.Error())
	}
	return Response{ID: req.ID, OK: true, Result: result}
}

// DecodeRequest rejects unexpected fields and multiple JSON values in a line.
func DecodeRequest(line []byte) (Request, error) {
	var req Request
	decoder := json.NewDecoder(bytes.NewReader(line))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&req); err != nil {
		return Request{}, err
	}
	var extra any
	if err := decoder.Decode(&extra); !errors.Is(err, io.EOF) {
		if err == nil {
			return Request{}, errors.New("multiple JSON values in one request")
		}
		return Request{}, err
	}
	return req, nil
}

// Serve processes many requests with one running Go renderer. It emits only
// JSON response lines to out. A malformed request is isolated from later lines.
// An oversized frame is consumed without buffering it fully in memory.
func Serve(in io.Reader, out io.Writer) error {
	reader := bufio.NewReaderSize(in, 64*1024)
	encoder := json.NewEncoder(out)
	state := &session{}
	for {
		line, oversized, err := readFrame(reader)
		if errors.Is(err, io.EOF) {
			return nil
		}
		if err != nil {
			return fmt.Errorf("read stdio request: %w", err)
		}
		if oversized {
			if err := encoder.Encode(failure("", "frame_too_large",
				fmt.Sprintf("request exceeds %d bytes", MaxFrameBytes))); err != nil {
				return err
			}
			continue
		}
		if len(bytes.TrimSpace(line)) == 0 {
			continue
		}
		req, err := DecodeRequest(line)
		if err != nil {
			if err := encoder.Encode(failure("", "invalid_json", err.Error())); err != nil {
				return err
			}
			continue
		}
		resp := state.handle(req)
		if err := encoder.Encode(resp); err != nil {
			return err
		}
		if req.Op == "shutdown" && resp.OK {
			return nil
		}
	}
}

// readFrame retains at most MaxFrameBytes while draining any oversized input
// through its newline, so one bad frame does not desynchronize the stream.
func readFrame(reader *bufio.Reader) ([]byte, bool, error) {
	var line []byte
	oversized := false
	for {
		part, err := reader.ReadSlice('\n')
		if !oversized && len(line)+len(part) > MaxFrameBytes {
			oversized = true
			line = nil
		}
		if !oversized {
			line = append(line, part...)
		}
		switch {
		case err == nil:
			return bytes.TrimSuffix(line, []byte{'\n'}), oversized, nil
		case errors.Is(err, bufio.ErrBufferFull):
			continue
		case errors.Is(err, io.EOF):
			if !oversized && len(line) == 0 {
				return nil, false, io.EOF
			}
			return line, oversized, nil
		default:
			return nil, false, err
		}
	}
}
