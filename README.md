# Secondly

Secondly is a configuration management plugin for Go projects. It taks care of
the app's configuration, specifically of updating it in runtime.

## Configuration

First we need to define a struct that will hold app's configuration. Let's make
it simple for demostration purposes.

```go
type Config struct {
    AppName string  `json:"app_name"`
    Version float32 `json:"version"`
}
```

Make sure you've defined `json` tags on each field.

Next, right where you will define your app's flags ask Secondly to add one for
configuration file.

```go
secondly.SetupFlags()
flag.Parse()
```

Now you can pass a configuration file to your program like this:

```
./app -config=config.json
```

Now we need to ask Secondly to take care of your configuration:

```go
var conf Config
secondly.Manage(&conf)

// or asynchronously
go secondly.Manage(&conf)
```

If you prefer to configure the app asynchronously then you'll probably want to
know when configuration is loaded, so there's a handly helper function just for
that:

```go
secondly.OnLoad(func(){
    log.Println("Configuration initialized")
})
```

Congratulations! You've just configured Secondly to read and initialize the
configuration of your app. But this is not what you came for, right?

If you want a configuration GUI, simply start Secondly's web server on a port
you want.

```go
secondly.StartServer("", 5500)
```

Tired of restarting the app every time you modify the config? You're not alone.

```go
secondly.HandleFSEvents()
```

Want some more control over when specifically the config will be reloaded? Ask
Secondly to listen for SIGHUP syscalls.

```go
secondly.HandleSIGHUP()
```

You can also set up callback functions on specific fields and receive a call
when this fields value changes.

```go
secondly.OnChange("NumWorkers", func(oldVal, newVal interface{}) {
    old := oldVal.(int)
    cur := newVal.(int)
    if cur > old {
        pool.AddWorkers(cur - old)
    } else {
        pool.StopWorkers(old - cur)
    }
    log.Println("Number of workers changed from %d to %d", old, cur)
}
})
```

Full example can be found [here](https://github.com/localhots/secondly/blob/master/demo/demo.go).

## Demo Screenshot

<img src="https://raw.githubusercontent.com/localhots/secondly/master/demo/screenshot.png" width="440">

## Building

The only thing to keep in mind when building Secndly is to convert assets into a
binary form so they could be kept in memory of your app and would not require
any additional web server configuration.

```
go get github.com/GeertJohan/go.rice/rice
rice embed-go
go build
```

## Licence

Secondly is distributed under the [MIT Licence](https://github.com/localhots/secondly/blob/master/LICENCE).
