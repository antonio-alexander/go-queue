package internal

import (
	"time"
)

//DefaultInfiniteTimeout provides a default time to send a signal
const DefaultSignalTimeout = time.Duration(0)

//ConfigSignalTimeout is a global variable that can be used to configure
// the signal timeout
var ConfigSignalTimeout = DefaultSignalTimeout
