## Changelog

All notable changes to this project will be documented in this section.

### [v1.4.1] - 2025-01-13
#### Fixed
- project module definition

### [v1.4.0] - 2025-01-13
#### Added
- Implemented support for multi-threaded cache access. Now it's thread safe!
- Added README, LICENSE and CHANGELOG files

### [v1.3.0] - 2025-01-13
#### Added
- Implemented Delete service.
- Implemented ListAll service.

### [v1.2.1] - 2025-01-13
#### Changed
- Small performance improvement in GET method.
- Small code quality improvement in `api_test.go` file

### [v1.2.0] - 2025-01-13
#### Added
- Support Get method.

### [v1.0.0] - 2025-02-12
#### Initial Release
- Implemented a basic n-way set associative cache.
- Supported LRU and MRU eviction policy by default.
- Implemented Put service.