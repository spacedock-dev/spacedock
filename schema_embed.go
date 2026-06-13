// ABOUTME: Embeds the canonical entity mdschema (the frontmatter SSOT) into the
// ABOUTME: binary so field-conformance validation has no runtime docs/ dependency.
package spacedock

import _ "embed"

// EntityMDSchema is the bytes of docs/schema/entity.mdschema.yml, the single
// source of truth for the entity frontmatter contract. It is embedded here at
// the module root because go:embed cannot reach a file above the embedding
// package's directory; embedding the canonical file (rather than a copy) keeps
// one source of truth that the validator reads at run time.
//
//go:embed docs/schema/entity.mdschema.yml
var EntityMDSchema []byte
