package loghelper

import "sync"

var logClient = &LogClient{}
var starter = sync.Once{}

func GetNotifyLogClient() *LogClient {
	starter.Do(func() {
		logClient.Init("COMMENT_LOG")
	})

	return logClient
}
