// Config is put into a different package to prevent cyclic imports in case
// it is needed in several locations

package config

import "time"

type Config struct {
	Period    time.Duration `config:"period"`
	Tuxdir    string        `config:"tuxdir"`
	PsCfgHome string        `config:"ps_cfg_home"`
	Domain    string        `config:"domain"`
}

var DefaultConfig = Config{
	Period:    1 * time.Second,
	Tuxdir:    "c:/psft/pt/bea/",
	PsCfgHome: "c:/Users/vagrant/psft/pt/8.55/",
	Domain:    "APPDOM",
}
