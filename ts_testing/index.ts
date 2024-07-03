import { CastType, FarcasterNetwork, KEY_GATEWAY_ADDRESS, ViemLocalEip712Signer, getFarcasterTime, getInsecureHubRpcClient, keyGatewayABI, makeCastAdd } from "@farcaster/hub-nodejs";
import { createWalletClient, fromBytes, http, publicActions, toHex } from "viem";
import { mnemonicToAccount, toAccount } from "viem/accounts";
import { optimism } from "viem/chains";

const HUB_URL = "38.242.131.38:2283" 
const MNEMONIC = Bun.env.MNEMONIC ?? ""
const OP_ENDPOINT = Bun.env.OP_ENDPOINT;
const FID = 10626
const NETWORK = FarcasterNetwork.MAINNET

const KeyContract = { abi: keyGatewayABI, address: KEY_GATEWAY_ADDRESS, chain: optimism };
const clt = getInsecureHubRpcClient(HUB_URL)
const account = mnemonicToAccount(MNEMONIC)
const localAccount = toAccount(account)
const signer = new ViemLocalEip712Signer(localAccount)

const walletClient = createWalletClient({
  account,
  chain: optimism,
  transport: http(OP_ENDPOINT!),
}).extend(publicActions);

const metadata = (await signer.getSignedKeyRequestMetadata({
  requestFid: BigInt(FID),
  key: account.getHdKey().privateKey!,
  deadline: BigInt(Math.floor(Date.now() / 1000) + 60 * 60),
}))._unsafeUnwrap()

const { request: signerAddRequest } = await walletClient.simulateContract({
  ...KeyContract,
  functionName: 'add',
  args: [1, toHex(account.getHdKey().publicKey!), 1, toHex(metadata)], // keyType, publicKey, metadataType, metadata
});

const signerAddTxHash = await walletClient.writeContract(signerAddRequest);
console.log(`Waiting for signer add tx to confirm: ${signerAddTxHash}`);
await walletClient.waitForTransactionReceipt({ hash: signerAddTxHash });
console.log("Sleeping 30 seconds to allow hubs to pick up the signer tx");

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