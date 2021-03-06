package og

// Code generated by github.com/tinylib/msgp DO NOT EDIT.

import (
	"github.com/annchain/OG/types"
	"github.com/tinylib/msgp/msgp"
)

// DecodeMsg implements msgp.Decodable
func (z *MessageCounter) DecodeMsg(dc *msgp.Reader) (err error) {
	var field []byte
	_ = field
	var zb0001 uint32
	zb0001, err = dc.ReadMapHeader()
	if err != nil {
		return
	}
	for zb0001 > 0 {
		zb0001--
		field, err = dc.ReadMapKeyPtr()
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		default:
			err = dc.Skip()
			if err != nil {
				return
			}
		}
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z MessageCounter) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 0
	err = en.Append(0x80)
	if err != nil {
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z MessageCounter) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 0
	o = append(o, 0x80)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *MessageCounter) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var field []byte
	_ = field
	var zb0001 uint32
	zb0001, bts, err = msgp.ReadMapHeaderBytes(bts)
	if err != nil {
		return
	}
	for zb0001 > 0 {
		zb0001--
		field, bts, err = msgp.ReadMapKeyZC(bts)
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		default:
			bts, err = msgp.Skip(bts)
			if err != nil {
				return
			}
		}
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z MessageCounter) Msgsize() (s int) {
	s = 1
	return
}

// DecodeMsg implements msgp.Decodable
func (z *MessageType) DecodeMsg(dc *msgp.Reader) (err error) {
	{
		var zb0001 uint64
		zb0001, err = dc.ReadUint64()
		if err != nil {
			return
		}
		(*z) = MessageType(zb0001)
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z MessageType) EncodeMsg(en *msgp.Writer) (err error) {
	err = en.WriteUint64(uint64(z))
	if err != nil {
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z MessageType) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	o = msgp.AppendUint64(o, uint64(z))
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *MessageType) UnmarshalMsg(bts []byte) (o []byte, err error) {
	{
		var zb0001 uint64
		zb0001, bts, err = msgp.ReadUint64Bytes(bts)
		if err != nil {
			return
		}
		(*z) = MessageType(zb0001)
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z MessageType) Msgsize() (s int) {
	s = msgp.Uint64Size
	return
}

// DecodeMsg implements msgp.Decodable
func (z *P2PMessage) DecodeMsg(dc *msgp.Reader) (err error) {
	var zb0001 uint32
	zb0001, err = dc.ReadArrayHeader()
	if err != nil {
		return
	}
	if zb0001 != 5 {
		err = msgp.ArrayError{Wanted: 5, Got: zb0001}
		return
	}
	{
		var zb0002 uint64
		zb0002, err = dc.ReadUint64()
		if err != nil {
			return
		}
		z.MessageType = MessageType(zb0002)
	}
	z.SourceID, err = dc.ReadString()
	if err != nil {
		return
	}
	z.Version, err = dc.ReadInt()
	if err != nil {
		return
	}
	err = z.Message.DecodeMsg(dc)
	if err != nil {
		return
	}
	if dc.IsNil() {
		err = dc.ReadNil()
		if err != nil {
			return
		}
		z.SourceHash = nil
	} else {
		if z.SourceHash == nil {
			z.SourceHash = new(types.Hash)
		}
		err = z.SourceHash.DecodeMsg(dc)
		if err != nil {
			return
		}
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z *P2PMessage) EncodeMsg(en *msgp.Writer) (err error) {
	// array header, size 5
	err = en.Append(0x95)
	if err != nil {
		return
	}
	err = en.WriteUint64(uint64(z.MessageType))
	if err != nil {
		return
	}
	err = en.WriteString(z.SourceID)
	if err != nil {
		return
	}
	err = en.WriteInt(z.Version)
	if err != nil {
		return
	}
	err = z.Message.EncodeMsg(en)
	if err != nil {
		return
	}
	if z.SourceHash == nil {
		err = en.WriteNil()
		if err != nil {
			return
		}
	} else {
		err = z.SourceHash.EncodeMsg(en)
		if err != nil {
			return
		}
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *P2PMessage) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// array header, size 5
	o = append(o, 0x95)
	o = msgp.AppendUint64(o, uint64(z.MessageType))
	o = msgp.AppendString(o, z.SourceID)
	o = msgp.AppendInt(o, z.Version)
	o, err = z.Message.MarshalMsg(o)
	if err != nil {
		return
	}
	if z.SourceHash == nil {
		o = msgp.AppendNil(o)
	} else {
		o, err = z.SourceHash.MarshalMsg(o)
		if err != nil {
			return
		}
	}
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *P2PMessage) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var zb0001 uint32
	zb0001, bts, err = msgp.ReadArrayHeaderBytes(bts)
	if err != nil {
		return
	}
	if zb0001 != 5 {
		err = msgp.ArrayError{Wanted: 5, Got: zb0001}
		return
	}
	{
		var zb0002 uint64
		zb0002, bts, err = msgp.ReadUint64Bytes(bts)
		if err != nil {
			return
		}
		z.MessageType = MessageType(zb0002)
	}
	z.SourceID, bts, err = msgp.ReadStringBytes(bts)
	if err != nil {
		return
	}
	z.Version, bts, err = msgp.ReadIntBytes(bts)
	if err != nil {
		return
	}
	bts, err = z.Message.UnmarshalMsg(bts)
	if err != nil {
		return
	}
	if msgp.IsNil(bts) {
		bts, err = msgp.ReadNilBytes(bts)
		if err != nil {
			return
		}
		z.SourceHash = nil
	} else {
		if z.SourceHash == nil {
			z.SourceHash = new(types.Hash)
		}
		bts, err = z.SourceHash.UnmarshalMsg(bts)
		if err != nil {
			return
		}
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z *P2PMessage) Msgsize() (s int) {
	s = 1 + msgp.Uint64Size + msgp.StringPrefixSize + len(z.SourceID) + msgp.IntSize + z.Message.Msgsize()
	if z.SourceHash == nil {
		s += msgp.NilSize
	} else {
		s += z.SourceHash.Msgsize()
	}
	return
}

// DecodeMsg implements msgp.Decodable
func (z *SendingType) DecodeMsg(dc *msgp.Reader) (err error) {
	{
		var zb0001 uint8
		zb0001, err = dc.ReadUint8()
		if err != nil {
			return
		}
		(*z) = SendingType(zb0001)
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z SendingType) EncodeMsg(en *msgp.Writer) (err error) {
	err = en.WriteUint8(uint8(z))
	if err != nil {
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z SendingType) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	o = msgp.AppendUint8(o, uint8(z))
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *SendingType) UnmarshalMsg(bts []byte) (o []byte, err error) {
	{
		var zb0001 uint8
		zb0001, bts, err = msgp.ReadUint8Bytes(bts)
		if err != nil {
			return
		}
		(*z) = SendingType(zb0001)
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z SendingType) Msgsize() (s int) {
	s = msgp.Uint8Size
	return
}

// DecodeMsg implements msgp.Decodable
func (z *StatusData) DecodeMsg(dc *msgp.Reader) (err error) {
	var field []byte
	_ = field
	var zb0001 uint32
	zb0001, err = dc.ReadMapHeader()
	if err != nil {
		return
	}
	for zb0001 > 0 {
		zb0001--
		field, err = dc.ReadMapKeyPtr()
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "ProtocolVersion":
			z.ProtocolVersion, err = dc.ReadUint32()
			if err != nil {
				return
			}
		case "NetworkId":
			z.NetworkId, err = dc.ReadUint64()
			if err != nil {
				return
			}
		case "CurrentBlock":
			err = z.CurrentBlock.DecodeMsg(dc)
			if err != nil {
				return
			}
		case "GenesisBlock":
			err = z.GenesisBlock.DecodeMsg(dc)
			if err != nil {
				return
			}
		case "CurrentId":
			z.CurrentId, err = dc.ReadUint64()
			if err != nil {
				return
			}
		default:
			err = dc.Skip()
			if err != nil {
				return
			}
		}
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z *StatusData) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 5
	// write "ProtocolVersion"
	err = en.Append(0x85, 0xaf, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e)
	if err != nil {
		return
	}
	err = en.WriteUint32(z.ProtocolVersion)
	if err != nil {
		return
	}
	// write "NetworkId"
	err = en.Append(0xa9, 0x4e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x49, 0x64)
	if err != nil {
		return
	}
	err = en.WriteUint64(z.NetworkId)
	if err != nil {
		return
	}
	// write "CurrentBlock"
	err = en.Append(0xac, 0x43, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x42, 0x6c, 0x6f, 0x63, 0x6b)
	if err != nil {
		return
	}
	err = z.CurrentBlock.EncodeMsg(en)
	if err != nil {
		return
	}
	// write "GenesisBlock"
	err = en.Append(0xac, 0x47, 0x65, 0x6e, 0x65, 0x73, 0x69, 0x73, 0x42, 0x6c, 0x6f, 0x63, 0x6b)
	if err != nil {
		return
	}
	err = z.GenesisBlock.EncodeMsg(en)
	if err != nil {
		return
	}
	// write "CurrentId"
	err = en.Append(0xa9, 0x43, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x49, 0x64)
	if err != nil {
		return
	}
	err = en.WriteUint64(z.CurrentId)
	if err != nil {
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *StatusData) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 5
	// string "ProtocolVersion"
	o = append(o, 0x85, 0xaf, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e)
	o = msgp.AppendUint32(o, z.ProtocolVersion)
	// string "NetworkId"
	o = append(o, 0xa9, 0x4e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x49, 0x64)
	o = msgp.AppendUint64(o, z.NetworkId)
	// string "CurrentBlock"
	o = append(o, 0xac, 0x43, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x42, 0x6c, 0x6f, 0x63, 0x6b)
	o, err = z.CurrentBlock.MarshalMsg(o)
	if err != nil {
		return
	}
	// string "GenesisBlock"
	o = append(o, 0xac, 0x47, 0x65, 0x6e, 0x65, 0x73, 0x69, 0x73, 0x42, 0x6c, 0x6f, 0x63, 0x6b)
	o, err = z.GenesisBlock.MarshalMsg(o)
	if err != nil {
		return
	}
	// string "CurrentId"
	o = append(o, 0xa9, 0x43, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x49, 0x64)
	o = msgp.AppendUint64(o, z.CurrentId)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *StatusData) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var field []byte
	_ = field
	var zb0001 uint32
	zb0001, bts, err = msgp.ReadMapHeaderBytes(bts)
	if err != nil {
		return
	}
	for zb0001 > 0 {
		zb0001--
		field, bts, err = msgp.ReadMapKeyZC(bts)
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "ProtocolVersion":
			z.ProtocolVersion, bts, err = msgp.ReadUint32Bytes(bts)
			if err != nil {
				return
			}
		case "NetworkId":
			z.NetworkId, bts, err = msgp.ReadUint64Bytes(bts)
			if err != nil {
				return
			}
		case "CurrentBlock":
			bts, err = z.CurrentBlock.UnmarshalMsg(bts)
			if err != nil {
				return
			}
		case "GenesisBlock":
			bts, err = z.GenesisBlock.UnmarshalMsg(bts)
			if err != nil {
				return
			}
		case "CurrentId":
			z.CurrentId, bts, err = msgp.ReadUint64Bytes(bts)
			if err != nil {
				return
			}
		default:
			bts, err = msgp.Skip(bts)
			if err != nil {
				return
			}
		}
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z *StatusData) Msgsize() (s int) {
	s = 1 + 16 + msgp.Uint32Size + 10 + msgp.Uint64Size + 13 + z.CurrentBlock.Msgsize() + 13 + z.GenesisBlock.Msgsize() + 10 + msgp.Uint64Size
	return
}
