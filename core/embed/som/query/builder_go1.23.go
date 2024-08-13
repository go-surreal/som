//go:build go1.23

package query

//func (b Builder[M, C]) Iterate(ctx context.Context, chunkSize int) iter.Seq2[*M, error] {
//	return func(yield func(*M, error) bool) {
//		initialOffset := b.query.Offset
//		limit := b.query.Limit
//
//		for page := 0; ; page++ {
//			currentOffset := chunkSize * page
//
//			if limit > 0 && currentOffset > limit {
//				currentOffset = limit
//			}
//
//			b.query.Offset = initialOffset + currentOffset
//			b.query.Limit = chunkSize
//
//			results, err := b.All(ctx)
//			if err != nil {
//				yield(nil, err)
//				return
//			}
//
//			if len(results) < 1 {
//				return
//			}
//
//			for _, result := range results {
//				yield(result, nil)
//			}
//
//			if len(results) < chunkSize || (limit > 0 && currentOffset >= limit) {
//				return
//			}
//		}
//	}
//}
