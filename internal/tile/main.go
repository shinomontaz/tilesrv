package tile

type IOsmData interface {
}

type Server struct {
	prefix string
	data   IOsmData
}

// New to create new tile server instance
func New(prefix string, reader IOsmData) *Server {
	return &Server{prefix: prefix}
}
