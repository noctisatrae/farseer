## plugin system
- [X] implement `TestIndividualLoader` to test `LoadHandler`
- [X] implement `TestListCompiledHandlers` to test `ListCompiledHandlers` (should list just the file names in `compiled_handlers`)
- [X] implement  `TestMultipleLoader` to test `LoadHandlersFromConf` (should add some debugging to the function)
- [X] enable plugins from `config.toml`
- [ ] develop basic plugins! 

## grafana
- [ ] find a way to get metrics from memory & go
- [ ] how to integrate/launch the server (should just provide the json for the dashboard?)

## plugin ideas
- [ ] Find a way to make a `JS`/`TS` sdk!
- [ ] Sink to RedisDB/Dragonfly/NoSQL DB
- [ ] Simple cast filter/cast tracker (for example of use of the handler API)