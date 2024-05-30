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
- [X] Sink to RedisDB/Dragonfly/NoSQL DB (sink to DB) => in process of doing it!
- [ ] Simple cast filter/cast tracker (for example of use of the handler API)

## libp2p stuff
- [ ] Do I need to regossip the messages? 

## PostgreSQL
- [ ] Disable saving certain types of messages in the chat from `config.toml`
- [ ] Ask around to see what kind of data modeling would be suitable to Hub messages in the DB
- [ ] Check out Shuttle/Neynar: one big table with IDs to differenciate the messages from one another
