package tsv

type Decoder interface {
	DecodeRecord(s string) error
}
