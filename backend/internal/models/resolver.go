package models

type StreamResolver interface {
	StreamURL(uri string) string
}
