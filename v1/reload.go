package reload

import (
	"log"
	"os"
	"syscall"
	"time"
)

func AutoReload() {

	now := time.Now()
	log.SetFlags(log.LstdFlags | log.Llongfile)
	p, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("go-hot-reload: get process path:", p, "boot time", now)
	fi, err := os.Stat(p)
	if err != nil {
		log.Fatal(err)
	}
	t1 := fi.ModTime()
	for {
		time.Sleep(time.Second * 3)
		fi, err = os.Stat(p)

		if err != nil {
			log.Fatal(err)
		}
		t2 := fi.ModTime()
		if t2 != t1 {
			log.Println("go-hot-reload: try reload:", p, "boot time", now, "uptime", time.Now().Sub(now).Seconds())
			err = syscall.Exec(p, os.Args, os.Environ())
			if err != nil {
				log.Fatal("err", err)
			}
		}

	}
}
