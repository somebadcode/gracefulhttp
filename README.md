#### Simple HTTP server with signal handling
This is a very simple HTTP server that only uses the standard libraries provided by the Go project.
It's not meant to be a production ready HTTP server but is meant to be a template you get start with
and replace components that could make it more suitable for running in production.

You can freely use this project as a starting point or as a reference on how can you use `context` to gracefully handle
signals and other concurrent code in general.
 
#### Things you *should* do
* Write unit tests for your code.
* Replace the logger with a logger that can provide more context for the errors such as [Zap](https://github.com/uber-go/zap) or any of the other good logging modules.
* Replace the HTTP router and handling with something much better such as [Gorilla Mux](https://github.com/gorilla/mux) or [Gin](https://github.com/gin-gonic/gin). There are many modules to choose from, so it's up to you to pick one for your project.
* Split the project into submodules before it becomes too complex.

### Contribution
Anyone can contribute as long as it doesn't add new functionality or anything that requires a change of license. I only accept fixes since this is only meant as a starting point or as a reference for other projects.