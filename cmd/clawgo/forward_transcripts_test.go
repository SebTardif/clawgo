package main

import (
	"bufio"
	"context"
	"encoding/json"
	"net"
	"testing"
	"time"

	"github.com/clawdbot/clawgo/internal/routing"
	"github.com/clawdbot/clawgo/modules/stt"
)

func TestForwardTranscriptsSkipsPartials(t *testing.T) {
	for _, tc := range []struct {
		name                 string
		agent, router, quick bool
	}{
		{name: "voice"}, {name: "agent", agent: true},
		{name: "voice-router", router: true}, {name: "agent-router", router: true, agent: true},
		{name: "quick-action", router: true, quick: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			writer, reader := net.Pipe()
			defer reader.Close()
			defer writer.Close()
			if err := reader.SetReadDeadline(time.Now().Add(5 * time.Second)); err != nil {
				t.Fatal(err)
			}
			client := &BridgeClient{conn: writer, logf: func(string, ...any) {}}
			cfg := NodeConfig{SessionKey: "synthetic-session", AgentRequest: tc.agent}
			var router routing.Router
			if tc.router {
				var err error
				router, err = routing.New("default", routing.Config{SessionKey: cfg.SessionKey, AgentRequest: tc.agent, QuickActions: tc.quick, DeliverChannel: "telegram", DeliverTo: "synthetic-destination"}, bridgeTransport{client: client}, nil)
				if err != nil {
					t.Fatal(err)
				}
			}
			text := "finished"
			if tc.quick {
				text = "telegram ping"
			}
			in := make(chan stt.Transcript, 4)
			in <- stt.Transcript{Text: text, Final: false}
			in <- stt.Transcript{Text: " ", Final: true}
			in <- stt.Transcript{Text: "", Final: false}
			in <- stt.Transcript{Text: "  " + text + "  ", Final: true}
			close(in)
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			go func() { forwardTranscripts(ctx, client, cfg, in, router); writer.Close() }()
			scanner := bufio.NewScanner(reader)
			var frames []map[string]any
			for scanner.Scan() {
				var frame map[string]any
				if err := json.Unmarshal(scanner.Bytes(), &frame); err != nil {
					t.Fatal(err)
				}
				frames = append(frames, frame)
			}
			if err := scanner.Err(); err != nil {
				t.Fatal(err)
			}
			if len(frames) != 1 {
				t.Fatalf("forwarded %d frames, want exactly one final utterance", len(frames))
			}
			frame := frames[0]
			if tc.quick {
				if frame["method"] != "send" {
					t.Fatalf("expected one quick action, got %v", frame)
				}
				return
			}
			wantEvent, field := "voice.transcript", "text"
			if tc.agent {
				wantEvent, field = "agent.request", "message"
			}
			if frame["event"] != wantEvent {
				t.Fatalf("event=%v, want %s", frame["event"], wantEvent)
			}
			var payload map[string]any
			if err := json.Unmarshal([]byte(frame["payloadJSON"].(string)), &payload); err != nil {
				t.Fatal(err)
			}
			if payload[field] != text || payload["sessionKey"] != cfg.SessionKey {
				t.Fatalf("unexpected payload %v", payload)
			}
		})
	}
}
