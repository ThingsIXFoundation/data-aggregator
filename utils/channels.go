package utils

func WaitForChannelsToClose[T any](chans ...chan T) {
	for _, v := range chans {
		<-v
	}
}
