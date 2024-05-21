package reload

import (
	"fmt"
	"io"
	"log"
	"net/http"
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
	go reloadServer(p)
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

func reloadServer(bin string) {

	reloadPort := os.Getenv("HOT_RELOAD_PORT")
	if reloadPort == "" {
		reloadPort = "8087"
	}
	http.HandleFunc("/upload", uploadFile(bin))

	log.Println("Starting hot reload server at :" + reloadPort)
	if err := http.ListenAndServe(":"+reloadPort, nil); err != nil {
		fmt.Println("Error starting server:", err)
	}
}

func uploadFile(f string) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		file, handler, err := r.FormFile("file")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer file.Close()

		// 创建临时文件
		tempFile, err := os.CreateTemp("", "upload-*.tmp")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// 将上传的文件内容写入临时文件
		_, err = io.Copy(tempFile, file)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		_ = tempFile.Close()

		fmt.Fprintf(w, "File uploaded successfully: %s\n", handler.Filename)
		fmt.Fprintf(w, "Saved to: %s\n", tempFile.Name())
		stat, err := os.Stat(tempFile.Name())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// 设置新的权限，添加执行权限
		newMode := stat.Mode() | 0111 // 添加用户、组和其他人的执行权限
		err = os.Chmod(tempFile.Name(), newMode)
		if err != nil {
			fmt.Fprintf(w, "chmod +x get error: %s\n", err.Error())
			return
		}
		fmt.Fprintf(w, "chmod +x %s\n", tempFile.Name())
		os.Rename(tempFile.Name(), f)
	}

}
