package qiniu

// Dntry qiniu store api URI for point to bucket and object path
type Dntry struct {
	Bucket string
	Key    string
}

// ToString Dntry return a string like 'Bucket:Key'
func (d *Dntry) ToString() string {
	return d.Bucket + ":" + d.Key
}

// EncodedEntryURI Dntry return a string like safeBase64Encode('Bucket:Key')
func (d *Dntry) EncodedEntryURI() string {
	return safeBase64Encode([]byte(d.ToString()))
}
