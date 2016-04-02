god = require("god")
log = require("log")
path = require("path")
os = require("os")
go = require("go")
log.info("hello")

god.ignore(".git", ".vscode")

god.init(function(){go.reload("./test")})

god.watch("./test/*.go",function(abs,rel) {
    go.reload(path.dir(rel),[],[],function(err) {
        if(err){log.error(err)}
    })
})


// function rebuild(name) {
//     os.system([tool, "build", "-o", name + "_tmp.exe"])
//     os.system(["taskkill", "/im", name + ".exe"], "/f")
//     os.rename(name + "_tmp.exe", name + ".exe")
//     os.exec(name + ".exe")
// }

// god.watch(["*_test.go", "**/*_test.go"],
//     function(abs, rel) {
//         os.system([tool, "test", path.dir(rel)])
//     }
// )

// god.watch("**/*.go",
//     function(abs, rel) {
//         if (rel.indexOf("_test.go") >= 0) return;
//         os.system([tool, "install", path.dir(rel)])
//     }
// )

// god.watch(["*.go"],
//     function(abs, rel) {
//         if (rel.indexOf("_test.go") >= 0) return;
//         rebuild(os.wdName())

//     }
// )


