package cmd

import (
	//"crypto/ecdsa"
	"errors"
	"flag"
	"fmt"
	"math/big"
	"os"
	"os/signal"
	"strings"

	ecdsa "github.com/core-coin/go-goldilocks"

	//"github.com/ethereum/go-ethereum/crypto"
	"github.com/core-coin/go-core/crypto"

	"github.com/todesstille/eth-faucet/internal/chain"
	"github.com/todesstille/eth-faucet/internal/server"
)

var (
	appVersion = "v1.1.0"
	chainIDMap = map[string]int{"ropsten": 3, "rinkeby": 4, "goerli": 5, "kovan": 42}

	httpPortFlag = flag.Int("httpport", 8080, "Listener port to serve HTTP connection")
	proxyCntFlag = flag.Int("proxycount", 0, "Count of reverse proxies in front of the server")
	queueCapFlag = flag.Int("queuecap", 100, "Maximum transactions waiting to be sent")
	versionFlag  = flag.Bool("version", false, "Print version number")

	payoutFlag   = flag.Int("faucet.amount", 1, "Number of Ethers to transfer per user request")
	intervalFlag = flag.Int("faucet.minutes", 1440, "Number of minutes to wait between funding rounds")
	netnameFlag  = flag.String("faucet.name", "testnet", "Network name to display on the frontend")

	keyJSONFlag  = flag.String("wallet.keyjson", os.Getenv("KEYSTORE"), "Keystore file to fund user requests with")
	keyPassFlag  = flag.String("wallet.keypass", "password.txt", "Passphrase text file to decrypt keystore")
	privKeyFlag  = flag.String("wallet.privkey", os.Getenv("PRIVATE_KEY"), "Private key hex to fund user requests with")
	providerFlag = flag.String("wallet.provider", os.Getenv("WEB3_PROVIDER"), "Endpoint for Ethereum JSON-RPC connection")
)

func init() {
	flag.Parse()
	if *versionFlag {
		fmt.Println(appVersion)
		os.Exit(0)
	}
}

func Execute() {
	privateKey, err := getPrivateKeyFromFlags()
	if err != nil {
		panic(fmt.Errorf("failed to read private key: %w", err))
	}

	var chainID *big.Int //3
	if value, ok := chainIDMap[strings.ToLower(*netnameFlag)]; ok {
		chainID = big.NewInt(int64(value))
	}
	fmt.Println(chainID)
	// chainID := big.NewInt(3)

	txBuilder, err := chain.NewTxBuilder(*providerFlag, privateKey, chainID)
	if err != nil {
		panic(fmt.Errorf("cannot connect to web3 provider: %w", err))
	}
	config := server.NewConfig(*netnameFlag, *httpPortFlag, *intervalFlag, *payoutFlag, *proxyCntFlag, *queueCapFlag)
	go server.NewServer(txBuilder, config).Run()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}

func getPrivateKeyFromFlags() (*ecdsa.PrivateKey, error) {
	if *privKeyFlag != "" {
		return crypto.HexToEDDSA(*privKeyFlag)
	} else { // if *keyJSONFlag == ""
		return nil, errors.New("missing private key")
	}
	/*
		keyfile, err := chain.ResolveKeyfilePath(*keyJSONFlag)
		if err != nil {
			return nil, err
		}
		password, err := os.ReadFile(*keyPassFlag)
		if err != nil {
			return nil, err
		}

		return chain.DecryptKeyfile(keyfile, strings.TrimRight(string(password), "\r\n"))
	*/
}
