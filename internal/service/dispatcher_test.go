package service

import (
	"testing"

	"github.com/M2rk13/Otus-327619/internal/model/api"
	"github.com/M2rk13/Otus-327619/internal/model/log"
)

func TestNewDispatcherService(t *testing.T) {
	d := NewDispatcherService()

	if d == nil {
		t.Fatal("NewDispatcherService returned nil")
	}
}

func TestDispatcherDispatchExampleDataEvenOdd(t *testing.T) {
	d := NewDispatcherService()
	reqCh := make(chan *api.Request, 2)
	respCh := make(chan *api.Response, 2)
	logCh := make(chan *log.ConversionLog, 2)

	d.DispatchExampleData(2, reqCh, respCh, logCh)
	req1 := <-reqCh
	resp1 := <-respCh
	lg1 := <-logCh

	if req1.From != "USD" || req1.To != "EUR" || req1.Amount <= 0 {
		t.Fatalf("bad request from dispatcher: %+v", req1)
	}

	if !resp1.Success || resp1.Query.Id != req1.Id || resp1.Result <= 0 {
		t.Fatalf("bad response (even) from dispatcher: %+v", resp1)
	}

	if lg1.Request.Id != req1.Id || lg1.Response.Id != resp1.Id {
		t.Fatalf("bad log (even): %+v", lg1)
	}

	d.DispatchExampleData(3, reqCh, respCh, logCh)
	req2 := <-reqCh
	resp2 := <-respCh
	lg2 := <-logCh

	if resp2.Success {
		t.Fatalf("expected Success=false on odd iteration, got true: %+v", resp2)
	}

	if lg2.Request.Id != req2.Id || lg2.Response.Id != resp2.Id {
		t.Fatalf("bad log (odd): %+v", lg2)
	}
}
