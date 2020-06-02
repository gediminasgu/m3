package httpd

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/uber-go/tally"
)

func TestHanlersSwitch(t *testing.T) {
	tests := []struct {
		enabled      bool
		query        string
		expectPromql bool
		expectM3Q    bool
	}{
		{enabled: true, query: "/query?query=sum(metric)", expectPromql: false, expectM3Q: true},
		{enabled: false, query: "/query?query=sum(metric)", expectPromql: true, expectM3Q: false},
		{enabled: true, query: "/query?query=sum(metric)&engine=prometheus", expectPromql: true, expectM3Q: false},
		{enabled: true, query: "/query?query=random(metric)", expectPromql: true, expectM3Q: false},
		{enabled: true, query: "/query?query=random(metric)&engine=m3query", expectPromql: false, expectM3Q: true},
	}

	for _, tt := range tests {
		t.Run(tt.query, func(t *testing.T) {
			promqlCalled := false
			promqlHandler := func(w http.ResponseWriter, req *http.Request) {
				promqlCalled = true
			}

			m3qCalled := false
			m3qHandler := func(w http.ResponseWriter, req *http.Request) {
				m3qCalled = true
			}

			scope := tally.NewTestScope("", nil)
			router := NewRouter(RouterOptions{
				Enabled:            tt.enabled,
				DefaultQueryEngine: "prometheus",
				PromqlHandler:      promqlHandler,
				M3QueryHandler:     m3qHandler,
				LogStats:           true,
				Scope:              scope,
			})

			req, err := http.NewRequest("GET", tt.query, nil)
			if err != nil {
				t.Fatal(err)
			}
			rr := httptest.NewRecorder()

			router.ServeHTTP(rr, req)

			assert.True(t, promqlCalled == tt.expectPromql)
			assert.True(t, m3qCalled == tt.expectM3Q)
		})
	}
}

func TestM3Rounting(t *testing.T) {
	tests := []struct {
		query string
		pass  bool
	}{
		{query: `sum(cpu_total)`, pass: true},
		{query: `sum(cpu_total{})`, pass: true},
		{query: `sum(cpu_total{par1="111-asd"})`, pass: true},
		{query: `sum(cpu_total{par1="111-asd",par2="222-qwe",cpu_no="111as000"})`, pass: true},
		{query: `sum(cpu_total{par1="111-asd", par2="222-qwe",cpu_no="111as000"})`, pass: true},
		{query: `sum({cpu_total="asd"})`, pass: true},
		{query: ` sum(cpu_total)`, pass: true},
		{query: `sum(cpu_total) `, pass: true},
		{query: `sum( cpu_total)`, pass: true},
		{query: `sum(cpu_total )`, pass: true},
		{query: `sum(cpu_total{})+sum(cpu_count{})`, pass: false},
		{query: `sum(max(metric))`, pass: true},
		{query: `sum(max by(id)(metric{label1="value-1",label2="value2"}))`, pass: true},
		{query: `sum(max by(id)(metric))`, pass: true},
		{query: `sum(min by(id)(metric))`, pass: true},
		{query: `sum(avg by(id)(metric))`, pass: true},
		{query: `sum(max by(id)(metric{}))`, pass: true},
		{query: `sum(max by(id, id2)(metric{}))`, pass: false}, // issue #15
		{query: `sum(max by(id) (metric{label1="value-1",label2="value2"}))`, pass: true},
		{query: `sum(max(metric{label1="value-1",label2="value2"}))by(id)`, pass: true},
		{query: `sum(max(metric{label1="value-1",label2="value2"})) by (id)`, pass: true},
		{query: `sum(topk(metric)) by (id)`, pass: false},
		{query: `sort(sum(metric) by (id))`, pass: false},
		{query: `sum(foo{bar=~"$^"})`, pass: false},   // issue #1
		{query: `foo{bar=~”$^”}`, pass: false},        // issue #1
		{query: `foo{bar=~”$^|baz”}`, pass: false},    // issue #1
		{query: `sum`, pass: false},                   // issue #2
		{query: `absent_over_time()`, pass: false},    // issue #6
		{query: `quantile(φ, multi_10)`, pass: false}, // issue #8
	}

	scope := tally.NewTestScope("", nil)
	r := NewRouter(RouterOptions{Scope: scope})

	for _, tt := range tests {
		pass := r.isSupportedByM3Query(tt.query)
		assert.Equal(t, tt.pass, pass, tt.query)
	}
}
