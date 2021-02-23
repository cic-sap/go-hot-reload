# go-hot-reload

Keep the pid of the process, replace the executable file

example code

```go
package main

import (
	"fmt"
	reload "github.com/cic-sap/go-hot-reload/v1"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"time"
)
/**

go build -o main main.go
./main &
# change your code
go build -o main main.go

 */
func main() {

	log.SetFlags(log.Llongfile | log.LstdFlags)
	ver := "9"
	go reload.AutoReload()

	// your app code ...
	go func() {
		for {
			time.Sleep(time.Second)
			log.Println("ver", ver, "pid", os.Getpid())
		}
	}()

	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, fmt.Sprintf("hello ver:%s", ver))
	})

	r.Run("0.0.0.0:2345")

}

```

How to replace k8s pod process online


example
```shell
set -ex
export KUBECONFIG=~/.kube/kubeconfig--gcp01.yaml

(GOOS=linux go build -o myapp main.go ;mv myapp myapp.tmp;gzip -fv myapp.tmp) &

name=`kubectl get pods  -n ns1 -l app=myapp -o=jsonpath='{.items[0].metadata.name}'`
wait
kubectl get pods  -n ns1 -l app=myapp &

kubectl -n ns1 cp  myapp.tmp.gz  $name:/usr/local/bin/myapp.tmp.gz
kubectl -n ns1 exec  $name --  sh -x -c 'cd /usr/local/bin/;pwd;ls -alh myapp*;gzip -f -d -v myapp.tmp.gz;mv -f -v myapp.tmp myapp;ls -alh myapp*;'
```