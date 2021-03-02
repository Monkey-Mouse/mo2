package importService

import (
	"fmt"
	"io/ioutil"
)

func ExampleTransform() {
	//t1:=time.Now()
	source, err := ioutil.ReadFile("2020-04-16-《生活十讲》-爱与情.md")
	if err != nil {
		panic(err)
	}
	blog := Transform([]byte(source))
	//fmt.Println(time.Since(t1))
	fmt.Println(blog.Title)
	fmt.Println(blog.EntityInfo.CreateTime)

	// Output:
	// 《生活十讲》——爱与情
	// 2020-04-16 00:00:00 +0000 UTC
}
