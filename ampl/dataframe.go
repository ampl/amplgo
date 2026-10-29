package ampl

/*
#cgo LDFLAGS: -L./lib/amd64 -lampl
#cgo CFLAGS: -I./include
#include "ampl/ampl_c.h"
#include <stdlib.h>
*/
import "C"
import (
  "fmt"
  "runtime"
  "unsafe"
)

type DataFrame struct {
  cdataframe *C.AMPL_DATAFRAME
}

func makeCStringArray(values []string) (**C.char) {
  cArray := C.malloc(C.size_t(len(values)) * C.size_t(unsafe.Sizeof(uintptr(0))))
  arr := (*[1 << 30]*C.char)(cArray)[:len(values):len(values)]
  for i, v := range values {
    arr[i] = C.CString(v)
  }
  return (**C.char)(cArray)
}

func freeCStringArray(cArray **C.char, n int) {
  arr := (*[1 << 30]*C.char)(unsafe.Pointer(cArray))[:n:n]
  for i := 0; i < n; i++ {
    C.free(unsafe.Pointer(arr[i]))
  }
  C.free(unsafe.Pointer(cArray))
}

func wrapDataFrame(cdataframe *C.AMPL_DATAFRAME) (*DataFrame) {
  wrapper := &DataFrame{cdataframe: cdataframe}

  runtime.SetFinalizer(wrapper, func(dataframe *DataFrame) {
    if dataframe.cdataframe != nil {
      C.AMPL_DataFrameFree(&dataframe.cdataframe)
      dataframe.cdataframe = nil
    }
  })

  return wrapper
}

// NewDataFrame allocates a DataFrame skeleton with the given number of
// index and data columns and the corresponding column headers.
func NewDataFrame(numberOfIndexColumns int, numberOfDataColumns int, headers []string) (*DataFrame, error) {
  var cdataframe *C.AMPL_DATAFRAME

  cHeaders := makeCStringArray(headers)
  defer freeCStringArray(cHeaders, len(headers))

  errInfo := C.AMPL_DataFrameCreate(&cdataframe, C.size_t(numberOfIndexColumns), C.size_t(numberOfDataColumns), cHeaders)
  if errInfo != nil {
    errorMsg := C.GoString(C.AMPL_ErrorInfoGetMessage(errInfo))
    return nil, fmt.Errorf("%s", errorMsg)
  }

  return wrapDataFrame(cdataframe), nil
}

// NewDataFrameIndexed allocates a DataFrame skeleton with the given number
// of index columns and no data columns.
func NewDataFrameIndexed(numberOfIndexColumns int) (*DataFrame, error) {
  var cdataframe *C.AMPL_DATAFRAME

  errInfo := C.AMPL_DataFrameCreate2(&cdataframe, C.size_t(numberOfIndexColumns))
  if errInfo != nil {
    errorMsg := C.GoString(C.AMPL_ErrorInfoGetMessage(errInfo))
    return nil, fmt.Errorf("%s", errorMsg)
  }

  return wrapDataFrame(cdataframe), nil
}

func (dataframe *DataFrame) Copy() (*DataFrame, error) {
  var ccopy *C.AMPL_DATAFRAME

  errInfo := C.AMPL_DataFrameCopy(&ccopy, dataframe.cdataframe)
  if errInfo != nil {
    errorMsg := C.GoString(C.AMPL_ErrorInfoGetMessage(errInfo))
    return nil, fmt.Errorf("%s", errorMsg)
  }

  return wrapDataFrame(ccopy), nil
}

func (dataframe *DataFrame) GetHeaders() ([]string, error) {
  var size C.size_t
  var cHeaders **C.char

  errInfo := C.AMPL_DataFrameGetHeaders(dataframe.cdataframe, &size, &cHeaders)
  if errInfo != nil {
    errorMsg := C.GoString(C.AMPL_ErrorInfoGetMessage(errInfo))
    return nil, fmt.Errorf("%s", errorMsg)
  }

  headers := make([]string, size)
  cArray := (*[1 << 30]*C.char)(unsafe.Pointer(cHeaders))[:size:size]
  for i := 0; i < int(size); i++ {
    headers[i] = C.GoString(cArray[i])
    C.AMPL_StringFree(&cArray[i])
  }
  C.free(unsafe.Pointer(cHeaders))

  return headers, nil
}

func (dataframe *DataFrame) Equals(other *DataFrame) (bool, error) {
  var equals C.int

  errInfo := C.AMPL_DataFrameEquals(dataframe.cdataframe, other.cdataframe, &equals)
  if errInfo != nil {
    errorMsg := C.GoString(C.AMPL_ErrorInfoGetMessage(errInfo))
    return false, fmt.Errorf("%s", errorMsg)
  }

  return equals != 0, nil
}

func (dataframe *DataFrame) ToString() (string, error) {
  var cOutput *C.char

  errInfo := C.AMPL_DataFrameToString(dataframe.cdataframe, &cOutput)
  if errInfo != nil {
    errorMsg := C.GoString(C.AMPL_ErrorInfoGetMessage(errInfo))
    return "", fmt.Errorf("%s", errorMsg)
  }

  output := C.GoString(cOutput)
  C.AMPL_StringFree(&cOutput)
  return output, nil
}

// Reserve allocates space for the given number of rows. The rows still
// need to be added via AddRow.
func (dataframe *DataFrame) Reserve(numRows int) error {
  errInfo := C.AMPL_DataFrameReserve(dataframe.cdataframe, C.size_t(numRows))
  if errInfo != nil {
    errorMsg := C.GoString(C.AMPL_ErrorInfoGetMessage(errInfo))
    return fmt.Errorf("%s", errorMsg)
  }
  return nil
}

// AddRow adds a row to the DataFrame. The size of the tuple must equal the
// total number of columns in the DataFrame.
func (dataframe *DataFrame) AddRow(value *Tuple) error {
  errInfo := C.AMPL_DataFrameAddRow(dataframe.cdataframe, value.ctuple)
  if errInfo != nil {
    errorMsg := C.GoString(C.AMPL_ErrorInfoGetMessage(errInfo))
    return fmt.Errorf("%s", errorMsg)
  }
  return nil
}

// SetColumnDouble sets the values of a numeric column.
func (dataframe *DataFrame) SetColumnDouble(header string, column []float64) error {
  var cColumn *C.double
  if len(column) > 0 {
    cColumn = (*C.double)(unsafe.Pointer(&column[0]))
  }

  errInfo := C.AMPL_DataFrameSetColumnArgDouble(dataframe.cdataframe, C.CString(header), cColumn, C.size_t(len(column)))
  if errInfo != nil {
    errorMsg := C.GoString(C.AMPL_ErrorInfoGetMessage(errInfo))
    return fmt.Errorf("%s", errorMsg)
  }
  return nil
}

// SetColumnString sets the values of a string column.
func (dataframe *DataFrame) SetColumnString(header string, column []string) error {
  cColumn := makeCStringArray(column)
  defer freeCStringArray(cColumn, len(column))

  errInfo := C.AMPL_DataFrameSetColumnArgString(dataframe.cdataframe, C.CString(header), cColumn, C.size_t(len(column)))
  if errInfo != nil {
    errorMsg := C.GoString(C.AMPL_ErrorInfoGetMessage(errInfo))
    return fmt.Errorf("%s", errorMsg)
  }
  return nil
}

// SetMatrixDoubles assigns a flattened 2-d matrix of doubles to an empty
// DataFrame with two indexing columns. values must be indices0 x indices1
// values in row-major order (indices0 varying slowest).
func (dataframe *DataFrame) SetMatrixDoubles(indices0 []string, indices1 []string, values []float64) error {
  var cValues *C.double
  if len(values) > 0 {
    cValues = (*C.double)(unsafe.Pointer(&values[0]))
  }

  cIndices0 := makeCStringArray(indices0)
  defer freeCStringArray(cIndices0, len(indices0))
  cIndices1 := makeCStringArray(indices1)
  defer freeCStringArray(cIndices1, len(indices1))

  errInfo := C.AMPL_DataFrameSetMatrixStringString(dataframe.cdataframe, cValues, C.size_t(len(indices0)), cIndices0, C.size_t(len(indices1)), cIndices1)
  return amplError(errInfo)
}

// SetValue sets the value at the row identified by rowIndex and the column
// identified by header.
func (dataframe *DataFrame) SetValue(rowIndex *Tuple, header string, value *Variant) error {
  errInfo := C.AMPL_DataFrameSetValue(dataframe.cdataframe, rowIndex.ctuple, C.CString(header), value.cvariant)
  if errInfo != nil {
    errorMsg := C.GoString(C.AMPL_ErrorInfoGetMessage(errInfo))
    return fmt.Errorf("%s", errorMsg)
  }
  return nil
}

// SetValueByIndex sets the value at the given 0-based row and column
// indices (columns include the indexing columns).
func (dataframe *DataFrame) SetValueByIndex(rowNumber int, colNumber int, value *Variant) error {
  errInfo := C.AMPL_DataFrameSetValueByIndex(dataframe.cdataframe, C.size_t(rowNumber), C.size_t(colNumber), value.cvariant)
  if errInfo != nil {
    errorMsg := C.GoString(C.AMPL_ErrorInfoGetMessage(errInfo))
    return fmt.Errorf("%s", errorMsg)
  }
  return nil
}

// AddColumnDouble adds a new numeric column with the given header.
func (dataframe *DataFrame) AddColumnDouble(header string, values []float64) error {
  var cValues *C.double
  if len(values) > 0 {
    cValues = (*C.double)(unsafe.Pointer(&values[0]))
  }

  errInfo := C.AMPL_DataFrameAddColumnDouble(dataframe.cdataframe, C.CString(header), cValues)
  if errInfo != nil {
    errorMsg := C.GoString(C.AMPL_ErrorInfoGetMessage(errInfo))
    return fmt.Errorf("%s", errorMsg)
  }
  return nil
}

// AddColumnString adds a new string column with the given header.
func (dataframe *DataFrame) AddColumnString(header string, values []string) error {
  cValues := makeCStringArray(values)
  defer freeCStringArray(cValues, len(values))

  errInfo := C.AMPL_DataFrameAddColumnString(dataframe.cdataframe, C.CString(header), cValues)
  if errInfo != nil {
    errorMsg := C.GoString(C.AMPL_ErrorInfoGetMessage(errInfo))
    return fmt.Errorf("%s", errorMsg)
  }
  return nil
}

// AddEmptyColumn adds a new empty column with the given header.
func (dataframe *DataFrame) AddEmptyColumn(header string) error {
  errInfo := C.AMPL_DataFrameAddEmptyColumn(dataframe.cdataframe, C.CString(header))
  if errInfo != nil {
    errorMsg := C.GoString(C.AMPL_ErrorInfoGetMessage(errInfo))
    return fmt.Errorf("%s", errorMsg)
  }
  return nil
}

// GetNumCols returns the total number of columns (indexarity + values).
func (dataframe *DataFrame) GetNumCols() (int, error) {
  var num C.size_t
  errInfo := C.AMPL_DataFrameGetNumCols(dataframe.cdataframe, &num)
  if errInfo != nil {
    errorMsg := C.GoString(C.AMPL_ErrorInfoGetMessage(errInfo))
    return 0, fmt.Errorf("%s", errorMsg)
  }
  return int(num), nil
}

// GetNumRows returns the number of data rows.
func (dataframe *DataFrame) GetNumRows() (int, error) {
  var num C.size_t
  errInfo := C.AMPL_DataFrameGetNumRows(dataframe.cdataframe, &num)
  if errInfo != nil {
    errorMsg := C.GoString(C.AMPL_ErrorInfoGetMessage(errInfo))
    return 0, fmt.Errorf("%s", errorMsg)
  }
  return int(num), nil
}

// GetNumIndices returns the number of indices (the indexarity).
func (dataframe *DataFrame) GetNumIndices() (int, error) {
  var num C.size_t
  errInfo := C.AMPL_DataFrameGetNumIndices(dataframe.cdataframe, &num)
  if errInfo != nil {
    errorMsg := C.GoString(C.AMPL_ErrorInfoGetMessage(errInfo))
    return 0, fmt.Errorf("%s", errorMsg)
  }
  return int(num), nil
}

// GetColumnIndex returns the column index of the column with the given
// header.
func (dataframe *DataFrame) GetColumnIndex(name string) (int, error) {
  var columnIndex C.size_t
  errInfo := C.AMPL_DataFrameGetColumnIndex(dataframe.cdataframe, C.CString(name), &columnIndex)
  if errInfo != nil {
    errorMsg := C.GoString(C.AMPL_ErrorInfoGetMessage(errInfo))
    return 0, fmt.Errorf("%s", errorMsg)
  }
  return int(columnIndex), nil
}

// GetIndexingTuple returns the tuple representation of the index of the
// given row.
func (dataframe *DataFrame) GetIndexingTuple(rowIndex int) (*Tuple, error) {
  var cTuple *C.AMPL_TUPLE
  errInfo := C.AMPL_DataFrameGetIndexingTuple(dataframe.cdataframe, C.size_t(rowIndex), &cTuple)
  if errInfo != nil {
    errorMsg := C.GoString(C.AMPL_ErrorInfoGetMessage(errInfo))
    return nil, fmt.Errorf("%s", errorMsg)
  }
  return &Tuple{ctuple: cTuple}, nil
}

// GetRowIndex returns the row index identified by the given tuple.
func (dataframe *DataFrame) GetRowIndex(index *Tuple) (int, error) {
  var rowIndex C.size_t
  errInfo := C.AMPL_DataFrameGetRowIndex(dataframe.cdataframe, index.ctuple, &rowIndex)
  if errInfo != nil {
    errorMsg := C.GoString(C.AMPL_ErrorInfoGetMessage(errInfo))
    return 0, fmt.Errorf("%s", errorMsg)
  }
  return int(rowIndex), nil
}

// Element returns the value at the given row and column indices.
func (dataframe *DataFrame) Element(rowIndex int, colIndex int) (*Variant, error) {
  var cvariant *C.AMPL_VARIANT
  errInfo := C.AMPL_DataFrameElement(dataframe.cdataframe, C.size_t(rowIndex), C.size_t(colIndex), &cvariant)
  if errInfo != nil {
    errorMsg := C.GoString(C.AMPL_ErrorInfoGetMessage(errInfo))
    return nil, fmt.Errorf("%s", errorMsg)
  }
  return &Variant{cvariant: cvariant}, nil
}
