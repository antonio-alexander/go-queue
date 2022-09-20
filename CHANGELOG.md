# Change Log

All notable changes to this service will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/)
and this project adheres to [Semantic Versioning](http://semver.org/).

## [1.2.2] - 09/10/22

- Added "Must" functions that will block until success
- Refactored the inifinite tests to reflect the expected behavior and to be less specific to MY implementation of the infinite queue
- Refactored the goqueue/finite tests to use the MUST functions, still need to do the Peek functions, but otherwise complete
- Updated close function(s) to handle nil/closed channels

## [1.2.1] - 05/01/22

- Updated remaining tests where the Example data wasn't used

## [1.2.0] - 05/01/22

- Somewhat broke the API, removed the "Info" interface, separating the Length and Capacity functions, goqueue owns Length while Capacity is owned by finite
- Updated the tests to use an example data type that also implements BinaryMarshal/BinaryUnmarshal
- Added example package that implements wrappers for common/known queue functions
- Removed the go-modules for infinite/finite, determined that the release cadence was generally coupled to goqueue, so no real reason to separate
- Refactored tests to be easier to integrate into existing test fixtures

## [1.1.2] - 12/03/21

- Updated implementation for SendSignal to NOT provide a timeout (since it's an unbuffered channel)
- Updated default timeout for signal/channel to be non-zero (one millisecond)

## [1.1.1] - 12/02/21

- Updated internal code for sending signal to make the timeout optional and removed the configuration for signal timeout to the package that implements it
- Fixed bug where the signal timeout wasn't configurable

## [1.1.0] - 11/16/21

- Added tests

## [1.0.0] - 09/27/21

- Initial version
