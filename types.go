package zoau

type BlockInFileArgs struct {
	Block    *string
	Market   *string
	InsAft   *string
	InsBef   *string
	Encoding *string
	State    *bool
	Lock     bool
	Force    bool
}

type Point struct {
	Start uint
	End   uint
}

type CompareArgs struct {
	Columnns   *Point
	Lines      *Point
	IgnoreCase bool
}

type CopyArgs struct {
	Alias      bool
	Executable bool
	Binary     bool
	TextMode   bool
	Force      bool
}

type DsType = string

const (
	DS_ORG_PDS   DsType = "PDS"
	DS_ORG_PDSE  DsType = "PDSE"
	DS_ORG_SEQ   DsType = "SEQ"
	DS_ORG_LDS   DsType = "LDS"
	DS_ORG_RRDS  DsType = "RRDS"
	DS_ORG_ESDS  DsType = "ESDS"
	DS_ORG_KSDS  DsType = "KSDS"
	DS_ORG_LARGE DsType = "LARGE"
)

type RcFormat = string

const (
	RC_FORMAT_FB  RcFormat = "FB"
	RC_FORMAT_FBA RcFormat = "FBA"
	RC_FORMAT_FBS RcFormat = "FBS"
	RC_FORMAT_U   RcFormat = "U"
	RC_FORMAT_VB  RcFormat = "VB"
	RC_FORMAT_VBA RcFormat = "VBA"
	RC_FORMAT_VBS RcFormat = "VBS"
)

type KeyPoint struct {
	KeyLength uint
	KeyOffset uint
}

type CreateArgs struct {
	Type                *DsType
	PrimarySpace        *string
	SecondarySpace      *string
	DirectoryBlocks     *uint
	BlockSize           *uint
	RecordFormat        *RcFormat
	RecordLength        *uint
	StorageClassName    *string
	DataClassName       *string
	ManagementClassName *string
	Keys                *KeyPoint
	Volumes             *string
}

type LineInFileArgs struct {
	Regex      *string
	InsAft     *string
	InsBef     *string
	Encoding   *string
	State      *bool
	FirstMatch bool
	Lock       bool
	Force      bool
	Backref    bool
}

type ListingArgs struct {
	NameOnly bool
	Migrate  bool
	Volume   *string
}

type ListingOutput struct {
	Name           string
	LastReferenced string
	Dsorg          string
	Recfm          string
	Lrecl          int
	BlockSize      int
	Volume         string
	UsedSpace      *int
	TotalSpace     int
}

type ReadArgs struct {
	Tail     *uint
	FromLine *uint
}

type SearchArgs struct {
	CountLines    bool
	DisplayLines  bool
	IgnoreCase    bool
	PrintDatasets bool
	Lines         *uint
}

type UnZipArgs struct {
	Dataset             bool
	Overwrite           bool
	SmsForTmp           bool
	Volume              *string
	Size                *string
	Include             *string
	Exclude             *string
	StorageClassName    *string
	ManagementClassName *string
	DestVolume          *string
	SrcVolume           *string
}

type ZipArgs struct {
	Dataset             bool
	Overwrite           bool
	Force               bool
	Volume              *string
	Size                *string
	Exclude             *string
	StorageClassName    *string
	ManagementClassName *string
	DestVolume          *string
	SrcVolume           *string
}
