package packets

import (
	_log "log"
	"os"
)

var log = _log.New(
	os.Stdout,
	"loser: ",
	_log.Ldate|
		_log.Ltime|
		_log.Lmicroseconds|
		_log.LUTC|
		_log.Lmsgprefix|
		_log.LstdFlags,
)
