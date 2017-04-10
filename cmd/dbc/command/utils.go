package command

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/funkygao/dbus"
	"github.com/funkygao/dbus/engine"
	"github.com/funkygao/dbus/pkg/cluster"
	czk "github.com/funkygao/dbus/pkg/cluster/zk"
	"github.com/funkygao/gafka/ctx"
	"github.com/funkygao/gorequest"
)

func openClusterManager(zone string) cluster.Manager {
	mgr := czk.NewManager(ctx.ZoneZkAddrs(zone), engine.Globals().ZrootCluster)
	swallow(mgr.Open())

	return mgr
}

func swallow(err error) {
	if err != nil {
		panic(err)
	}
}

func callAPI(p cluster.Participant, api string, method string, body string) (string, []error) {
	r := gorequest.New()
	uri := fmt.Sprintf("%s/api/v1/%s", p.APIEndpoint(), api)
	switch strings.ToUpper(method) {
	case "PUT":
		r = r.Put(uri)
	case "POST":
		r = r.Post(uri)
	case "GET":
		r = r.Get(uri)
	}

	reply, replyBody, errs := r.
		Set("User-Agent", fmt.Sprintf("dbus-%s", dbus.Revision)).
		SendString(body).
		End()
	if reply.StatusCode != http.StatusOK {
		return "", []error{fmt.Errorf("status %d", reply.StatusCode)}
	}
	return replyBody, errs
}
