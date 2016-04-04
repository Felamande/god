//modules located in github.com/Felamande/jsvm/modules 
//and github.com/Felamande/god/modules
//you can write modules yourself if you're familiar with otto.
god = require("god")
log = require("log")
path = require("path")
os = require("os")
go = require("go")

//changes of ignored files or dirs will not be watched 
god.ignore(".git", ".vscode")

var buildArgs = []
var installArgs = []
var binArgs = []
var testArgs = [] //args for reloaded binaries 

//will be call after god starts and before god watches changes.
god.init(function() {console.log("hello");go.reload(".", buildArgs, binArgs)})

// event 
// event.rel, relative path of matched file or directory
// event.abs, absolute path of matched file or directory
// event.dir, relative parent directory of matched file or directory
// path seperator will be slash on windows.
god.watch(["*_test.go", "**/*_test.go"], true,
    function(event) {
        log.info("test",event.dir)
        go.test(event.dir, testArgs, function(err) { log.error(err) })
    }
)

// ** will match ONE or more directories
// * will match just ONE directory or as many chars as possible except slash .
god.watch("**/*.go", false,
    function(event) {
        log.info("install",event.dir)
        go.install(event.dir, installArgs, function(err) { log.error(err) })
    }
)

god.watch(["*.go","**/*.go"], false,
    function(event) {
        log.info("reload", event.dir)
        go.reload(".", buildArgs, binArgs, function(err) { if (err) { log.error(err) } })
    }
)

