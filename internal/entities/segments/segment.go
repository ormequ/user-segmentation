package segments

import "errors"

var (
	ErrEmptySlug  = errors.New("slug cannot be empty")
	ErrSlugToLong = errors.New("slug is too long")
)

type Segment struct {
	Slug string
}

func New(slug string) (Segment, error) {
	if len(slug) == 0 {
		return Segment{}, ErrEmptySlug
	}
	if len(slug) > 255 {
		return Segment{}, ErrSlugToLong
	}
	return Segment{
		Slug: slug,
	}, nil
}
