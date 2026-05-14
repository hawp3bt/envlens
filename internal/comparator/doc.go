// Package comparator implements multi-environment .env comparison.
//
// Unlike the differ package — which compares exactly two env maps — comparator
// accepts an arbitrary number of labelled environments and produces a unified
// KeyStatus matrix showing presence and value divergence across all of them.
//
// Typical usage:
//
//	result := comparator.Compare(map[string]map[string]string{
//		"dev":     devMap,
//		"staging": stagingMap,
//		"prod":    prodMap,
//	})
//	for _, ks := range result.Statuses {
//		if ks.Diverges {
//			fmt.Println("divergent key:", ks.Key)
//		}
//	}
package comparator
