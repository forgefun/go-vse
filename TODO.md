## TODO
* [ ] Refactor Portfolio methods
  * [ ] Use stock ticker to build fuid
  * [ ] Better segregation between Portfolio and Order structs and their methods
* [x] Requests should be made directly in Client object
  * [x] doRequest() should handle query parameters
  * [x] doRequest() should better handle requests other than GET
* [ ] Unit tests
  * [ ] CI with Travis
  * [ ] Use mock server with, potentially with gorilla/mux
  * [ ] Use test-fixtures to mock page
  * [ ] Use golden files to update test-fixtures
* [ ] Document source for readability in godoc
* [ ] Use godeps to cache dependency in source. No leftpad shenanigans
* [ ] Clean up  and improve logging
