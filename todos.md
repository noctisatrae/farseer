## plugin system
- [X] implement `TestIndividualLoader` to test `LoadHandler`
- [X] implement `TestListCompiledHandlers` to test `ListCompiledHandlers` (should list just the file names in `compiled_handlers`)
- [X] implement  `TestMultipleLoader` to test `LoadHandlersFromConf` (should add some debugging to the function)
- [X] enable plugins from `config.toml`
- [X] develop basic plugins! 

## grafana
- [ ] find a way to get metrics from memory & go
- [ ] how to integrate/launch the server (should just provide the json for the dashboard?)

## plugin ideas
- [X] Find a way to make a `JS`/`TS` sdk! => gRPC
- [X] Sink to RedisDB/Dragonfly/NoSQL DB (sink to DB) => in process of doing it!
- [X] Simple cast filter/cast tracker (for example of use of the handler API)

## libp2p stuff
- [X] Do I need to regossip the messages? => Answer from V: no

## PostgreSQL
- [X] Be **SURE** message type & FID filtering works! 
- [X] Spend some time on data modelling (how to make it simple but efficient) & send a message to Alex for advices!
- [X] Implement LINKS type message (for follows & unfollows)
- [X] Disable saving certain types of messages in the chat from `config.toml`
- [X] Ask around to see what kind of data modeling would be suitable to Hub messages in the DB
- [X] Check out Shuttle/Neynar: one big table with IDs to differenciate the messages from one another
- [X] How to compute cast hashes so you can query them later? 
- [ ] Implement `VerificationAdd` for the message

## Hub stuff
- [X] Absolute path for `config.toml`. If relay is executed in a folder, search the config from the context of execution.
- [X] Figure out which timestamp Farcaster uses? => farcaster time
- [X] Do we receive new messages from the network or sync messages? => find out using timestamps
- [X] gRPC API to write message & act more as real hub.
- [X] save private key to file for persistent identity
- [ ] Create a tool to generate the `hub_identity` file so you can get persistence inside the container

## gRPC
- [ ] `config.toml` to set options of the server
- [ ] implement getCurrentPeers endpoint
- [ ] TLS auth for getSecureSSLClient()

## The project in itself
- [ ] Branding & asserts for README.md
- [X] Choose a license that allows you to make a living out of this
- [X] Make pre-compiled binaries available in the release

## CI/CD & Automatic release
- [ ] Write a `Makefile` to automatically build the binaries, compress them and release to GitHub
- [ ] Create a Github job to build for different platforms (Mac/Windows)
- [ ] Write documentation so people can extend the jobs & `Makefile` to build customs plug-ins