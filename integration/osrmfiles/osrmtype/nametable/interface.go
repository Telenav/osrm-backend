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
	Pronuciation string
	Ref          string
	Exits        string
}

// NameIDOffset is implementation convention that how to get next/previous NameID from current one.
// nextNameID = currNameID + NameIDOffset
// previousNameID = currentNameID - NameIDOffset
// In C++ implementation,
// https://github.com/Telenav/osrm-backend/blob/b24b8a085dc10bea279ffb352049330beae23791/src/extractor/extractor_callbacks.cpp#L323,
// https://github.com/Telenav/osrm-backend/blob/b24b8a085dc10bea279ffb352049330beae23791/include/extractor/name_table.hpp#L38
// Above 5 Names items store together in memory, and always use the index of Names.Name as NameID.
// Way string data is stored in blocks based on `id` as follows:
//
// | name | destination | pronunciation | ref | exits
//                      ^               ^
//                      [range)
//                       ^ id + 2
//
// `id + offset` gives us the range of chars.
//
// Offset 0 is name, 1 is destination, 2 is pronunciation, 3 is ref, 4 is exits
const NameIDOffset = 5

// GetNamesForID try to get Names for an NameID.
func (i IndexedData) GetNamesForID(id osrmtype.NameID) (Names, error) {
	if id == osrmtype.InvalidNameID || id%NameIDOffset != 0 {
		return Names{}, fmt.Errorf("invalid name id: %d", id)
	}

	var name Names
	var b []byte
	var err error
	if b, err = i.at(uint32(id)); err != nil {
		return name, err
	}
	name.Name = string(b)

	if b, err = i.at(uint32(id + 1)); err != nil {
		return name, err
	}
	name.Destinations = string(b)

	if b, err = i.at(uint32(id + 2)); err != nil {
		return name, err
	}
	name.Pronuciation = string(b)

	if b, err = i.at(uint32(id + 3)); err != nil {
		return name, err
	}
	name.Ref = string(b)

	if b, err = i.at(uint32(id + 4)); err != nil {
		return name, err
	}
	name.Exits = string(b)

	return name, nil
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
