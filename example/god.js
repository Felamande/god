//modules located in github.com/Felamande/jsvm/modules 
//and github.com/Felamande/god/modules
//you can write modules yourself if you're familiar with otto.
god = require("god")
log = require("log")
path = require("path")
os = require("os")
go = require("go")
hk  = require("hotkey")


//bind hotkey as you like
hk.bind("ctrl+shift+o",function() {
    log.info("hotkey","ctrl+shift+o")
})
hk.bind("ctrl+shift+k",function() {
    log.info("hotkey","ctrl+shift+k")
})

//changes of ignored files or dirs will not be watched 
god.ignore(".git", ".vscode")

var buildArgs = []
var installArgs = []
var binArgs = []
var testArgs = [] //args for reloaded binaries 

//will be called immediately after god starts.
god.init(function() {console.log("hello")})


// define your subcommand, flags and arguments will be passed to the callback function.
// (god) subcmd "-willnot=-be-parsed" name=what stringvalue -key=value -testarg=-test.v -godebug=gctrace=1 -boolval
// will be parsed as
// nargs = ["-willnot=-be-parsed", "name=what", "stringvalue"], 
// flags = {"key":"value", "testarg":"-test.v", "godebug":"gctrace=1", "boolvar":true}
god.subcmd("print",function(nargs,flags){
   log.info(JSON.stringify(nargs),JSON.stringify(flags)) 
})

god.subcmd("eval",function (nargs,flags) {
    console.log(eval(nargs[0]))
})

god.subcmd("test",function(pkgs,flags){
       log.info("test",pkgs[0])
        go.test(pkgs[0], testArgs, function(err) { log.error(err) })  
})

god.subcmd("install",function(pkgs,flags){
       log.info("test",pkgs[0])
        go.test(pkgs[0], testArgs, function(err) { log.error(err) })  
})

god.subcmd("exec",function(nargs,flags){os.exec(nargs)})

// function watch(name, wildcard, isUnique, callback)
// if isUnique, the event which matches multiple wildcards will only be sent to the unique callback.
//  
// function callback(event) 
// event.rel, relative path of matched file or directory
// event.abs, absolute path of matched file or directory
// event.dir, relative parent directory of matched file or directory
//
// path seperator will be slash on windows.
// watch tasks will not start until you type the subcommand "watch [taskname...]", 
// after that tasks will run in a goroutine.

god.watch("btest","*_test.go", true,
    function(event) {
        log.info("test",event.dir)
        go.test(event.dir, testArgs, function(err) { log.error(err) })
    }
)

// ** will match ONE or more directories
// * will match just ONE directory or as many chars as possible except the slash.
god.watch("ptest", "**/*_test.go", true,
    function(event) {
        log.info("test",event.dir)
        go.test(event.dir, testArgs, function(err) { log.error(err) })
    }
)

god.watch("pinstall","**/*.go", false,
    function(event) {
        log.info("install",event.dir)
        go.install(event.dir, installArgs, function(err) { log.error(err) })
    }
)

god.watch("breload","*.go", false,
    function(event) {
        log.info("reload", event.dir)
        go.reload(".", buildArgs, binArgs, function(err) { if (err) { log.error(err) } })
    }
)

// TODO:
// 1.the way to unwatch tasks.
// 2.separate normal tasks from watch tasks.
