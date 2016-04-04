god = require("god")
log = require("log")
path = require("path")
os = require("os")
go = require("go")
log.info("hello")

god.ignore("**/.git/**", "**/.vscode/**")

god.init(function() { go.reload("./test") })


//function watch(wildcard, unique, eventCallback)
//wildcard, match the path or file
//unique, if unique is true, only this callback will be called and the others will be ignored. 
god.watch("./test/*.go", false,
    function(event) {
        go.reload(path.dir(event.rel), [], [], function(err) { log.error(err) })
    }
)


god.watch(["*_test.go", "**/*_test.go"], true,
    function(event) {
        go.test(path.dir(event.rel), [], function(err) { log.error(err) })
    }
)

//This will not be called if test files are changed. 
god.watch("**/*.go", false,
    function(event) {
        go.install(path.dir(event.rel), [], function(err) { log.error(err) })
    }
)

god.watch(["*.go"], false,
    function(event) {
        log.info("reload", event.abs)
        go.reload(".", [], [], function(err) { if (err) { log.error(err) } })
    }
)


