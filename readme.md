# farseer - another kind of Farcaster hub
farseer is a lightweight re-implementation of a Farcaster hub that's extensible through plugins (custom message handlers). In short, you can bring your own DB, logic & infrastructure and harvest data from the protocol.

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
├── hub_identity <== SECRET private key of the hub
├── relay.binary <== what you'll run
```
### Easy mode (Docker)
[Compiling the plugins for Docker](#compiling-plugins-for-docker)
1. Go into the root of the directory & start the container stack:
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
go build -buildmode=plugin -o ./compiled_handlers postgresql/postgresql.go
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
MessageTypesAllowed = [1, 2]
# who are you tracking?
FidsAllowed = [10626]
```

## Compiling plugins for Docker
```dockerfile
FROM golang:1.22

WORKDIR /usr/src/app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# 1. Make sure your plugin is included into the image
COPY . .
# RUN go build -buildmode=plugin -o ./compiled_handlers [your source code for the plugin] 
# 2. Example for the postgresql plugin
RUN go build -buildmode=plugin -o ./compiled_handlers postgresql/postgresql.go
# Then, build the hub itself
RUN go build -v -o /usr/local/bin/app ./relay

CMD ["app"]
```