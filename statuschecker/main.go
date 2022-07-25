package main

import (
	"fmt"
	"net/http"
	"time"
)

func main() {

	links := []string{
		"http://google.com",
		"http://facebook.com",
		"http://stackoverflow.com",
		"http://golang.org",
		"http://amazon.com",
	}

	// Prepare a channel for communication of strings
	c := make(chan string)

	// First for loop: making sure we spin up a go routine for each link to achieve concurrency
	for _, link := range links {
		// Spin up a child go routine and run this function within it
		go checkLink(link, c)
	}

	// Second for loop: make sure each go routine we spun up in the first for loop is going to be repeated forever
	// We will use this for loop to deal with the idea of having many go routines that communicate
	// with the main go routine on a common channel
	// This to fix the problem of the main go routine terminatiing upon waking up to receive something from
	// the channel.
	for l := range c {
		// Waiting on data coming from a channel will block the main go routine and put it to sleep
		// This will block the for loop until the next wakeup and that's what we need exactly
		// We have to register channel waits greater than or equal the number of go routines we spin up

		// go checkLink(l, c)
		// time.Sleep(2 * time.Second)

		// Little gotcha: variable l in the function literal is referencing the value in the main routine scope
		// we have to pass down the variable l to the function literal so that the checkLink would work correctly
		// with the correct version of l
		go func(link string) {
			time.Sleep(5 * time.Second)
			go checkLink(link, c)
		}(l)
	}
}

func checkLink(link string, c chan string) {
	_, err := http.Get(link)
	if err != nil {
		fmt.Println(link, "might be down!")

		// Pushing the link to the channel in order to deal with the trick of repeating routines
		c <- link
		return
	}

	fmt.Println(link, "is up")
	c <- link
}
