package tools

// File contains the basic data about a file.
type File struct {
	Name            string
	Size            int
	PieceSize       int
	Key             string
	BufferMapLength int
	Peers           map[string]*Peer
} //BufferMap BufferMap

// Piece contains the data about a piece (the actual data).
type Piece struct {
	Index int
	Data  Data
}

// Data contains the actual data of a piece.
// The `Length` attribute is the length of the bit sequence, not of the `BitSequence` array, which is padded to the byte.
// Not to be confused with the BufferMap struct which contains the bits telling whether a peer has pieces.
type Data struct {
	Length      int
	BitSequence []byte
}

// ListData is the struct returned by ListCheck.
type ListData struct {
	Files []File
}

// InterestedData is the struct returned by InterestedCheck.
type InterestedData struct {
	Key string
}

// HaveData is the struct returned by HaveCheck.
type HaveData struct {
	Key       string
	BufferMap BufferMap
}

// GetPiecesData is the struct returned by GetPiecesCheck.
type GetPiecesData struct {
	Key    string
	Pieces []int
}

// DataData is the struct returned by DataCheck.
type DataData struct {
	Key    string
	Pieces []Piece
}

// PeersData is the struct returned by PeersCheck.
type PeersData struct {
	Key   string
	Peers []Peer
}

// Peer is a struct used by PeersData.
type Peer struct {
	IP         string
	Port       int
	BufferMaps map[string]*BufferMap
}

// BufferMap tells whether the peer has the pieces of a file.
// The `Length` attribute is the length of the bit sequence, not of the `BitSequence` array, which is padded to the byte.
type BufferMap struct {
	Length      int
	BitSequence []byte
}

// FileCmp tells if two File are equal.
func FileCmp(f1 File, f2 File) bool {
	return f1.Key == f2.Key // && f1.Name == f2.Name && (f1.Size == f2.Size) && (f1.PieceSize == f2.PieceSize) && (&f1.Peers == &f2.Peers) //BufferMapCmp(f1.BufferMap, f2.BufferMap)
}

// ListDataCmp tells if two ListData are equal.
func ListDataCmp(lD1 ListData, lD2 ListData) bool {
	if len(lD1.Files) != len(lD2.Files) {
		return false
	}
	for i, f := range lD1.Files {
		if !FileCmp(f, lD2.Files[i]) {
			return false
		}
	}
	return true
}

// InterestedCmp tells if two InterestedData are equal.
func InterestedCmp(iD1 InterestedData, iD2 InterestedData) bool {
	return iD1.Key == iD2.Key
}

// HaveCmp tells if two HaveData are equal.
func HaveCmp(hD1 HaveData, hD2 HaveData) bool {
	return hD1.Key == hD2.Key && BufferMapCmp(hD1.BufferMap, hD2.BufferMap)
}

// DataStructCmp tells if two Data are equal.
func DataStructCmp(d1 Data, d2 Data) bool {
	if d1.Length != d2.Length {
		return false
	}
	for i, v := range d1.BitSequence {
		if v != d2.BitSequence[i] {
			return false
		}
	}
	return true
}

// PieceCmp tells if two Piece are equal.
func PieceCmp(p1 Piece, p2 Piece) bool {
	return p1.Index == p2.Index && DataStructCmp(p1.Data, p2.Data)
}

// DataCmp tells if two DataData are equal.
func DataCmp(dD1 DataData, dD2 DataData) bool {
	if dD1.Key != dD2.Key || len(dD1.Pieces) != len(dD2.Pieces) {
		return false
	}
	for i, piece := range dD1.Pieces {
		if !PieceCmp(piece, dD2.Pieces[i]) {
			return false
		}
	}
	return true
}

// GetPiecesCmp tells if two GetPiecesData are equal.
func GetPiecesCmp(gPD1 GetPiecesData, gPD2 GetPiecesData) bool {
	if gPD1.Key != gPD2.Key {
		return false
	}
	for i, p := range gPD1.Pieces {
		if p != gPD2.Pieces[i] {
			return false
		}
	}
	return true
}

// BufferMapCmp tells if two BufferMap are equal.
func BufferMapCmp(bM1 BufferMap, bM2 BufferMap) bool {
	if bM1.Length != bM2.Length {
		return false
	}
	if bM1.Length == 0 {
		return true
	}
	for i := 0; i < bM1.Length/8; i++ {
		if bM1.BitSequence[i] != bM2.BitSequence[i] {
			return false
		}
	}
	return true
}

func PeersCmp(pD1 PeersData, pD2 PeersData) bool {
	if pD1.Key != pD2.Key || len(pD1.Peers) != len(pD2.Peers) {
		return false
	}
	for i, peer := range pD1.Peers {
		if !PeerCmp(peer, pD2.Peers[i]) {
			return false
		}
	}
	return true
}

func PeerCmp(p1 Peer, p2 Peer) bool {
	return p1.IP == p2.IP && p1.Port == p2.Port
}

func MapPeersCmp(mp1 map[string]*Peer, mp2 map[string]*Peer) bool {
	for k, arr := range mp1 {
		arr2, err := mp2[k]
		if err {
			return false
		}
		if arr != arr2 {
			return false
		}
	}
	return true
}
