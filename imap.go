package main

import "github.com/BrianLeishman/go-imap"

func emailClient() *imap.Dialer {
	// Defaults to false. This package level option turns on or off debugging output, essentially.
	// If verbose is set to true, then every command, and every response, is printed,
	// along with other things like error messages (before the retry limit is reached)
	imap.Verbose = IMAP_DEBUG

	// Defaults to 10. Certain functions retry; like the login function, and the new connection function.
	// If a retried function fails, the connection will be closed, then the program sleeps for an increasing amount of time,
	// creates a new connection instance internally, selects the same folder, and retries the failed command(s).
	// You can check out github.com/StirlingMarketingGroup/go-retry for the retry implementation being used
	imap.RetryCount = 3

	// Create a new instance of the IMAP connection you want to use
	im, err := imap.New(IMAP_USER, IMAP_PASS, IMAP_SERVER, IMAP_PORT)
	check(err)
	return im
}
