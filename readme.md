# farseer - another kind of Farcaster hub
farseer is a lightweight re-implementation of a Farcaster hub that **does not require syncing to operate** & is extendable through plugins. In short, you can bring your own DB, logic & infrastructure and harvest data from the protocol while fostering decentralization! If you like this project, **consider giving a star to the repository**: it's a real help for motivation! Looking for PRs, issues & feedback, so feel free to write something and send it to me! 

See the [todos](./todos.md) to see what's left to do :)

## How to get started?
### General architecture
There are four components to farseer: a config file (`config.toml`), plugins, the identity file & the hub itself. Here's a quick rundown of how it works:
```
.
├── compiled_handlers <== where you'll put your compiled plugins
│   └── postgresql.so 
├── config.toml <== configure the behaviour of the hubs & the plugin
├── docker-compose.yml <== infrastructure example
├── Dockerfile <== automatization of the process
├── hub_identity <== SECRET private key of the hub (needs to be generated for docker-compose)
├── relay <== what you'll run (binary)
```
### Easy mode (Docker)
[Compiling the plugins for Docker](#compiling-plugins-for-docker)
1. Generate a `hub_identity` using the latest utility found in the [release section](https://github.com/noctisatrae/farseer/releases) or run the code in the `identity` folder. **Don't forget to put in the root of the repository!**
2. Run this command to start the containers!
```sh
docker-compose up -d
```
This will start the hub with the default behavior & plug-ins.

### The DIY way
1. Get the source from somewhere:
```sh
git clone https://github.com/noctisatrae/farseer.git
```

2. Compile with Go 1.22+ and produce a binary **in the same directory** where `config.toml` is:
```sh
go build -v -o app ./relay
```

3. Compile your plugins/custom handlers using the *plugin mode* of `go build` (here we'll compile the example `postgresql` plugin):
```sh
go build -buildmode=plugin -o ./compiled_handlers/postgresql.so postgresql/postgresql.go
```
There's a lot going here but essentially, we tell Go to build a plugin from the postgresql.go file output it in the `compiled_handlers` folder that will be read by the hub to exectute the custom logic.

4. Now, you'll start the hub by running: 
```sh
./app
```
5. For fine-tuning the behavior of the hub, see the [configuration section](#configuration)

## Configuration
```toml
[hub]
# How can other peers reach your hub!
PublicHubIp = "92.158.95.48"
GossipPort = 2282
# Who will be your first contacts?
# Quick rundown of the libp2p multiaddr format: 
# /typeOfAddr/addr/protocol/port/p2p/publicIdentity
BootstrapPeers = [
  # Those are the peers used by the farcaster dev team, nemes.farcaster.xyz is the public one & the others will certainly not make the connection with you!
  # As of the 25th of July 2024, nemes.farcaster.xyz is not working
  # "/dns/lamia.farcaster.xyz/tcp/2283/p2p/12D3KooWJECuSHn5edaorpufE9ceAoqR5zcAuD4ThoyDzVaz77GV",
  # "/dns/nemes.farcaster.xyz/tcp/2283/p2p/12D3KooWMQrf6unpGJfLBmTGy3eKTo4cGcXktWRbgMnfbZLXqBbn",
  # "/dns/hoyt.farcaster.xyz/tcp/2283/p2p/12D3KooWRnSZUxjVJjbSHhVKpXtvibMarSfLSKDBeMpfVaNm1Joo",
]
# Super handy when things go wrong!
Debug = false
# Not sure of the usefulness of this, it's something I have yet to experiment with
BufferSize = 128
ContactInterval = 30

# The interesting part!
# To define the behavior of a plugin in `compiled_handlers`, you write:
# [handlers.(pluginName)]
[handlers.postgresql]
# This is common to all plugins: do you want to enable it?
Enabled = true
# Below, the options are specific:
# The options below are determined to by the developer of the plugin. They manage how the arguments are parsed and used!
DbAddress = "postgres://postgres:example@db:5432/postgres"
# refer to the enum l.60 in message.proto for the integer of msg types | here we only want to save the casts & deletions
# delete a filter to not use it!
MessageTypesAllowed = [1, 2]
# who are you tracking?
FidsAllowed = [10626]
```
## Plugins
## Handler API
At some point, you'll want to make your own plug-ins. To get started, you should look at `handlers/handlers.go`! A plugin exports a Handler `struct` defining its own function to handle the message; here's an excerpt from the `struct`:
```go
type Handler struct {
	Name string
  // Used to make a connection to the DB. Go to handlers/handlers.go to see a method to pass down variables to the functions.
	InitHandler               InitBehaviour
  // Those functions will handle incoming messages! It's up to you to define those you need.
	CastAddHandler            HandlerBehaviour
	CastRemoveHandler         HandlerBehaviour
	FrameActionHandler        HandlerBehaviour
	ReactionAddHandler        HandlerBehaviour
	ReactionRemoveHandler     HandlerBehaviour
	LinkAddHandler            HandlerBehaviour
	LinkRemoveHandler         HandlerBehaviour
	VerificationAddHandler    HandlerBehaviour
	VerificationRemoveHandler HandlerBehaviour
}

var PluginHandler = handler.Handler{
  // .... amazing stuff here
}

// Then you compile & put it in compiled_handlers!
```
It's up to you to define & verify the paramaters that will be used in `config.toml`.
### Compiling plugins for Docker
You can edit the project's Dockerfile to add your plugin build command! 
```diff
FROM golang:1.22

WORKDIR /usr/src/app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# 1. Make sure your plugin is included into the image
COPY . .
+ RUN go build -buildmode=plugin -o ./compiled_handlers/[plugin name].so [your source code for the plugin] 
# 2. Example for the postgresql plugin
RUN go build -buildmode=plugin -o ./compiled_handlers/postgresql.so postgresql/postgresql.go
# Then, build the hub itself
RUN go build -v -o /usr/local/bin/app ./relay

CMD ["app"]
```