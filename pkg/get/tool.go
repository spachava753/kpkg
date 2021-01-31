package get

// UrlConstructor is an interface that will construct a url to download files from
type UrlConstructor interface {
	Construct() (string, error)
}
