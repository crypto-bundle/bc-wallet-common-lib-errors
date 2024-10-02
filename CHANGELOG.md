# Change Log

## [v0.0.6] - 30.09.2024
### Added
* Added value based error format service component, currently support next kind of values:
  * KindDetails
  * KindScope 
  * KindCode 
  * KindPublicCode
* Added unit-tests for value-based formatter
* Added new public receiver-methods for all formatters:
  * ErrorWithCode/ErrWithCode
  * ErrorGetCode/ErrGetCode
### Changed
* Removed usage of getFuncName(), now all errors wrapped without name of caller-function in error text 

## [v0.0.4, v0.0.5] - 27.09.2024
### Added
* Added ErrNoWrap/ErrorNoWrap method to all implementations of error formatters
### Changed
* Fixed linter issues
* Added new settings to wrapcheck linter

## [v0.0.3] - 29.08.2024
### Added
* Improvements for scope-based error wrappers

## [v0.0.2] - 28.08.2024
### Added
* Added new scoped error service-component
* Added scoped exported functions

## [v0.0.1] - 28.08.2024 23:56 MSK
### Added
* Added new error formatter library
* Added LICENSE file
* Added go-linters to project