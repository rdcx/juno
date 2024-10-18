package transaction

func ContainsKey(keys []TransactionKey, key TransactionKey) bool {
	for _, k := range keys {
		if k == key {
			return true
		}
	}
	return false
}
