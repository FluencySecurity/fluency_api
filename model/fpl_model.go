package model

import (
	"encoding/json"
)

type RunFPLV2Report struct {
	Entry *ReportTask `json:"entry,omitempty"`

	ReportName string                      `json:"reportName,omitempty"`
	Arguments  []*ProgramArgument `json:"arguments,omitempty"`
}

type GetTaskResponse struct {
	Entry *ReportTask `json:"entry,omitempty"`
	// Logs  []*fplmodel.LogEvent       `json:"logs,omitempty"`
	FPL   string                     `json:"fpl,omitempty"`
}

type GetMetricFPLResultResponse struct {
	Task   *ReportTask  `json:"task,omitempty"`
	FPL    string                      `json:"fpl,omitempty"`
	Result *MetricFPLResponse `json:"result,omitempty"`

	//ReportConfig *fplreportModel.TaskReportConfig `json:"reportConfig,omitempty"`
}

type ReportTask struct {
	ID          uint   `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`

	//Host        string          `json:"host,omitempty"`
	State string `json:"state,omitempty"`
	//Percent     float64         `json:"percent"`
	Error      string          `json:"error,omitempty"`
	Kargs      json.RawMessage `json:"kargs,omitempty"`
	FPLVersion int             `json:"fplVersion"`
	FPL string `json:"fpl,omitempty"`

	Took int64   `json:"took"`
	Cost float64 `json:"cost,omitempty"`

	//ReportName   string               `json:"reportName,omitempty"`
	Arguments    []*ProgramArgument `json:"arguments,omitempty"`
	// ReportConfig *TaskReportConfig           `json:"reportConfig,omitempty"`
	// ExportPipes  []string                    `json:"exportPipes,omitempty"`
}

type ProgramArgument struct {
	Name         string   `json:"name,omitempty"`
	Description  string   `json:"description,omitempty"`
	Value        *string  `json:"value,omitempty"`
	DefaultValue string   `json:"defaultValue,omitempty"`
	Optional     bool     `json:"optional,omitempty"`
	Type         string   `json:"type,omitempty"` // integer, float, string, boolean
	IsList       bool     `json:"isList,omitempty"`
	EnumValues   []string `json:"enumValues,omitempty"`
	Enum         bool     `json:"enum,omitempty"`
}

type MetricFPLResponse struct {
	TaskID  uint                   `json:"taskID,omitempty"`
	Objects []*FplObject           `json:"objects"`
	// Logs    []*LogEvent            `json:"logs"`
	Error   string                 `json:"error"`
	Env     map[string]interface{} `json:"env,omitempty"`
	Cost    float64                `json:"cost,omitempty"`
	// runtime in miliseconds
	Took int64 `json:"took,omitempty"`
	// for the final merged result
	//Accounts []string `json:"accounts,omitempty"`
	//Regions  []string `json:"regions,omitempty"`
}

type FplObject struct {
	Name        string            `json:"name,omitempty"`
	Description string            `json:"description,omitempty"`
	Table       *FplTable         `json:"table,omitempty"`
	Stream      *FplStream        `json:"stream,omitempty"`
	Alert       *FplAlertResult   `json:"alert,omitempty"`
	List        []*PrimitiveValue `json:"list,omitempty"`
	// Value       *FplStackValue    `json:"value,omitempty"`
}

type FplTable struct {
	RowCount uint64                   `json:"rowCount,omitempty"`
	Columns  []*FplColumnInfo         `json:"columns,omitempty"`
	Rows     []map[string]interface{} `json:"rows,omitempty"`
	FplRows  []*FplRowState           `json:"fplRows,omitempty"`

	// if table is loaded from a resource, this field is set
	Resource string `json:"resource,omitempty"`

	keyMap map[string]int `json:"-"`
	colMap map[string]int `json:"-"`
}

type FplColumnInfo struct {
	Name string `json:"name"`
	Unit string `json:"unit,omitempty"`

	ColumnType    string `json:"columnType,omitempty"`
	AggregateType string `json:"aggType,omitempty"`
	IsVariable    bool   `json:"isVariable"`
	IsHidden      bool   `json:"isHidden,omitempty"`
	Dynamic       bool   `json:"dynamic,omitempty"`
}

type FplValueState struct {
	V *PrimitiveValue `json:"v,omitempty"`
}
type FplRowState struct {
	Key    *PrimitiveValue  `json:"key,omitempty"`
	Values []*FplValueState `json:"values,omitempty"`
	index  int              `json:"-"`
}

type PrimitiveValue struct {
	Type      string `json:"type"`
	BoolValue *bool  `json:"bool,omitempty"`
	// NameValue   *string  `json:"name,omitempty"`
	IntValue    *int64                     `json:"int,omitempty"`
	FloatValue  *float64                   `json:"float,omitempty"`
	StringValue *string                    `json:"string,omitempty"`
	List        []*PrimitiveValue          `json:"list,omitempty"`
	Map         map[string]*PrimitiveValue `json:"map,omitempty"`
	JsonObj     map[string]interface{}     `json:"jsonObj,omitempty"`
	JsonArray   []interface{}              `json:"jsonArray,omitempty"`
	Blob        []byte                     `json:"blob,omitempty"`
}

type FplStream struct {
	//Count      uint64   `json:"count,omitempty"`
	From       int64    `json:"from,omitempty"`
	To         int64    `json:"to,omitempty"`
	Interval   int64    `json:"interval,omitempty"`
	Slots      []int64  `json:"slots,omitempty"`
	Dimensions []string `json:"dimensions,omitempty"`

	Metrics []*FplStreamMetric `json:"metrics,omitempty"`
	IsBool  bool               `json:"isBool,omitempty"`

	DimensionName string                   `json:"dimensionName,omitempty"`
	MetricMap     map[int]*FplStreamMetric `json:"-"`

	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`

	Unit      string `json:"unit,omitempty"`
	Metric    string `json:"metric,omitempty"`
	Namespace string `json:"namespace,omitempty"`

	DimensionTableMap map[string]*FplTable `json:"-"`

	//Separator string `json:"separator,omitempty"`
	//dimensionHandle *DimensionHandle `json:"-"`
	//Columns  []*FplColumnInfo         `json:"columns,omitempty"`
	//Rows     []map[string]interface{} `json:"rows,omitempty"`
	//FplRows  []*FplRowState           `json:"fplRows,omitempty"`

}

type FplStreamMetric struct {
	Flags  []int             `json:"flags,omitempty"`
	Values []float64         `json:"values,omitempty"`
	Bools  []bool            `json:"bools,omitempty"`
	Key    string            `json:"key,omitempty"`
	ID     int               `json:"id,omitempty"`
	Tags   []*TableColumnTag `json:"tags,omitempty"`
	Aggs   []float64         `json:"-"` // for sorting purpose
}

type TableColumnTag struct {
	Key    string   `json:"key,omitempty"`
	Values []string `json:"values,omitempty"`
}

type FplAlertResult struct {
	Alerts      []*FplAlert `json:"alerts,omitempty"`
	Slots       []int64     `json:"slots,omitempty"`
	From        int64       `json:"from,omitempty"`
	To          int64       `json:"to,omitempty"`
	Interval    int64       `json:"interval,omitempty"`
	SliderStart int64       `json:"sliderStart,omitempty"`
}

type FplAlert struct {
	// Options     *AnomalyOptions
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	IsAlert     bool   `json:"isAlert,omitempty"`

	Eval *AnomalyEvalT `json:"eval,omitempty"`

	Key string `json:"key,omitempty"`

	Account       string                `json:"account,omitempty"`
	AccountID     string                `json:"accountID,omitempty"`
	Region        string                `json:"region,omitempty"`
	DimensionTags []*StreamDimensionTag `json:"dimensionTags,omitempty"`
	Strategy      string                `json:"strategy,omitempty"`

	Timestamp int64 `json:"timestamp,omitempty"`

	AlertMetric *FplStreamMetric `json:"alertMetric,omitempty"`

	Values     []*StreamHistogram `json:"values,omitempty"`
	Associates []*StreamHistogram `json:"associates,omitempty"`
}

type StreamDimensionTag struct {
	Dimension string            `json:"dimension,omitempty"`
	Key       string            `json:"key,omitempty"`
	Tags      []*TableColumnTag `json:"tags,omitempty"`
}

type StreamHistogram struct {
	Name   string           `json:"name,omitempty"`
	Unit   string           `json:"unit,omitempty"`
	Metric *FplStreamMetric `json:"metric,omitempty"`
}

type AnomalyEvalT struct {
	alertCount int
	lastAlert  int64
	zScore     float64
	value      float64
	key        string
	id         int
	strategy   string
}
