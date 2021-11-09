module parallel

go 1.17

require (
	github.com/allaman/toolbox/es-query/auxiliary v0.0.0-00010101000000-000000000000
	github.com/elastic/go-elasticsearch/v7 v7.13.1
	github.com/tidwall/gjson v1.9.1
)

require (
	github.com/tidwall/match v1.0.3 // indirect
	github.com/tidwall/pretty v1.2.0 // indirect
)

replace github.com/allaman/toolbox/es-query/auxiliary => ../auxiliary
