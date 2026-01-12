// Package slogutil provides slog utilities using the standard library.
package slogutil

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"sync"
)

// IndentedJSONHandler is a slog.Handler that outputs indented JSON.
type IndentedJSONHandler struct {
	opts   slog.HandlerOptions
	w      io.Writer
	mu     *sync.Mutex
	attrs  []slog.Attr
	groups []string
}

// NewIndentedJSONHandler creates a new handler that outputs indented JSON.
func NewIndentedJSONHandler(w io.Writer, opts *slog.HandlerOptions) *IndentedJSONHandler {
	if opts == nil {
		opts = &slog.HandlerOptions{}
	}
	return &IndentedJSONHandler{
		opts: *opts,
		w:    w,
		mu:   &sync.Mutex{},
	}
}

func (h *IndentedJSONHandler) Enabled(_ context.Context, level slog.Level) bool {
	minLevel := slog.LevelInfo
	if h.opts.Level != nil {
		minLevel = h.opts.Level.Level()
	}
	return level >= minLevel
}

func (h *IndentedJSONHandler) Handle(_ context.Context, r slog.Record) error {
	// Build the log entry map
	entry := make(map[string]any)

	entry["time"] = r.Time.Format("2006-01-02T15:04:05.000Z07:00")
	entry["level"] = r.Level.String()
	entry["msg"] = r.Message

	if h.opts.AddSource && r.PC != 0 {
		// Add source info if requested
		fs := r.Source()
		if fs != nil {
			entry["source"] = map[string]any{
				"function": fs.Function,
				"file":     fs.File,
				"line":     fs.Line,
			}
		}
	}

	// Add handler-level attrs
	current := entry
	for _, g := range h.groups {
		nested := make(map[string]any)
		current[g] = nested
		current = nested
	}
	for _, a := range h.attrs {
		addAttr(current, a, h.opts.ReplaceAttr)
	}

	// Add record attrs
	r.Attrs(func(a slog.Attr) bool {
		addAttr(current, a, h.opts.ReplaceAttr)
		return true
	})

	// Marshal with indentation
	data, err := json.MarshalIndent(entry, "", "  ")
	if err != nil {
		return err
	}

	h.mu.Lock()
	defer h.mu.Unlock()
	_, err = h.w.Write(append(data, '\n'))
	return err
}

func (h *IndentedJSONHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	newAttrs := make([]slog.Attr, len(h.attrs)+len(attrs))
	copy(newAttrs, h.attrs)
	copy(newAttrs[len(h.attrs):], attrs)
	return &IndentedJSONHandler{
		opts:   h.opts,
		w:      h.w,
		mu:     h.mu,
		attrs:  newAttrs,
		groups: h.groups,
	}
}

func (h *IndentedJSONHandler) WithGroup(name string) slog.Handler {
	if name == "" {
		return h
	}
	newGroups := make([]string, len(h.groups)+1)
	copy(newGroups, h.groups)
	newGroups[len(h.groups)] = name
	return &IndentedJSONHandler{
		opts:   h.opts,
		w:      h.w,
		mu:     h.mu,
		attrs:  h.attrs,
		groups: newGroups,
	}
}

func addAttr(m map[string]any, a slog.Attr, replace func([]string, slog.Attr) slog.Attr) {
	if replace != nil {
		a = replace(nil, a)
	}
	if a.Key == "" {
		return
	}

	v := a.Value.Resolve()

	switch v.Kind() {
	case slog.KindGroup:
		attrs := v.Group()
		if len(attrs) == 0 {
			return
		}
		nested := make(map[string]any)
		for _, ga := range attrs {
			addAttr(nested, ga, replace)
		}
		m[a.Key] = nested
	default:
		m[a.Key] = v.Any()
	}
}
