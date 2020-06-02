package httpd

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/m3db/m3/src/query/api/v1/options"
	"github.com/uber-go/tally"
)

var (
	blackListedPatterns = []string{
		`=~\s*"\$`,
		`[^\w]+by\s*\([^),]+,[^)]+\)`, // issue #15
	}

	whiteListedPatterns = []string{
		`^sum\(\s*([a-z0-9_]+)?\s*(\{[^}]*\})?\s*\)$`,
		`^sum\(\s*(max|min|avg)\s*(by\s*\([^)]+\))?\s*\(([a-z0-9_]+)?\s*(\{[^}]*\})?\s*\)\)\s*(by\s*\([^)]+\))?$`,
	}
)

type RouterOptions struct {
	Enabled            bool
	DefaultQueryEngine options.QueryEngine
	PromqlHandler      func(http.ResponseWriter, *http.Request)
	M3QueryHandler     func(http.ResponseWriter, *http.Request)
	LogStats           bool
	Scope              tally.Scope
}

type Router struct {
	enabled               bool
	whiteList             []*regexp.Regexp
	blackList             []*regexp.Regexp
	promqlHandler         func(http.ResponseWriter, *http.Request)
	m3QueryHandler        func(http.ResponseWriter, *http.Request)
	defaultQueryEngine    options.QueryEngine
	logStats              bool
	promqlProcessingTime  tally.Timer
	m3QueryProcessingTime tally.Timer
}

func NewRouter(opts RouterOptions) *Router {
	whiteList := make([]*regexp.Regexp, 0, len(whiteListedPatterns))
	for _, p := range whiteListedPatterns {
		f := regexp.MustCompile(p)
		whiteList = append(whiteList, f)
	}

	blackList := make([]*regexp.Regexp, 0, len(blackListedPatterns))
	for _, p := range blackListedPatterns {
		f := regexp.MustCompile(p)
		blackList = append(blackList, f)
	}

	defaultEngine := opts.DefaultQueryEngine
	if defaultEngine != options.PrometheusEngine && defaultEngine != options.M3QueryEngine {
		defaultEngine = options.PrometheusEngine
	}

	scope := opts.Scope.SubScope("router")
	promql := scope.Tagged(map[string]string{"engine": "prometheus"})
	m3query := scope.Tagged(map[string]string{"engine": "m3query"})

	return &Router{
		enabled:               opts.Enabled,
		whiteList:             whiteList,
		blackList:             blackList,
		defaultQueryEngine:    defaultEngine,
		promqlHandler:         opts.PromqlHandler,
		m3QueryHandler:        opts.M3QueryHandler,
		logStats:              opts.LogStats,
		promqlProcessingTime:  promql.Timer("query-processing-time"),
		m3QueryProcessingTime: m3query.Timer("query-processing-time"),
	}
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	engine := strings.ToLower(req.Header.Get(engineHeaderName))
	urlParam := strings.ToLower(req.URL.Query().Get(engineURLParam))

	if len(urlParam) > 0 {
		engine = urlParam
	}

	autorouted := false
	// if engine is not forced then try to auto-route it
	if !options.IsQueryEngineSet(engine) {
		engine = string(r.defaultQueryEngine)

		// do auto-routing only if promql is the default engine
		if r.enabled && engine == string(options.PrometheusEngine) {
			err := req.ParseForm()
			if err == nil {
				autorouted = true
				query := req.Form.Get("query")
				if len(query) > 0 && r.isSupportedByM3Query(query) {
					engine = string(options.M3QueryEngine)
				}
			}
		}
	}

	w.Header().Add(engineHeaderName, engine)

	if r.logStats {
		start := time.Now()
		defer func() {
			elapsed := time.Since(start)
			fmt.Printf("QSTATS:\t%c\t%c\t%s?%s\t%s\n", engine[0], strconv.FormatBool(autorouted)[0], req.URL.Path, req.URL.RawQuery, elapsed)
		}()
	}

	if engine == string(options.M3QueryEngine) {
		stopwatch := r.m3QueryProcessingTime.Start()
		r.m3QueryHandler(w, req)
		stopwatch.Stop()
		return
	}

	stopwatch := r.promqlProcessingTime.Start()
	r.promqlHandler(w, req)
	stopwatch.Stop()
}

func (r *Router) isSupportedByM3Query(query string) bool {
	query = strings.Trim(strings.ToLower(query), " ")

	for _, f := range r.blackList {
		if f.MatchString(query) {
			return false
		}
	}

	for _, f := range r.whiteList {
		if f.MatchString(query) {
			return true
		}
	}
	return false
}
