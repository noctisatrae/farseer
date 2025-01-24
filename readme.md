# farseer - another kind of Farcaster hub
farseer is a lightweight re-implementation of a Farcaster hub that **does not require syncing to operate**. In short, you can bring your own DB, logic & infrastructure and harvest data from the protocol while fostering decentralization! If you like this project, **consider giving a star to the repository**: it's a real help for motivation! Looking for PRs, issues & feedback, so feel free to write something and send it to me! 

See the [todos](./todos.md) to see what's left to do :)

## How to get started?
### General architecture
There are three components to farseer: a config file (`config.toml`), the identity file & the hub itself. Here's a quick rundown of how it works:
```
.
├── config.toml <== configure the behaviour of the hubs & the plugin
├── docker-compose.yml <== infrastructure example
├── Dockerfile <== automatization of the process
├── hub_identity <== SECRET private key of the hub (needs to be generated for docker-compose)
├── relay <== what you'll run (binary)
```
### Easy mode (Docker)
1. Generate a `hub_identity` using the latest utility found in the [release section](https://github.com/noctisatrae/farseer/releases) or run the code in the `identity` folder. **Don't forget to put in the root of the repository!**
2. Change your public IP address in the `config.toml` file so other peers can connect to you!
3. Run this command to start the containers!
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

2. Now, you'll start the hub by running: 
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
```