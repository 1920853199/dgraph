package fbx

import (
	"github.com/dgraph-io/dgraph/fb"
	flatbuffers "github.com/google/flatbuffers/go"
)

type Facet struct {
	builder *flatbuffers.Builder

	key       flatbuffers.UOffsetT
	value     flatbuffers.UOffsetT
	valueType fb.FacetValueType
	tokens    flatbuffers.UOffsetT
	alias     flatbuffers.UOffsetT
}

func NewFacet() *Facet {
	return &Facet{
		builder: flatbuffers.NewBuilder(bufSize),
	}
}

func (f *Facet) SetKey(key string) *Facet {
	f.key = f.builder.CreateString(key)
	return f
}

func (f *Facet) SetValue(value []byte) *Facet {
	f.value = f.builder.CreateByteVector(value)
	return f
}

func (f *Facet) SetValueType(valueType fb.FacetValueType) *Facet {
	f.valueType = valueType
	return f
}

func (f *Facet) SetTokens(tokens []string) *Facet {
	offsets := make([]flatbuffers.UOffsetT, len(tokens))
	for i, token := range tokens {
		offsets[i] = f.builder.CreateString(token)
	}

	fb.FacetStartTokensVector(f.builder, len(tokens))
	for i := len(tokens) - 1; i >= 0; i-- {
		f.builder.PrependUOffsetT(offsets[i])
	}
	f.tokens = f.builder.EndVector(len(tokens))

	return f
}

func (f *Facet) SetAlias(alias string) *Facet {
	f.alias = f.builder.CreateString(alias)
	return f
}

func (f *Facet) Build() *fb.Facet {
	facet := f.buildOffset()
	f.builder.Finish(facet)

	buf := f.builder.FinishedBytes()
	return fb.GetRootAsFacet(buf, 0)
}

func (f *Facet) buildOffset() flatbuffers.UOffsetT {
	fb.FacetStart(f.builder)
	fb.FacetAddKey(f.builder, f.key)
	fb.FacetAddValue(f.builder, f.value)
	fb.FacetAddValueType(f.builder, f.valueType)
	fb.FacetAddTokens(f.builder, f.tokens)
	fb.FacetAddAlias(f.builder, f.alias)
	return fb.FacetEnd(f.builder)
}
