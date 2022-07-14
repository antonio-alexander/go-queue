package infinite

import "time"

//DefaultSignalTimeout provides a default time to send a signal
const DefaultSignalTimeout = time.Duration(time.Millisecond)

//ConfigSignalTimeout is a global variable that can be used to configure
// the signal timeout
var ConfigSignalTimeout = DefaultSignalTimeout
