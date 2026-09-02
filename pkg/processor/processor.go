package processor

import "github.com/pettai/segment"

type Processor struct {
	le *levenshtein
}

type Line struct {
	Input  []byte
	Tokens []string
	IsNew  bool
}

func NewProcessor(maxDiffRatio float64) *Processor {
	le := newLevenshtein(maxDiffRatio)
	return &Processor{le: le}
}

func (p *Processor) Process(in Line) Line {
	s := segment.NewSegmenterDirect(in.Input)
	in.Tokens = make([]string, 0)
	for s.Segment() {
		t := s.Text()
		// Every type below is variable data, so each collapses to a placeholder
		// and lines differing only in those values compare equal.
		//
		// The tokenizer is a fork of blevesearch/segment that recognizes these
		// shapes as single tokens.
		switch s.Type() {
		case segment.Number:
			t = "<N>"
		case segment.Timestamp:
			t = "<DATETIME>"
		case segment.IPv4:
			t = "<IPV4>"
		case segment.UUID:
			t = "<UUID>"
		case segment.Email:
			t = "<EMAIL>"
		case segment.MAC:
			t = "<MAC>"
		}
		in.Tokens = append(in.Tokens, t)
	}
	in = p.le.process(in)
	return in
}
