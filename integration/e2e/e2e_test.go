package e2e

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"syscall"
	"time"

	docker "github.com/fsouza/go-dockerclient"
	"github.com/hyperledger/fabric/integration/nwo"
	"github.com/hyperledger/fabric/integration/nwo/commands"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
	"github.com/tedsuo/ifrit"
)

const LongEventualTimeout = time.Minute

var _ = Describe("e2e", func() {
	var (
		testDir   string
		client    *docker.Client
		network   *nwo.Network
		chaincode nwo.Chaincode
		process   ifrit.Process
	)

	BeforeEach(func() {
		var err error
		testDir, err = ioutil.TempDir("", "e2e")
		Expect(err).NotTo(HaveOccurred())

		client, err = docker.NewClientFromEnv()
		Expect(err).NotTo(HaveOccurred())

		chaincode = nwo.Chaincode{
			Name:    "wasmcc",
			Version: "0.0",
			Path:    "github.com/hyperledger-labs/fabric-chaincode-wasm/wasmcc",
			Ctor:    `{"Args":[]}`,
			Policy:  `AND ('Org1MSP.member')`,
		}
		network = nwo.New(SimpleSoloNetwork(), testDir, client, 30000, components)
		network.GenerateConfigTree()
		network.Bootstrap()

		networkRunner := network.NetworkGroupRunner()
		process = ifrit.Invoke(networkRunner)
		Eventually(process.Ready(), network.EventuallyTimeout).Should(BeClosed())

	})
	AfterEach(func() {
		if process != nil {
			process.Signal(syscall.SIGTERM)
			Eventually(process.Wait(), network.EventuallyTimeout).Should(Receive())
		}
		if network != nil {
			network.Cleanup()
		}
		os.RemoveAll(testDir)
	})
	It("is able to deploy and execute wasmcc contract", func() {
		By("getting the orderer by name")
		orderer := network.Orderer("orderer")

		By("setting up the channel")
		network.CreateAndJoinChannel(orderer, "testchannel")

		By("deploying the chaincode")
		nwo.DeployChaincode(network, "testchannel", orderer, chaincode)

		By("getting the client peer by name")
		peer := network.Peer("Org1", "peer0")

		By("querying installed chaincodes")
		RunQueryInstalledChaincodesEmpty(network, orderer, peer, "testchannel")

		By("installing sample Rust based wasm chaincodes : rust-balancewasm")
		sess, err := network.PeerUserSession(peer, "User1", commands.ChaincodeInvoke{
			ChannelID: "testchannel",
			Orderer:   network.OrdererAddress(orderer, nwo.ListenPort),
			Name:      "wasmcc",
			Ctor:      fmt.Sprintf(`{"Args":["create","rust-balancewasm","%s","account1","100","account2","1000"]}`, ReadAssetTransferWASMHex()),
			PeerAddresses: []string{
				network.PeerAddress(network.Peer("Org1", "peer0"), nwo.ListenPort),
			},
			WaitForEvent: true,
		})
		Expect(err).NotTo(HaveOccurred())
		Eventually(sess, LongEventualTimeout).Should(gexec.Exit(0))

		By("verify rust-balancewasm chaincode installed")
		RunQueryInstalledChaincodesExists(network, orderer, peer, "testchannel")
		By("transferring account balance")
		sess, err = network.PeerUserSession(peer, "User1", commands.ChaincodeInvoke{
			ChannelID: "testchannel",
			Orderer:   network.OrdererAddress(orderer, nwo.ListenPort),
			Name:      "wasmcc",
			Ctor:      fmt.Sprintf(`{"Args":["execute","rust-balancewasm","invoke","account2","account1","10"]}`),
			PeerAddresses: []string{
				network.PeerAddress(network.Peer("Org1", "peer0"), nwo.ListenPort),
			},
			WaitForEvent: true,
		})
		Expect(err).NotTo(HaveOccurred())
		Eventually(sess, LongEventualTimeout).Should(gexec.Exit(0))

		By("querying updated account balance")
		sess, err = network.PeerUserSession(peer, "User1", commands.ChaincodeQuery{
			ChannelID: "testchannel",
			Name:      "wasmcc",
			//get()
			Ctor: fmt.Sprintf(`{"Args":["execute","rust-balancewasm","query","account1"]}`),
		})
		Expect(err).NotTo(HaveOccurred())
		Eventually(sess, LongEventualTimeout).Should(gexec.Exit(0))
		Expect(sess.Out).To(gbytes.Say("110"))

	})
})

func RunQueryInstalledChaincodesExists(n *nwo.Network, orderer *nwo.Orderer, peer *nwo.Peer, channel string) {
	sess, err := n.PeerUserSession(peer, "User1", commands.ChaincodeQuery{
		ChannelID: channel,
		Name:      "wasmcc",
		Ctor:      `{"Args":["installedChaincodes",""]}`,
	})
	Expect(err).NotTo(HaveOccurred())
	Eventually(sess, n.EventuallyTimeout).Should(gexec.Exit(0))
	Expect(sess).To(gbytes.Say("balancewasm\n"))

}

func RunQueryInstalledChaincodesEmpty(n *nwo.Network, orderer *nwo.Orderer, peer *nwo.Peer, channel string) {
	sess, err := n.PeerUserSession(peer, "User1", commands.ChaincodeQuery{
		ChannelID: channel,
		Name:      "wasmcc",
		Ctor:      `{"Args":["installedChaincodes",""]}`,
	})
	Expect(err).NotTo(HaveOccurred())
	Eventually(sess, n.EventuallyTimeout).Should(gexec.Exit(0))
	Expect(sess).To(gbytes.Say(""))

}
func ReadAssetTransferWASMHex() string {

	file, err := ioutil.ReadFile("../../sample-wasm-chaincode/chaincode_example02/rust/app_main.wasm")
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	return hex.EncodeToString(file)
}
