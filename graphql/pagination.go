package graphql

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"

	"github.com/solher/toolbox/graphql/types"
)

// PaginationArgs are standardized pagination arguments.
type PaginationArgs struct {
	First, Last   *int32
	After, Before *types.ID
}

// PaginationArgsToOffsetLimit converts pagination arguments into offset/limit to allow
// easy SQL querying.
func PaginationArgsToOffsetLimit(args PaginationArgs) (offset, limit uint64, err error) {
	if args.After != nil && args.Before != nil {
		return 0, 0, errors.New("cannot paginate forward and backward at the same time")
	}
	limit = 20
	switch {
	case args.After != nil:
		offset = FromCursor(*args.After)
		if args.First != nil {
			if *args.First < 0 {
				return 0, 0, errors.New("the 'first' argument cannot be negative")
			}
			limit = uint64(*args.First)
		}
	case args.Before != nil:
		offset = FromCursor(*args.Before)
		if args.Last != nil {
			if *args.Last < 0 {
				return 0, 0, errors.New("the 'last' argument cannot be negative")
			}
			limit = uint64(*args.Last)
		}
		if limit > offset {
			limit = offset
		}
		offset = offset - limit
	default:
		if args.First != nil {
			limit = uint64(*args.First)
		}
	}
	return offset, limit, nil
}

// FromCursor returns an offset from an encoded cursor.
func FromCursor(cursor types.ID) (offset uint64) {
	offsetStr, err := base64.StdEncoding.DecodeString(string(cursor))
	if err != nil {
		return 0
	}
	offset, err = strconv.ParseUint(string(offsetStr), 10, 64)
	if err != nil {
		return 0
	}
	return offset
}

// ToCursor generates an encoded cursor from an offset.
func ToCursor(offset uint64) (cursor types.ID) {
	str := fmt.Sprintf("%d", offset)
	encStr := base64.StdEncoding.EncodeToString([]byte(str))
	return types.ID(encStr)
}

// NewPageInfoResolver returns an instance of a PageInfoResolver.
func NewPageInfoResolver(offset, limit, totalCount uint64) *PageInfoResolver {
	r := &PageInfoResolver{
		hasPreviousPage: (offset > 0),
		hasNextPage:     (offset+limit < totalCount),
	}
	if r.hasPreviousPage {
		c := ToCursor(offset)
		r.startCursor = &c
	}
	if r.hasNextPage {
		c := ToCursor(offset + limit)
		r.endCursor = &c
	}
	return r
}

// PageInfoResolver resolves page information.
type PageInfoResolver struct {
	hasNextPage, hasPreviousPage bool
	startCursor, endCursor       *types.ID
}

// HasNextPage indicates if a next result page is available.
func (r *PageInfoResolver) HasNextPage() bool {
	return r.hasNextPage
}

// HasPreviousPage indicates if a previous result page is available.
func (r *PageInfoResolver) HasPreviousPage() bool {
	return r.hasPreviousPage
}

// StartCursor returns the cursor indicating the start of the page.
func (r *PageInfoResolver) StartCursor() *types.ID {
	return r.startCursor
}

// EndCursor returns the cursor indicating the end of the page.
func (r *PageInfoResolver) EndCursor() *types.ID {
	return r.endCursor
}
