package zoau

import (
	"errors"
	"fmt"
)

func Apf(args ApfArgs) (string, int, error) {
	options := make([]string, 0)
	persistentOption := make([]string, 0)
	if args.Persistent != nil {
		if args.Persistent.AddDs == nil && args.Persistent.DelDs == nil {
			return "", -1, errors.New("AddDs and/or DelDs is required with persistent option")
		}
		if args.Persistent.Market != nil {
			persistentOption = append(persistentOption, "-M", *args.Persistent.Market)
		}
		if args.Persistent.AddDs != nil {
			persistentOption = append(persistentOption, "-P", *args.Persistent.AddDs)
		}
		if args.Persistent.DelDs != nil {
			persistentOption = append(persistentOption, "-R", *args.Persistent.DelDs)
		}
	}
	if args.Opt != nil {
		switch *args.Opt {
		case OPT_ADD | OPT_DEL:
			if args.DsName != nil {
				return "", -1, errors.New(fmt.Sprintf("DsName is required with %v operation", *args.Opt))
			}
			if args.ForceDynamic {
				options = append(options, "-f")
			}
			if *args.Opt == OPT_ADD {
				options = append(options, "-A")
			} else {
				options = append(options, "-D")
			}
			dsn := *args.DsName
			if args.Sms {
				dsn += ",sms"
			} else if args.Volumen != nil {
				dsn += "," + *args.Volumen
			}
			if args.Persistent != nil {
				options = append(options, persistentOption...)
			}
			options = append(options, dsn)
		case OPT_CHECK_FORMAT | OPT_SET_DYNAMIC | OPT_SET_STATIC:
			options = append(options, "-F")
			if *args.Opt == OPT_SET_DYNAMIC {
				options = append(options, "DYNAMIC")
			} else if *args.Opt == OPT_SET_STATIC {
				options = append(options, "STATIC")
			}
		case OPT_LIST:
			options = append(options, "-lj")
		}
	} else if len(args.Batch) != 0 {
		if args.ForceDynamic {
			options = append(options, "-f")
		}
		for _, b := range args.Batch {
			switch b.Opt {
			case OPT_ADD | OPT_DEL:
				if b.Opt == OPT_ADD {
					options = append(options, "-A")
				} else {
					options = append(options, "-B")
				}
				dsn := b.DsName
				if args.Sms {
					dsn += ",sms"
				} else if args.Volumen != nil {
					dsn += "," + *args.Volumen
				}
				options = append(options, dsn)
			default:
				return "", -1, errors.New(fmt.Sprintf("Invalid operation: %v", ApfOptString(b.Opt)))
			}
		}
		if args.Persistent != nil {
			options = append(options, persistentOption...)
		}
	} else {
		return "", -1, errors.New("i")
	}
	stdout, rc, err := execZaouCmd("apfadm", options)
	return stdout, rc, err
}

func FindLinkList(member string) (string, error) {
	return execSimpleStringCmd("llwhence", []string{member})
}

func FindParmLib(member string) (string, error) {
	return execSimpleStringCmd("parmwhence", []string{member})
}

func FindProcLib(member string) (string, error) {
	return execSimpleStringCmd("procwhence", []string{member})
}

func ListLinkList() ([]string, error) {
	return execSimpleStringListCmd("pll", nil)
}

func ListParmList() ([]string, error) {
	return execSimpleStringListCmd("pparm", nil)
}

func ListProcLib() ([]string, error) {
	return execSimpleStringListCmd("pproc", nil)
}

func ReadConsole(options *rune) (string, error) {
	opt := 'r'
	if options == nil {
		opt = *options
	}
	switch opt {
	case 'h' | 'r' | 'l' | 'd' | 'w' | 'm' | 'y' | 'a':
		return execSimpleStringCmd("pcon", []string{fmt.Sprintf("-%c", opt)})
	default:
		return "", errors.New(fmt.Sprintf("Invalid option -%c", opt))
	}
}

func SearchParamLib(find string) (string, error) {
	return execSimpleStringCmd("parmgrep", []string{find})
}

func SearchProcLib(find string) (string, error) {
	return execSimpleStringCmd("procgrep", []string{find})
}
