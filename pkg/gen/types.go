package gen

import _ "embed"

//go:embed defs/watchlist/watchlist.swagger.json
var WatchlistSwagger string

//go:embed defs/chat/chat.swagger.json
var ChatSwagger string

//go:embed defs/user/user.swagger.json
var UserSwagger string
