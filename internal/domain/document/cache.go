package document

type Cache interface {
	Get(id string) (*Document, bool)
	Set(id string, doc *Document)
	Invalidate(id string)
}
