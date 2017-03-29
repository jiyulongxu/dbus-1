package engine

import (
	"io"
	"net/http"

	"github.com/funkygao/dbus/pkg/cluster"
	log "github.com/funkygao/log4go"
)

const (
	maxRPCBodyLen = 128 << 10
)

func (e *Engine) doLocalRebalance(w http.ResponseWriter, r *http.Request) {
	if r.ContentLength > maxRPCBodyLen {
		log.Warn("too large RPC request body: %d", r.ContentLength)
		w.WriteHeader(http.StatusRequestEntityTooLarge)
		return
	}

	buf := make([]byte, r.ContentLength)
	_, err := io.ReadFull(r.Body, buf)
	if err != nil {
		log.Error("%v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	resources := cluster.RPCResources(buf)
	log.Trace("got %d resources: %v", len(resources), resources)

	// merge resources by input plugin name
	resourceMap := make(map[string][]string) // inputName:resources
	for _, res := range resources {
		if _, present := resourceMap[res.InputPlugin]; !present {
			resourceMap[res.InputPlugin] = []string{res.Name}
		} else {
			resourceMap[res.InputPlugin] = append(resourceMap[res.InputPlugin], res.Name)
		}
	}

	// dispatch decision to input plugins
	for inputName, resources := range resourceMap {
		ir, ok := e.InputRunners[inputName]
		if !ok {
			// should never happen
			panic(inputName + " not found")
		}

		ir.feedResources(resources)
		ir.rebalance()
	}

}
