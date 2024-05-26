https://raw.githubusercontent.com/farcasterxyz/allowlist-mainnet/main/networkConfig.js
This is the list of the peers accepted for Farcaster!

- [X] Find already existing peers id for libp2p
- [X] Find how gossipsub works inside the box

https://github.com/OpenFarcaster/teleport
Another implementation in Rust of a farcaster hub

IDEAS of what I can do with this:
- Sink the data to redis/memecached/any DB => just need to write an adapter
  - <del>REQUIREMENT: allows to write an extandable behaviour so people can build stuff on top of it</del>
  - <del>**RECEIVE MESSAGE => SERIALIZE => GET THE TYPE OF THE MESSAGE => PASS IT TO THE HANDLER** => *HANDLER DO STUFF*</del>
  - HANDLER API IS DONE 

IDEAS OF HANDLERS:
- Apple APNS server who pushes notification when a certain condition on the network is met.

TODO:
- JS/TS SDK
- Grafana dashboard
- Binary/Library?