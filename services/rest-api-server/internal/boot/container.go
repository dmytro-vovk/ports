package boot

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type Container struct {
	config     *Config
	container  sync.Map
	shutdownFn []serviceShutdownFn
	m          sync.Mutex
}

type serviceShutdownFn struct {
	name string
	fn   func()
}

func New(configFile string) (*Container, func(), error) {
	config, err := readConfig(configFile)
	if err != nil {
		return nil, nil, err
	}

	log.Printf("Config read from %s", configFile)

	c := Container{config: config}

	c.arm(syscall.SIGINT, syscall.SIGTERM)

	return &c, c.shutdown, nil
}

func (c *Container) arm(signals ...os.Signal) {
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, signals...)

	go func() {
		s := <-sc

		log.Printf("Got %v, shutting down...", s)

		c.shutdown()

		log.Printf("Shutdown complete")

		os.Exit(0)
	}()
}

func (c *Container) shutdown() {
	c.m.Lock()

	for i := len(c.shutdownFn) - 1; i >= 0; i-- {
		log.Printf("Shutting down %s", c.shutdownFn[i].name)

		c.shutdownFn[i].fn()
	}

	c.shutdownFn = c.shutdownFn[0:0]

	c.m.Unlock()
}

func (c *Container) get(key string) interface{} {
	if value, ok := c.container.Load(key); ok {
		return value
	}

	return nil
}

func (c *Container) set(key string, value interface{}, shutdown func()) *Container {
	c.container.Store(key, value)

	if shutdown != nil {
		c.shutdownFn = append(c.shutdownFn, serviceShutdownFn{
			name: key,
			fn:   shutdown,
		})
	}

	return c
}

func (c *Container) Config() *Config { return c.config }
