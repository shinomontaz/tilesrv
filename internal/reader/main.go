package reader

type Reader struct {
}

//New to create new reader instance
func New(filename string) *Reader {
	return &Reader{}
}

// Init to read to memory and process provided pbf file, create spatial index
func (r *Reader) Init() {

}
