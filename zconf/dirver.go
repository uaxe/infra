package zconf

const (
	YamlName = "yaml"
	JsonName = "json"
	TomlName = "toml"
)

type Driver interface {
	Name() string
	Marshal(any) ([]byte, error)
	Unmarshal([]byte, any) error
}
