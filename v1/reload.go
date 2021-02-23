package reload

import (
	"log"
	"os"
	"syscall"
	"time"
)

func AutoReload() {

	p, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("get process path:", p)
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
			log.Println("try reload:", p)
			err = syscall.Exec(p, os.Args, os.Environ())
			if err != nil {
				log.Fatal("err", err)
			}
		}

	}
}
