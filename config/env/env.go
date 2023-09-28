package env

import (
	"fmt"
	"os"
	"time"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
)

var Envs *Env

type Env struct {
	Port                        int           `env:"PORT" envDefault:"8080"`
	AppEnv                      string        `env:"APP_ENV" envDefault:"development"`
	MisesTestEndpoint           string        `env:"MISES_TEST_ENDPOINT" envDefault:""`
	LogLevel                    string        `env:"LOG_LEVEL" envDefault:"INFO"`
	MongoURI                    string        `env:"MONGO_URI" envDefault:"mongodb://localhost:27017"`
	DBUser                      string        `env:"DB_USER"`
	DBPass                      string        `env:"DB_PASS"`
	DBName                      string        `env:"DB_NAME" envDefault:"mises"`
	AssetHost                   string        `env:"ASSET_HOST" envDefault:"http://localhost/"`
	StorageHost                 string        `env:"STORAGE_HOST" envDefault:"http://localhost/"`
	StorageKey                  string        `env:"STORAGE_KEY" envDefault:""`
	StorageSalt                 string        `env:"STORAGE_SALT" envDefault:""`
	StorageProvider             string        `env:"STORAGE_PROVIDER" envDefault:"local"`
	JWTSecret                   string        `env:"JWT_SECRET" envDefault:"jwt secret"`
	TokenDuration               time.Duration `env:"TOKEN_DURATION" envDefault:"720h"`
	AllowOrigins                string        `env:"ALLOW_ORIGINS" envDefault:""`
	MisesEndpoint               string        `env:"MISES_ENDPOINT" envDefault:""`
	MisesChainID                string        `env:"MISES_CHAIN_ID" envDefault:""`
	DebugMisesPrefix            string        `env:"DEBUG_MISES_PREFIX" envDefault:""`
	DebugAirdropPrefix          string        `env:"DEBUG_AIRDROP_PREFIX" envDefault:""`
	GOTWI_API_KEY               string        `env:"GOTWI_API_KEY" envDefault:""`
	GOTWI_API_KEY_SECRET        string        `env:"GOTWI_API_KEY_SECRET" envDefault:""`
	GooglePlayAppID             string        `env:"GOOGLE_PLAY_APPID" envDefault:""`
	AppStoreID                  string        `env:"AppStoreID" envDefault:""`
	OpenseaApiKey               string        `env:"OpenseaApiKey" envDefault:""`
	TwitterAuthSuccessCallback  string        `env:"TwitterAuthSuccessCallback" envDefault:""`
	TWEET_TAG                   string        `env:"TWEET_TAG"`
	ReplyText                   string        `env:"ReplyText" envDefault:"@Mises001"`
	VALID_TWITTER_REGISTER_DATE string        `env:"VALID_TWITTER_REGISTER_DATE"`
	MinCheckFollowers           uint64        `env:"MinCheckFollowers" envDefault:"350"`
	MaxCheckFollowers           uint64        `env:"MaxCheckFollowers" envDefault:"10000"`
	RootPath                    string
}

func init() {
	fmt.Println("socialsvc env initializing...")
	//_, b, _, _ := runtime.Caller(0)
	appEnv := os.Getenv("APP_ENV")
	projectRootPath, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	envPath := projectRootPath + "/.env"
	appEnvPath := envPath + "." + appEnv
	localEnvPath := appEnvPath + ".local"
	_ = godotenv.Load(filtePath(localEnvPath, appEnvPath, envPath)...)
	Envs = &Env{}
	err = env.Parse(Envs)
	if err != nil {
		panic(err)
	}
	Envs.RootPath = projectRootPath
	fmt.Println("socialsvc env root " + projectRootPath)
	fmt.Println("socialsvc env chain id " + Envs.MisesChainID)
	fmt.Println("socialsvc env debug prefix " + Envs.DebugMisesPrefix)
	fmt.Println("socialsvc env loaded...")
}

func filtePath(paths ...string) []string {
	result := make([]string, 0)
	for _, path := range paths {
		if _, err := os.Stat(path); err == nil {
			result = append(result, path)
		}
	}
	return result
}
