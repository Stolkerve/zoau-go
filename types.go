package zoau

type BlockInFileArgs struct {
	// The line(s) to insert inside the marker lines separated by '\\n'. (e.g. "line 1\\nline 2\\nline 3")
	Block *string

	// The marker line template in this format <marker_begin>\\n<marker_end>\\n< {mark} marker>
	// The template should be 3 sections separated by '\\n'. (e.g "OPEN\\nCLOSE\\n# {mark} IBM BLOCK")
	// "{mark}" should be included in the < {mark} marker> (default="# {mark} MANAGED BLOCK") section
	// and will be replaced with <marker_begin> (default="BEGIN") and <marker_end> (default="END").
	// The two marker lines will be surrounding the lines that are going to be inserted.
	// Marker lines can only be used once. If marker lines already exist in the target dataset or HFS file,
	// they will be removed with the surrounded lines before new block get inserted
	Market *string

	// Insert block after matching regex pattern
	// The special value "EOF" will insert the block at the end of the target dataset or HFS file
	InsAft *string

	// Insert block before matching regex pattern
	// The special value "BOF" will insert the block at the beginning of the target dataset or HFS file.
	InsBef *string

	// Encoding of the dataset
	Encoding *string

	// Defaults to True.
	//	- state=True -> Insert or replace block
	//	- state=False -> Remove block
	State *bool

	// Obtain exclusive lock for the dataset.
	Lock bool

	// Force open. Open dataset member in DISP=SHR mode. Default is DISP=OLD mode when False.
	Force bool
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
	// If the source data set has aliases, they will be recreated in the target data set.
	Alias bool

	// This should be set if the source data set is an executable.
	Executable bool

	Binary   bool
	TextMode bool

	// Forces the copy. IMPORTANT: Use of this option can lead to permanent loss of the original target information.
	Force bool
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
	// Mutally inclusive with key_offset. Required for KSDS datasets.
	KeyLength uint
	// Mutally inclusive with key_length. Required for KSDS datasets.
	KeyOffset uint
}

type CreateArgs struct {
	// Type of dataset (also known as dsorg).
	Type *DsType

	// Space to allocate for the dataset. Defaults to 5M
	PrimarySpace *string

	// Secondary (extent) space to allocate for the dataset.
	// Defaults to 1/10 of primary space.
	SecondarySpace *string

	// Directory blocks for PDS-type datasets. Default is 5.
	DirectoryBlocks *uint

	// Block size of dataset.
	// Default varies on record format: FBA=32718, FB=32720, VBA=32743, VB/U=32760
	BlockSize *uint

	// Record format of dataset.
	// FB (default), F, FBA, FBS, U, VB, VBA, VBS
	RecordFormat *RcFormat

	// Logical record length, expressed in bytes
	//	- Defaults vary on format. F/FB/FBS=80, FBA=133, VB/VBA/VBS=137, U=0.
	//	- For variable datasets, the length must include the 4-byte prefix area.
	RecordLength *uint

	// The storage class for an SMS-managed dataset.
	//	- Required for SMS-managed datasets that do not match an SMS-rule.
	//	- Not valid for datasets that are not SMS-managed.
	//	- Note that all non-linear VSAM datasets are SMS-managed.
	StorageClassName *string

	// Data class name for dataset.
	DataClassName *string

	// The management class for an SMS-managed dataset.
	//	- Optional for SMS-managed datasets that do not match an SMS-rule.
	//	- Not valid for datasets that are not SMS-managed.
	//	- Note that all non-linear VSAM datasets are SMS-managed.
	ManagementClassName *string

	// Required for KSDS datasets
	Keys *KeyPoint

	// Comma separated list of volume serials. Offline volumes are not considered.
	Volumes *string
}

type LineInFileArgs struct {
	// The regular expression to look for in every line of the dataset or HFS file.
	//	- For state = True, the pattern to replace if found. Only the last line found will be replaced.
	//	- For state = False, the pattern of the line(s) to remove.
	//	- If the regular expression is not matched, the line will be added to the dataset or HFS file in keeping with ins_aft or ins_bef settings.
	Regex *string

	// Insert line after matching regex pattern.
	//	- The special value “EOF” will insert the line at the end of the target dataset or HFS file.
	//	- If regex is provided, ins_aft is only honored if no match for regex is found.
	//	- ins_bef will be ignored if provided.
	InsAft *string

	// Insert line before matching regex pattern
	//	- The special value “BOF” will insert the line at the beginning of the target dataset or HFS file.
	//	- If regex is provided, ins_bef is only honored if no match for regex is found.
	//	- ins_bef will be ignored if ins_aft is provided.
	InsBef *string

	// Encoding of the dataset.
	Encoding *string

	// Defaults to True.
	//	- state=True -> Insert or replace block
	//	- state=False -> Remove block
	State *bool

	// If set, ins_aft and ins_bef will work with the first line that matches the given regular expression.
	FirstMatch bool

	// Obtain exclusive lock for the dataset.
	Lock bool

	// Force open. Open dataset member in DISP=SHR mode. Default is DISP=OLD mode when False.
	Force bool
}

type ListingArgs struct {
	// If True, only Dataset names are populated in the returned object.
	NameOnly bool

	// Display migrated datasets. Requires name_only to be True.
	Migrate bool

	// Filter dataset by volume name.
	Volume *string
}

// Struct that represents the z/OS dataset.
type Dataset struct {
	// Name of the dataset.
	Name string

	// Record format of the dataset.
	LastReferenced string

	// Dataset organization of the dataset.
	Dsorg string

	// Record format of the dataset.
	Recfm string

	// Record length of the dataset.
	Lrecl int

	// Block size of the dataset.
	BlockSize int

	// Volume the dataset resides on.
	Volume string

	// Estimated used space of the dataset. nil if unknown.
	UsedSpace *int

	// Estimated total space of the dataset.
	TotalSpace int
}

type ReadArgs struct {
	// Read the last tail lines from the dataset.
	Tail *uint

	// Returns lines from the given line.
	FromLine *uint
}

type SearchArgs struct {
	// Print only a count of matched lines in the dataset.
	CountLines bool

	// Display the line number for each match.
	DisplayLines bool

	// Ignore case for search.
	IgnoreCase bool

	// Print names of datasets being searched.
	PrintDatasets bool

	// Number of lines to be shown before and after each match.
	Lines *uint
}

type UnZipArgs struct {
	// Src is a dataset
	Dataset bool

	// Overwrite existing data sets with the same name on target device
	Overwrite bool

	// Specifies the SMS classes specified with -S and/or -m should be used when creating temporary datasets.
	SmsForTmp bool

	// Unzip volume (default is dataset).
	Volume *string

	// specify how large to allocate datasets. Valid units are: CYL, TRK, K, M, G. Defaults to bytes if no unit provided.
	Size *string

	// Include particular data set patterns from dzip binary in unzipped contents.
	Include *string

	// Exclude particular data set patterns from dzip binary in unzipped contents.
	Exclude *string

	// specifies the user-desired storage class that is to replace the source storage class as input to the ACS routines.
	StorageClassName *string

	// specifies the user-desired management class that is to replace the source management class as input to the ACS routines.
	ManagementClassName *string

	// specifies a particular volume should be used when creating temporary datasets
	DestVolume *string

	//
	SrcVolume *string
}

type ZipArgs struct {
	// Dump to data set instead of file.
	Dataset bool

	// Overwrite file or dataset destination if it already exists.
	Overwrite bool

	// Specifies potentially recoverable errors should be tolerated.
	Force bool

	// Dump a volume instead of datasets.
	// If a volume is provided along with dataset information and -V is
	// not specified, dzip will look for provided data set patterns on the
	// provided volume.
	Volume *string

	// Specify how large to allocate datasets. Valid units are: CYL, TRK, K, M, G. Defaults to bytes if no unit provided.
	Size *string

	// Exclude pattern for data sets, this option is ignored if dumping a volume.
	Exclude *string

	// Specifies the user-desired storage class is to be used when creating temporary and target datasets.
	StorageClassName *string

	// Specifies the user-desired management class that is to be used when creating temporary and target datasets.
	ManagementClassName *string

	// Specifies a particular volume should be used when creating temporary and target datasets.
	DestVolume *string

	//
	SrcVolume *string
}
