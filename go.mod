module github.com/mmb/tmpbbs

go 1.24.1

require (
	github.com/derekparker/trie v0.0.0-20230829180723-39f4de51ef7d
	github.com/google/uuid v1.6.0
	github.com/kyokomi/emoji/v2 v2.2.13
	github.com/microcosm-cc/bluemonday v1.0.27
	github.com/russross/blackfriday/v2 v2.1.0
	github.com/skip2/go-qrcode v0.0.0-20200617195104-da1b6568686e
	github.com/spf13/pflag v1.0.6
	github.com/spf13/viper v1.20.1
	golang.org/x/text v0.23.0
	google.golang.org/grpc v1.71.1
	google.golang.org/protobuf v1.36.6
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/aymerick/douceur v0.2.0 // indirect
	github.com/fsnotify/fsnotify v1.8.0 // indirect
	github.com/go-viper/mapstructure/v2 v2.2.1 // indirect
	github.com/gorilla/css v1.0.1 // indirect
	github.com/pelletier/go-toml/v2 v2.2.3 // indirect
	github.com/sagikazarmark/locafero v0.7.0 // indirect
	github.com/sourcegraph/conc v0.3.0 // indirect
	github.com/spf13/afero v1.12.0 // indirect
	github.com/spf13/cast v1.7.1 // indirect
	github.com/subosito/gotenv v1.6.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/net v0.37.0 // indirect
	golang.org/x/sys v0.31.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250324211829-b45e905df463 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
)

replace github.com/kyokomi/emoji/v2 => github.com/mmb/kyokomiemoji/v2 v2.0.0-20250405050457-8f6070a155ee
