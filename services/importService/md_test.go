package importService

import (
	"io/ioutil"
	"strings"
	"time"

	"github.com/Monkey-Mouse/mo2/mo2utils"
)

func ExampleTransform() {
	c := make(chan []byte)
	done := make(chan bool)
	count := 0
	t1 := time.Now()

	go func() {
		for {
			res, ok := <-c
			if ok {
				blog := Transform(res)
				println(blog.Title)
				println(count, "write", time.Since(t1))
				count--
			} else {
				println("finished!")
				done <- true
				return
			}
		}
	}()

	mo2utils.ProcessAllFiles("./", "", func(parameter ...string) {
		if strings.HasSuffix(parameter[0], ".md") {
			source, err := ioutil.ReadFile(parameter[0])
			if err != nil {
				panic(err)
			}
			c <- source
			count++
			println(count, " read", time.Since(t1))
		}
	})
	close(c)
	<-done
	// Output:
}
