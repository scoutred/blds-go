package blds

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"time"
)

type IsoDate struct {
	time.Time
}

func (d *IsoDate) UnmarshalCSV(s string) error {
	if s == "" {
		return nil
	}
	tim, err := time.Parse("2006-01-02", s)
	if err != nil {
		tim, err = time.Parse("1/2/2006", s)
		if err != nil {
			return err
		}
	}

	d.Time = tim

	return nil
}

func (d *IsoDate) UnmarshalJSON(b []byte) error {
	var s string
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}

	if len(s) == 0 {
		return nil
	}

	tim, err := time.Parse("2006-01-02", s)
	if err != nil {
		tim, err = time.Parse("1/2/2006", s)
		if err != nil {
			return err
		}
	}

	d = &IsoDate{tim}

	return nil
}

type Int int

func (i *Int) UnmarshalJSON(b []byte) error {
	var s string
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}

	if len(s) == 0 {
		return nil
	}

	ii, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return err
	}

	*i = Int(ii)
	return nil
}

type Record struct {
	PermitNum        string
	Description      string
	AppliedDate      IsoDate
	IssueDate        IsoDate
	CompletedDate    IsoDate
	StatusCurrent    string
	OriginalAddress1 string
	OriginalAddress2 string
	OriginalCity     string
	OriginalState    string
	OriginalZip      string

	// Reccomended Fields
	Jurisdiction          string
	PermitClass           string
	PermitClassMapped     string
	StatusCurrentMapped   string
	WorkClass             string
	WorkClassMapped       string
	PermitType            string
	PermitTypeMapped      string
	PermitTypeDesc        string
	StatusDate            IsoDate
	TotalSqFt             Int
	Link                  string
	Latitude              float64
	Longitude             float64
	EstProjCost           Int
	HousingUnits          Int
	PIN                   string
	ContractorCompanyName string
	ContractorTrade       string
	ContractorTradeMapped string
	ContractorLicNum      string
	ContractorStateLic    string

	// Optional
	ProposedUse           string
	AddedSqFt             Int
	RemovedSqFt           Int
	MasterPermitNum       string
	ExpiresDate           IsoDate
	COIssuedDate          IsoDate
	HoldDate              IsoDate
	VoidDate              IsoDate
	ProjectName           string
	ProjectID             string
	TotalFinieshedSqFt    string
	TotalUnfinishedSqFt   string
	TotalHeadtedSqFt      string
	TotalUnHeatedSqFt     string
	TotalAccSqFt          string
	TotalSprinkledSqFt    string
	ExtraFields           interface{}
	Publisher             string
	Fee                   float64 `json:",string"`
	ContractorFullName    string
	ContractorCompanyDesc string
	ContractorPhone       string
	ContractorAddress1    string
	ContractorAddress2    string
	ContractorCity        string
	ContractorState       string
	ContractorZip         string
	ContractorEmail       string
}

type Unmarshaler interface {
	UnmarshalCSV(string) error
}

func FromCSV(r io.Reader) ([]Record, error) {

	rcsv := csv.NewReader(r)
	header, err := rcsv.Read()
	if err != nil {
		return nil, err
	}

	recordT := reflect.TypeOf(Record{})
	csv2struct := make([]int, len(header))
	for i, field := range header {
		sfield, ok := recordT.FieldByName(field)
		if !ok {
			csv2struct[i] = -1
		} else {
			csv2struct[i] = sfield.Index[0]
		}
	}

	ret := make([]Record, 0, 100)

	for {
		record, err := rcsv.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		recordRet := reflect.New(recordT).Elem()
		for i, v := range record {

			fieldIdx := csv2struct[i]

			if fieldIdx < 0 || v == "" {
				// field not present in csv or not in record
				continue
			}

			field := recordRet.Field(fieldIdx)

			if vv, ok := field.Interface().(Unmarshaler); ok {
				err = vv.UnmarshalCSV(v)
				if err != nil {
					return nil, fmt.Errorf("could not unmarshal field %v: %v",
						recordT.Field(fieldIdx).Name, err)
				}

				continue
			} else if vv, ok := field.Addr().Interface().(Unmarshaler); ok {
				err = vv.UnmarshalCSV(v)
				if err != nil {
					return nil, fmt.Errorf("could not unmarshal field %v: %v",
						recordT.Field(fieldIdx).Name, err)
				}

				continue
			}

			typeErr := fmt.Errorf("cannot unmarshal %v into field %v of type %T",
				v, recordT.Field(fieldIdx).Name, field.Interface())

			switch field.Kind() {
			case reflect.Bool:
				b, err := strconv.ParseBool(v)
				if err != nil {
					return nil, typeErr
				}

				field.SetBool(b)

			case reflect.Int,
				reflect.Int8,
				reflect.Int16,
				reflect.Int32,
				reflect.Int64:

				i, err := strconv.ParseInt(v, 10, 64)
				if err != nil {
					return nil, typeErr
				}

				field.SetInt(i)

			case reflect.Uint,
				reflect.Uint8,
				reflect.Uint16,
				reflect.Uint32,
				reflect.Uint64,
				reflect.Uintptr:

				i, err := strconv.ParseUint(v, 10, 64)
				if err != nil {
					return nil, typeErr
				}

				field.SetUint(i)

			case reflect.Float32,
				reflect.Float64:

				f, err := strconv.ParseFloat(v, 64)
				if err != nil {
					return nil, typeErr
				}

				field.SetFloat(f)

			case reflect.String:
				field.SetString(v)

			default:
				return nil, typeErr
			}
		}

		ret = append(ret, recordRet.Interface().(Record))
	}

	return ret, nil
}

type jsonResponse struct {
	Result struct {
		Records []Record `json:"records"`
	} `json:"result"`
}

func FromJSON(r io.Reader) ([]Record, error) {
	var ret jsonResponse
	err := json.NewDecoder(r).Decode(&ret)
	if err != nil {
		return nil, err
	}

	return ret.Result.Records, nil
}
