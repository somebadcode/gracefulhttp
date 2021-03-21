#### Simple HTTP server with graceful shutdown
This is a simple HTTP server that allows for graceful shutdowns by using context. I've only included  one third party
library which I use for logging, so it has nothing to do with the solution in and of itself. I have also included a
small example of how to embed files into the final binary and in this case it's part of a small demo web application.

You can freely use this project as a starting point or as a reference on how can you use `context` to gracefully handle
signals and other concurrent code in general.

### Contribution
Anyone can contribute as long as it doesn't add new functionality or anything that requires a change of license.
I only accept fixes and improved solutions since this is only meant as a starting point or as a reference for other
projects.