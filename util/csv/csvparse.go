package csv

import (
	"encoding/csv"
	"errors"
	"io"
	"os"
	"strconv"

	"github.com/axgle/mahonia"
)

func ConvertString(src string) string {
	return mahonia.NewDecoder("gbk").ConvertString(src)
}

type CsvTable struct {
	FileName string
	Records  []CsvRecord
}

type CsvRecord struct {
	Record map[string]string
}

func (c *CsvRecord) GetInt(field string) (int, error) {
	var r int
	var err error
	if r, err = strconv.Atoi(c.Record[field]); err != nil {
		return 0, err
	}
	return r, nil
}

func (c *CsvRecord) GetString(field string) string {
	data, ok := c.Record[field]
	if ok {
		return data
	} else {
		return ""
	}
}

func LoadCsvCfg(filename string, row int) (*CsvTable, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	if reader == nil {
		return nil, err
	}
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}
	if len(records) < row {
		return nil, errors.New("data less than row")
	}
	colNum := len(records[0])
	recordNum := len(records)
	var allRecords []CsvRecord
	for i := row; i < recordNum; i++ {
		record := &CsvRecord{make(map[string]string)}
		for k := 0; k < colNum; k++ {
			record.Record[records[0][k]] = records[i][k]
		}
		allRecords = append(allRecords, *record)
	}
	var result = &CsvTable{
		filename,
		allRecords,
	}
	return result, nil
}

func ReadLineCSV(filename string, f func(int, []string)) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	if reader == nil {
		return err
	}
	reader.FieldsPerRecord = -1
	for i := 0; ; i++ {
		records, err := reader.Read()
		if err != nil && err != io.EOF {
			return err
		}
		if err == io.EOF || len(records) == 0 {
			return nil
		}
		f(i, records)
	}
	return nil
}
