package conf

const (
	YAML_NAME = "yaml"
	JSON_NAME = "json"
	TOML_NAME = "toml"
)

type Driver interface {
	Name() string
	Marshal(any) ([]byte, error)
	Unmarshal([]byte, any) error
}
