package chain

import (
	"github.com/ava-labs/avalanchego/database"
)

var (
	_ UnsignedTransaction = &SetTx{}
)

type SetTx struct {
	*BaseTx `serialize:"true"`
	Key     []byte `serialize:"true"`
	Value   []byte `serialize:"true"`
}

func (s *SetTx) Verify(db database.Database, blockTime int64) error {
	if len(s.Key) == 0 {
		return ErrKeyEmpty
	}
	if len(s.Key) > maxKeyLength {
		return ErrKeyTooBig
	}
	if len(s.Value) > maxValueLength {
		return ErrValueTooBig
	}
	i, has, err := GetPrefixInfo(db, s.Prefix)
	if err != nil {
		return err
	}
	// Cannot set key if prefix doesn't exist
	if !has {
		return ErrPrefixMissing
	}
	// Prefix cannot be updated if not owned by modifier
	if i.Owner != s.Sender.Address() {
		return ErrUnauthorized
	}
	// Prefix cannot be updated if expired
	if i.Expiry < blockTime {
		return ErrPrefixExpired
	}
	// If we are trying to delete a key, make sure it previously exists.
	if len(s.Value) > 0 {
		return s.accept(db, blockTime)
	}
	has, err = HasPrefixKey(db, s.Prefix, s.Key)
	if err != nil {
		return err
	}
	// Cannot delete non-existent key
	if !has {
		return ErrKeyMissing
	}
	return s.accept(db, blockTime)
}

func (s *SetTx) accept(db database.Database, blockTime int64) error {
	i, _, err := GetPrefixInfo(db, s.Prefix)
	if err != nil {
		return err
	}
	timeRemaining := (i.Expiry - i.LastUpdated) * i.Keys
	if len(s.Value) == 0 {
		i.Keys--
		if err := DeletePrefixKey(db, s.Prefix, s.Key); err != nil {
			return err
		}
	} else {
		i.Keys++
		if err := PutPrefixKey(db, s.Prefix, s.Key, s.Value); err != nil {
			return err
		}
	}
	newTimeRemaining := timeRemaining / i.Keys
	i.LastUpdated = blockTime
	i.Expiry = blockTime + newTimeRemaining
	return PutPrefixInfo(db, s.Prefix, i)
}
