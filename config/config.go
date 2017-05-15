// Config is put into a different package to prevent cyclic imports in case
// it is needed in several locations

package config

import "time"

type Config struct {
	Period time.Duration `config:"period"`
	Tuxdir string        `config:"tuxdir"`
	Home   string        `config:"home"`
}

var DefaultConfig = Config{
	Period: 1 * time.Second,
	Tuxdir: "c:/psft/pt/bea/",
	Home:   "c:/Users/vagrant/psft/pt/8.55/appserv/APPDOM/PSTUXCFG",
}
