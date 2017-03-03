package input

import (
	"fmt"
	"time"

	"github.com/funkygao/dbus/engine"
	"github.com/funkygao/dbus/pkg/myslave"
	conf "github.com/funkygao/jsconf"
	log "github.com/funkygao/log4go"
)

type MysqlbinlogInput struct {
	stopChan chan struct{}

	slave *myslave.MySlave
}

func (this *MysqlbinlogInput) Init(config *conf.Conf) {
	this.stopChan = make(chan struct{})
	this.slave = myslave.New().LoadConfig(config)

	// so that KafkaOutput can reuse
	key := fmt.Sprintf("myslave.%s", config.String("name", ""))
	engine.Globals().Register(key, this.slave)
}

func (this *MysqlbinlogInput) Stop() {
	log.Trace("stopping...")
	close(this.stopChan)
	this.slave.StopReplication()
}

func (this *MysqlbinlogInput) Run(r engine.InputRunner, h engine.PluginHelper) error {
	backoff := time.Second * 5
	for {
	RESTART_REPLICATION:

		log.Info("starting replication")

		ready := make(chan struct{})
		go this.slave.StartReplication(ready)
		select {
		case <-ready:
		case <-this.stopChan:
			log.Trace("yes sir!")
			return nil
		}

		rows := this.slave.EventStream()
		errors := this.slave.Errors()
		for {
			select {
			case <-this.stopChan:
				log.Trace("yes sir!")
				return nil

			case err := <-errors:
				// e,g.
				// ERROR 1236 (HY000): Could not find first log file name in binary log index file
				// ERROR 1236 (HY000): Could not open log file
				log.Error("backoff %s: %v", backoff, err)
				this.slave.StopReplication()

				select {
				case <-time.After(backoff):
				case <-this.stopChan:
					log.Trace("yes sir!")
					return nil
				}
				goto RESTART_REPLICATION

			case pack, ok := <-r.InChan():
				if !ok {
					log.Trace("yes sir!")
					return nil
				}

				select {
				case err := <-errors:
					// TODO is this neccessary?
					log.Error("backoff %s: %v", backoff, err)
					this.slave.StopReplication()

					select {
					case <-time.After(backoff):
					case <-this.stopChan:
						log.Trace("yes sir!")
						return nil
					}
					goto RESTART_REPLICATION

				case row, ok := <-rows:
					if !ok {
						log.Info("event stream closed")
						return nil
					}

					pack.Payload = row
					r.Inject(pack)

				case <-this.stopChan:
					log.Trace("yes sir!")
					return nil
				}
			}
		}
	}

	return nil
}

func init() {
	engine.RegisterPlugin("MysqlbinlogInput", func() engine.Plugin {
		return new(MysqlbinlogInput)
	})
}
