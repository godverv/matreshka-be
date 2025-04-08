package matreshka

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.redsock.ru/evon"
	"gopkg.in/yaml.v3"

	"go.vervstack.ru/matreshka/server"
)

func Test_Servers(t *testing.T) {
	t.Run("YAML", func(t *testing.T) {
		t.Run("Marshal_OK", func(t *testing.T) {
			t.Parallel()

			var cfgIn AppConfig
			cfgIn.Servers = getConfigServersFull()

			marshaled, err := cfgIn.Marshal()
			require.NoError(t, err)

			var actual map[any]any

			require.NoError(t, yaml.Unmarshal(marshaled, &actual))

			expected := map[any]any{
				"servers": map[any]any{
					8080: map[string]any{
						"/{FS}": map[string]any{
							"dist": "web/dist",
						},
						"name": "MASTER",
					},
					50051: map[string]any{
						"/{GRPC}": map[string]any{
							"module":  "pkg/matreshka_be_api",
							"gateway": "/api",
						},
						"name": "MASTER2",
					},
				},
			}

			require.Equal(t, expected, actual)
		})
		t.Run("Marshal_With_Name_OK", func(t *testing.T) {
			t.Parallel()

			var cfgIn AppConfig
			cfgIn.Servers = getConfigServersFull()
			cfgIn.Servers[8080].Name = "Main"
			cfgIn.Servers[50051].Name = "Grpc"

			marshaled, err := cfgIn.Marshal()
			require.NoError(t, err)

			var actual map[any]any

			require.NoError(t, yaml.Unmarshal(marshaled, &actual))

			expected := map[any]any{
				"servers": map[any]any{
					8080: map[string]any{
						"name": "Main",
						"/{FS}": map[string]any{
							"dist": "web/dist",
						},
					},
					50051: map[string]any{
						"name": "Grpc",
						"/{GRPC}": map[string]any{
							"module":  "pkg/matreshka_be_api",
							"gateway": "/api",
						},
					},
				},
			}

			require.Equal(t, expected, actual)
		})
		t.Run("Unmarshal_OK", func(t *testing.T) {
			t.Parallel()

			cfg, err := ParseConfig(apiConfig)
			require.NoError(t, err)

			servers := getConfigServersFull()
			require.Equal(t, cfg.Servers, servers)
		})
		t.Run("Unmarshal_With_Name_OK", func(t *testing.T) {
			t.Parallel()

			cfg, err := ParseConfig(apiConfigWithName)
			require.NoError(t, err)

			servers := getConfigServersFull()
			servers[8080].Name = "Main"
			require.Equal(t, cfg.Servers, servers)
		})

		t.Run("Unmarshal_Error", func(t *testing.T) {
			t.Run("Port_Is_Not_Int", func(t *testing.T) {
				cfg := AppConfig{}
				err := cfg.Unmarshal(apiInvalidPortConfig)
				require.Equal(t, err.Error(), "strconv.Atoi: parsing \"string\": invalid syntax\nerror converting port to int\n")
			})
			t.Run("Invalid_Struct", func(t *testing.T) {
				cfg := AppConfig{}
				err := cfg.Unmarshal(apiInvalidStructConfig)
				require.Equal(t, err.Error(), "yaml: unmarshal errors:\n  line 2: cannot unmarshal !!seq into map[string]yaml.Node\nerror unmarshalling to yaml.Nodes\n")
			})
			t.Run("Invalid_Item", func(t *testing.T) {
				cfg := AppConfig{}
				err := cfg.Unmarshal(apiInvalidItemConfig)
				require.Equal(t, err.Error(), "yaml: unmarshal errors:\n  line 3: cannot unmarshal !!seq into map[string]yaml.Node\nerror unmarshaling YAML\nerror decoding server\n")
			})
		})
	})

	t.Run("ENV", func(t *testing.T) {
		t.Run("Marshal", func(t *testing.T) {
			t.Parallel()

			var cfgIn AppConfig
			cfgIn.Servers = getConfigServersFull()

			marshaledNodes, err := cfgIn.Servers.MarshalEnv("MATRESHKA_SERVERS")
			require.NoError(t, err)

			marshalledBytes := evon.Marshal(marshaledNodes)
			require.Equal(t, string(apiEnvConfig), string(marshalledBytes))
		})
		t.Run("Unmarshal", func(t *testing.T) {
			t.Parallel()
			cfg := NewEmptyConfig()
			err := evon.UnmarshalWithPrefix("MATRESHKA", apiEnvConfig, &cfg)
			require.NoError(t, err)

			servers := getConfigServersFull()
			require.Equal(t, cfg.Servers, servers)
		})
	})

}

func getConfigServersFull() Servers {
	return Servers{
		8080: {
			Name: "MASTER",
			Port: "8080",
			GRPC: map[string]*server.GRPC{},
			FS: map[string]*server.FS{
				"/{FS}": {
					Dist: "web/dist",
				},
			},
			HTTP: map[string]*server.HTTP{},
		},
		50051: {
			Name: "MASTER2",
			Port: "50051",
			GRPC: map[string]*server.GRPC{
				"/{GRPC}": {
					Module:  "pkg/matreshka_be_api",
					Gateway: "/api",
				},
			},
			FS:   map[string]*server.FS{},
			HTTP: map[string]*server.HTTP{},
		},
	}
}
