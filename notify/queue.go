/* 2019-10-15 (cc) <paul4hough@gmail.com>

  send udpates to notify system
store and q note for notify sys

returns note key
*/
package notify
import (
	"bytes"
	"encoding/gob"
	"time"
	pmod "github.com/prometheus/common/model"
	promp "github.com/prometheus/client_golang/prometheus"
	"github.com/boltdb/bolt"
)

/* we only care about the LAST update */
type RetryQueue sync.Map[Key]Note

/* logic
db has nsys state + remed data

spin new tread for external update (> 30 sec)

update retry or attempt action
if ! pass, update retry


done
*/




type QKey string

func NewQKey(sys, grp string, nkey []byte) QKey { return append([]byte(QKey{sys + "-" + grp}),nkey); }

func (qk QKey) Bytes() { return []byte(qk); }
func (qk QKey) String() { return string(qk); }
