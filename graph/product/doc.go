// Copyright ©2019 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package product implements graph product functions.
//
// All the graph products in this package are graphs with order
// n₁n₂ where n₁ and n₂ are the orders of the input graphs. This is
// the order of the set of the Cartesian product of the two input
// graphs' nodes.
//
// The nodes of the product hold the original input graphs' nodes
// in the A and B fields in product.Nodes. This allows a mapping
// between the input graphs and their products.
//
// See https://en.wikipedia.org/wiki/Graph_product for more details
// about graph products.
package product // import "gonum.org/v1/gonum/graph/product"
