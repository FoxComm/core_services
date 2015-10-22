package backuper

type Strategy func(settings Settings) error

var Strategies = map[string]Strategy{
	"core_database":     CoreDatabase,
	"feature_databases": FeatureDatabases,
	"origin_database":   OriginDatabase,
	"assets":            Assets,
}
