// Package nametable stores names and provides interface to query names for NameID.
// In OSRM's C++ implementation, it uses a `IndexedDataImpl` as internal storage to reduce memory usage.
// This package follows the storage policy too, and implemented same query method.
// C++ implementation: https://github.com/Telenav/osrm-backend/blob/b24b8a085dc10bea279ffb352049330beae23791/include/extractor/name_table.hpp#L116
package nametable

import (
	"bytes"
	"fmt"

	"github.com/Telenav/osrm-backend/integration/osrmfiles/meta"
	"github.com/Telenav/osrm-backend/integration/osrmfiles/osrmtype"
)

// IndexedData stores indexed(something like compressed) names data.
// It stores data same with its C++ implementation, but only provides read operation.
// C++ implementation: https://github.com/Telenav/osrm-backend/blob/b24b8a085dc10bea279ffb352049330beae23791/include/util/indexed_data.hpp#L262
type IndexedData struct {
	BlocksMeta   meta.Num
	ValuesMeta   meta.Num
	BlocksBuffer bytes.Buffer
	ValuesBuffer bytes.Buffer

	// these two are compact storage that not for human read, so we choose not to expose them.
	blocks []blockReference
	values []byte
}

// Names contains name related stuff that OSRM uses.
// OSRM retrieves them from OSM at
// 	https://github.com/Telenav/osrm-backend/blob/b24b8a085dc10bea279ffb352049330beae23791/profiles/lib/way_handlers.lua#L32
// 	https://github.com/Telenav/osrm-backend/blob/b24b8a085dc10bea279ffb352049330beae23791/profiles/lib/way_handlers.lua#L112
type Names struct {
	Name         string
	Destinations string
	Exits        string
	Ref          string
	Pronuciation string
}

// GetNamesForID try to get Names for an NameID.
func (i IndexedData) GetNamesForID(id osrmtype.NameID) (Names, error) {
	//TODO:
	return Names{}, nil
}

// Validate checks whether IndexedData valid or not after assemble.
func (i IndexedData) Validate() error {

	if uint64(i.BlocksMeta) != uint64(len(i.blocks)) {
		return fmt.Errorf("blocks meta not match, count in meta %d, but actual blocks count %d", i.BlocksMeta, len(i.blocks))
	}
	if uint64(i.ValuesMeta) != uint64(len(i.values)) {
		return fmt.Errorf("values meta not match, count in meta %d, but actual values count %d", i.ValuesMeta, len(i.values))
	}

	return nil
}

// Assemble process the loaded IndexedData to easy to Get Names for NameID.
func (i *IndexedData) Assemble() error {

	i.values = i.ValuesBuffer.Bytes()
	i.ValuesBuffer = bytes.Buffer{}

	return i.assembleBlockReferences()
}
