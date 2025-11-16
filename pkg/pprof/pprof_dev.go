//go:build !prd
// +build !prd

package pprof

import (
	"net/http"
	_ "net/http/pprof"
)

func InitPprof() {
	// Dev 환경에서만 pprof 활성화
	go func() {
		http.ListenAndServe("0.0.0.0:6060", nil)
	}()
}
