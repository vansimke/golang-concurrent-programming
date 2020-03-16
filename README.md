To use the test runner:

* navigate into the ./cmd directory
* if run with the defaults (`go run .`), the application will run in asynchronous mode and write it's output directly to the shell
    * the `-out` flag allows a destination file to be specified
    * the `-async` flag determines if the logger will asynchronously or not. It is async by default, which will not work properly in the application's initial condition.

To run the tests

* use the `go test` command from the command line
* to focus on a specific test, use the `-run` flag followed by a pattern that matches the test name
    * e.g. `go test -run MessageChannel` will run the first test of the first module
