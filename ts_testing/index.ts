import { CastType, FarcasterNetwork, NobleEd25519Signer, getInsecureHubRpcClient, makeCastAdd } from "@farcaster/hub-nodejs";
import { mnemonicToAccount, signMessage } from "viem/accounts";

const HUB_URL = "localhost:2283" 
const MNEMONIC = Bun.env.MNEMONIC ?? ""
const NETWORK = FarcasterNetwork.MAINNET
const FID = 10626

const account = mnemonicToAccount(MNEMONIC)
const clt = getInsecureHubRpcClient(HUB_URL);
const privateKey = account.getHdKey().privateKey ?? new Uint8Array()
const signer = new NobleEd25519Signer(privateKey)

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
    fid: FID,
    network: NETWORK
  },
  signer
)

clt.submitMessage(castToSend._unsafeUnwrap())