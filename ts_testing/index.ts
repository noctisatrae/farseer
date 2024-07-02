import { CastType, EthersEip712Signer, FarcasterNetwork, NobleEd25519Signer, getFarcasterTime, getInsecureHubRpcClient, makeCastAdd } from "@farcaster/hub-nodejs";
import { Wallet } from "ethers";

const HUB_URL = "38.242.131.38:2283" 
const MNEMONIC = Bun.env.MNEMONIC ?? ""
const NETWORK = FarcasterNetwork.MAINNET
const FID = 10626

const wallet = Wallet.fromPhrase(MNEMONIC)
const signer = new EthersEip712Signer(wallet)
const clt = getInsecureHubRpcClient(HUB_URL);

const castToSend = await makeCastAdd(
  {
    text: "Sent from my custom hub! Yay :)",
    embeds: [],
    embedsDeprecated: [],
    mentions: [],
    mentionsPositions: [],
    type: CastType.CAST
  },
  {
    timestamp: getFarcasterTime()._unsafeUnwrap(),
    fid: FID,
    network: NETWORK
  },
  signer
)

const res = await clt.submitMessage(castToSend._unsafeUnwrap())
console.debug(res.isOk() ? 'Submission was successful!' : `Request failed! Reason=${res.error}`)