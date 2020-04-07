package applog

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

type Severity string

const (
	Default   Severity = "DEFAULT"  // (0) The log entry has no assigned severity level.
	Debug     Severity = "DEBUG"    // (100) Debug or trace information.
	Info      Severity = "INFO"     // (200) Routine information, such as ongoing status or performance.
	Notice    Severity = "NOTICE"   // (300) Normal but significant events, such as start up, shut down, or a configuration change.
	Warning   Severity = "WARNING"  // (400) Warning events might cause problems.
	Error     Severity = "ERROR"    // (500) Error events are likely to cause problems.
	Critical  Severity = "CRITICAL" // (600) Critical events cause more severe problems or outages.
	Alert     Severity = "ALERT"    // (700) A person must take an action immediately.
	Emergency Severity = "EMERGENCY"
)

const keyTraceContext = "trace-context"

func TraceContextMiddleware(projectID string) mux.MiddlewareFunc {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var trace string
			if projectID != "" {
				traceHeader := r.Header.Get("X-Cloud-Trace-Context")
				traceParts := strings.Split(traceHeader, "/")
				if len(traceParts) > 0 && len(traceParts[0]) > 0 {
					trace = fmt.Sprintf("projects/%s/traces/%s", projectID, traceParts[0])
				}
			}

			handler.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), keyTraceContext, trace)))
		})
	}
}

func (s Severity) Log(ctx context.Context, component, message string) {
	trace, _ := ctx.Value(keyTraceContext).(string)
	log.Println(Entry{
		Severity:  string(s),
		Message:   message,
		Component: component,
		Trace:     trace,
	})
}

// Entry defines a log entry.
type Entry struct {
	Message  string `json:"message"`
	Severity string `json:"severity,omitempty"`
	Trace    string `json:"logging.googleapis.com/trace,omitempty"`

	// Stackdriver Log Viewer allows filtering and display of this as `jsonPayload.component`.
	Component string `json:"component,omitempty"`
}

// String renders an entry structure to the JSON format expected by Stackdriver.
func (e Entry) String() string {
	if e.Severity == "" {
		e.Severity = "INFO"
	}
	out, err := json.Marshal(e)
	if err != nil {
		log.Printf("json.Marshal: %v", err)
	}
	return string(out)
}
