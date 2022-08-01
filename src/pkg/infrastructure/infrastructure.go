package infrastructure

import (
	"log"
	"net/http"

	"github.com/ethereum/go-ethereum/ethclient"
	"gorm.io/gorm"

	"hms/gateway/pkg/compressor"
	"hms/gateway/pkg/config"
	"hms/gateway/pkg/docs/service/processing"
	"hms/gateway/pkg/indexer"
	"hms/gateway/pkg/keystore"
	"hms/gateway/pkg/localDB"
	"hms/gateway/pkg/storage"
	"hms/gateway/pkg/storage/ipfs"
)

type Infra struct {
	LocalDB            *gorm.DB
	Keystore           *keystore.KeyStore
	HttpClient         *http.Client
	EthClient          *ethclient.Client
	IpfsClient         *ipfs.Client
	Index              *indexer.Index
	LocalStorage       storage.Storager
	Compressor         compressor.Interface
	CompressionEnabled bool
}

func New(cfg *config.Config) *Infra {
	sc := storage.NewConfig(cfg.StoragePath)
	storage.Init(sc)

	db, err := localDB.New(cfg.DB.FilePath)
	if err != nil {
		log.Fatal(err)
	}

	db.AutoMigrate(&processing.Request{})
	db.AutoMigrate(&processing.BlockchainTx{})

	ks := keystore.New(cfg.KeystoreKey)

	ehtClient, err := ethclient.Dial(cfg.Contract.Endpoint)
	if err != nil {
		log.Fatal(err)
	}

	ipfsClient, err := ipfs.NewClient(cfg.Storage.Ipfs.EndpointURL)
	if err != nil {
		log.Fatal(err)
	}

	return &Infra{
		LocalDB:            db,
		Keystore:           ks,
		HttpClient:         http.DefaultClient,
		EthClient:          ehtClient,
		IpfsClient:         ipfsClient,
		Index:              indexer.New(cfg.Contract.Address, cfg.Contract.PrivKeyPath, ehtClient),
		LocalStorage:       storage.Storage(),
		Compressor:         compressor.New(cfg.CompressionLevel),
		CompressionEnabled: cfg.CompressionEnabled,
	}
}
