package github

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/google/go-github/v93/github"
)

func newSBOMTestClient(t *testing.T, handler http.HandlerFunc) *github.Client {
	t.Helper()
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)
	client, err := github.NewClient(
		github.WithHTTPClient(server.Client()),
		github.WithEnterpriseURLs(server.URL+"/api/v3/", server.URL+"/api/uploads/"),
		github.WithAuthToken("test-token"),
	)
	if err != nil {
		t.Fatal(err)
	}
	return client
}

func TestGetRepositorySBOM(t *testing.T) {
	var generations, fetches, downloads atomic.Int32
	download := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		downloads.Add(1)
		if r.Header.Get("Authorization") != "" {
			t.Error("download must not contain GitHub authentication")
		}
		fmt.Fprint(w, `{"SPDXID":"SPDXRef-DOCUMENT","spdxVersion":"SPDX-2.3","name":"owner/repo","packages":[{"name":"example","versionInfo":"1.0"}]}`)
	}))
	defer download.Close()
	client := newSBOMTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.Header.Get("Authorization") != "Bearer test-token" {
			t.Errorf("unexpected API request: %s, authenticated=%t", r.Method, r.Header.Get("Authorization") != "")
		}
		switch r.URL.Path {
		case "/api/v3/repos/owner/repo/dependency-graph/sbom/generate-report":
			generations.Add(1)
			w.WriteHeader(http.StatusCreated)
			fmt.Fprint(w, `{"sbom_url":"https://api.github.com/repos/owner/repo/dependency-graph/sbom/fetch-report/report-id"}`)
		case "/api/v3/repos/owner/repo/dependency-graph/sbom/fetch-report/report-id":
			if fetches.Add(1) <= 2 {
				w.WriteHeader(http.StatusAccepted)
				return
			}
			http.Redirect(w, r, download.URL, http.StatusFound)
		default:
			t.Errorf("unexpected endpoint: %s", r.URL.Path)
			http.NotFound(w, r)
		}
	})

	sbom, err := getRepositorySBOM(context.Background(), client, "owner", "repo", time.Millisecond)
	if err != nil {
		t.Fatal(err)
	}
	if sbom.SBOM.GetSPDXID() != "SPDXRef-DOCUMENT" || sbom.SBOM.GetName() != "owner/repo" || len(sbom.SBOM.Packages) != 1 {
		t.Fatalf("unexpected SBOM: %v", sbom)
	}
	if generations.Load() != 1 || fetches.Load() != 3 || downloads.Load() != 1 {
		t.Fatalf("requests: generation=%d, fetch=%d, download=%d", generations.Load(), fetches.Load(), downloads.Load())
	}
}

func TestGetRepositorySBOMErrors(t *testing.T) {
	for _, tt := range []struct {
		name           string
		generation     string
		genStatus      int
		fetchStatus    int
		downloadStatus int
		document       string
		want           string
	}{
		{name: "generation forbidden", genStatus: 403, want: "403"},
		{name: "fetch not found", fetchStatus: 404, want: "404"},
		{name: "missing report URL", generation: `{}`, want: "no valid report identifier"},
		{name: "missing report identifier", generation: `{"sbom_url":"https://api.github.com/fetch-report/"}`, want: "no valid report identifier"},
		{name: "invalid report URL", generation: `{"sbom_url":"%"}`, want: "invalid SBOM report URL"},
		{name: "download failure", downloadStatus: 502, want: "502"},
		{name: "invalid document", document: `{`, want: "unexpected EOF"},
		{name: "empty document", document: `null`, want: "empty document"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			download := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if tt.downloadStatus != 0 {
					w.WriteHeader(tt.downloadStatus)
				}
				fmt.Fprint(w, tt.document)
			}))
			defer download.Close()
			client := newSBOMTestClient(t, func(w http.ResponseWriter, r *http.Request) {
				if strings.HasSuffix(r.URL.Path, "/generate-report") {
					if tt.genStatus != 0 {
						w.WriteHeader(tt.genStatus)
						return
					}
					w.WriteHeader(http.StatusCreated)
					if tt.generation != "" {
						fmt.Fprint(w, tt.generation)
					} else {
						fmt.Fprint(w, `{"sbom_url":"https://api.github.com/fetch-report/report-id"}`)
					}
					return
				}
				if tt.fetchStatus != 0 {
					w.WriteHeader(tt.fetchStatus)
					return
				}
				http.Redirect(w, r, download.URL, http.StatusFound)
			})
			sbom, err := getRepositorySBOM(context.Background(), client, "owner", "repo", time.Millisecond)
			if sbom != nil || err == nil || !strings.Contains(err.Error(), tt.want) {
				t.Fatalf("got (%v, %v), want an error containing %q", sbom, err, tt.want)
			}
		})
	}
}

func TestGetRepositorySBOMCancellation(t *testing.T) {
	for _, timeout := range []bool{false, true} {
		t.Run(fmt.Sprintf("timeout=%t", timeout), func(t *testing.T) {
			var fetches atomic.Int32
			var ctx context.Context
			var cancel context.CancelFunc
			want := context.Canceled
			if timeout {
				ctx, cancel = context.WithTimeout(context.Background(), time.Second)
				want = context.DeadlineExceeded
			} else {
				ctx, cancel = context.WithCancel(context.Background())
			}
			defer cancel()
			client := newSBOMTestClient(t, func(w http.ResponseWriter, r *http.Request) {
				if strings.HasSuffix(r.URL.Path, "/generate-report") {
					w.WriteHeader(http.StatusCreated)
					fmt.Fprint(w, `{"sbom_url":"https://api.github.com/fetch-report/report-id"}`)
					return
				}
				fetches.Add(1)
				w.WriteHeader(http.StatusAccepted)
				if !timeout {
					cancel()
				}
			})
			sbom, err := getRepositorySBOM(ctx, client, "owner", "repo", time.Hour)
			if sbom != nil || !errors.Is(err, want) {
				t.Fatalf("got (%v, %v), want %v", sbom, err, want)
			}
			if fetches.Load() != 1 {
				t.Fatalf("got %d fetches, want 1", fetches.Load())
			}
		})
	}
}
