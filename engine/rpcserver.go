package engine

import (
	"net"
	"net/http"
	"time"

	log "github.com/funkygao/log4go"
	"github.com/gorilla/mux"
)

func (e *Engine) launchRPCServer() {
	e.rpcRouter = mux.NewRouter()
	e.rpcServer = &http.Server{
		Addr:         e.participantID,
		Handler:      e.rpcRouter,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
	}

	e.setupRPCRoutings()

	var err error
	if e.rpcListener, err = net.Listen("tcp", e.rpcServer.Addr); err != nil {
		panic(err)
	}

	go e.rpcServer.Serve(e.rpcListener)
	log.Info("RPC server ready on http://%s", e.rpcServer.Addr)
}

func (e *Engine) stopRPCServer() {
	if e.rpcListener != nil {
		e.rpcListener.Close()
		log.Info("RPC server stopped")
	}
}

func (e *Engine) setupRPCRoutings() {
	e.rpcRouter.HandleFunc("/v1/rebalance", e.doLocalRebalance).Methods("POST")
}
