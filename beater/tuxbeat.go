package beater

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"time"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
	"github.com/elastic/beats/libbeat/publisher"

	"github.com/psadmin-io/ps-tuxbeat/config"
)

type Tuxbeat struct {
	done   chan struct{}
	config config.Config
	client publisher.Client
	tuxdir string
	home   string

	lastIndexTime time.Time
}

// New beater
func New(b *beat.Beat, cfg *common.Config) (beat.Beater, error) {
	config := config.DefaultConfig
	if err := cfg.Unpack(&config); err != nil {
		return nil, fmt.Errorf("Error reading config file: %v", err)
	}

	bt := &Tuxbeat{
		done:   make(chan struct{}),
		config: config,
	}
	return bt, nil
}

/// *** Beater interface methods ***///

func (bt *Tuxbeat) Run(b *beat.Beat) error {
	logp.Info("lsbeat is running! Hit CTRL-C to stop it.")

	bt.client = b.Publisher.Connect()
	ticker := time.NewTicker(bt.config.Period)

	for {
		now := time.Now()
		// bt.listDir(bt.config.Path, b.Name) // call lsDir
		bt.captureDomainStatus(bt.config.Tuxdir, bt.config.PsCfgHome, bt.config.Domain)
		bt.returnTuxDir(bt.config.Tuxdir, bt.config.PsCfgHome, bt.config.Domain, b.Name)
		bt.lastIndexTime = now // mark Timestamp

		logp.Info("Event sent")

		select {
		case <-bt.done:
			return nil
		case <-ticker.C:
		}
	}
}

func (bt *Tuxbeat) Stop() {
	bt.client.Close()
	close(bt.done)
}

func (bt *Tuxbeat) captureDomainStatus(tuxdir string, cfgHome string, domain string) {
	PSTUXCFG := cfgHome + "/appsrv/" + domain + "/PSTUXCFG"
	cmd := exec.Command("$env:TUXCONFIG=\""+PSTUXCFG+"\"", "")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Setting environment variable: %q\n", out.String())
	fmt.Println("")
	// out, err := exec.Command(tuxdir+"/tmadmin", "-r").Output()
}

func (bt *Tuxbeat) returnTuxDir(tuxdir string, cfgHome string, domain string, beatname string) {

	PSTUXCFG := cfgHome + "/appsrv/" + domain + "/PSTUXCFG"
	event := common.MapStr{
		"@timestamp":  common.Time(time.Now()),
		"type":        beatname,
		"tuxdir":      tuxdir,
		"ps_cfg_home": cfgHome,
		"domain":      domain,
		"PSTUXCFG":    PSTUXCFG,
	}
	bt.client.PublishEvent(event)
}

// func (bt *Tuxbeat) listDir(dirFile string, beatname string) {
// 	files, _ := ioutil.ReadDir(dirFile)
// 	for _, f := range files {
// 		t := f.ModTime()
// 		path := filepath.Join(dirFile, f.Name())

// 		if t.After(bt.lastIndexTime) {
// 			event := common.MapStr{
// 				"@timestamp": common.Time(time.Now()),
// 				"type":       beatname,
// 				"modtime":    common.Time(t),
// 				"filename":   f.Name(),
// 				"path":       path,
// 				"directory":  f.IsDir(),
// 				"filesize":   f.Size(),
// 			}

// 			bt.client.PublishEvent(event)
// 		}

// 		if f.IsDir() {
// 			bt.listDir(path, beatname)
// 		}
// 	}
// }
