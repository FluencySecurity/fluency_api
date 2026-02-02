package model

type BehaviorSummary struct {
	ID       string `json:"id" bson:"id"`
	From     int64  `json:"from" bson:"from"`
	To       int64  `json:"to" bson:"to"`
	Count    int    `json:"count" bson:"count"`
	Key      string `json:"key" bson:"key"`
	KeyType  string `json:"keyType,omitempty" bson:"keyType,omitempty"`
	DayIndex string `json:"dayIndex,omitempty" bson:"dayIndex,omitempty"`
	Interval string `json:"interval,omitempty" bson:"interval,omitempty"`

	BehaviorRules []string `json:"behaviorRules" bson:"behaviorRules"`
	Behaviors     []string `json:"behaviors" bson:"behaviors"`
	Risks         []string `json:"risks,omitempty" bson:"risks,omitempty"`
	RiskScore     int      `json:"riskScore" bson:"riskScore"`
	MitreTags     []string `json:"mitreTags,omitempty" bson:"mitreTags,omitempty"`

	SummaryList []*BehaviorRuleSummary `json:"summaryList,omitempty" bson:"summaryList,omitempty"`
	// Slots       []*TimeSlotSummary     `json:"slots,omitempty" bson:"slots,omitempty"`

	ScoreLevel string `json:"scoreLevel" bson:"scoreLevel"`

	// KeyContext *recommend.EntityContext `json:"keyContext" bson:"keyContext"`

	// incident management
	Comments    []*UserComment `json:"comments" bson:"comments"`
	Status      string         `json:"status,omitempty" bson:"status,omitempty"`
	Incident    bool           `json:"incident" bson:"incident"`
	ScoreAdjust int            `json:"scoreAdjust"`
	UpdatedOn   int64          `json:"updatedOn" bson:"updatedOn"`

	NotifyFlag bool `json:"-"`

	// Investigations []*Investigation `json:"investigations,omitempty" bson:"investigations,omitempty"`

	// openAI embedding APIs
	VectorData      []float64 `json:"vectorData,omitempty" bson:"vectorData,omitempty"`
	FingerprintHash string    `json:"fingerprintHash,omitempty" bson:"fingerprintHash,omitempty"`
	Fingerprint     string    `json:"fingerprint,omitempty" bson:"fingerprint,omitempty"`
	RuleFP          string    `json:"ruleFP,omitempty" bson:"ruleFP,omitempty"`
}

type UserComment struct {
	Content  string   `json:"content"`
	Actions  []string `json:"actions"`
	Username string   `json:"username"`
	// NewState  string    `json:"newState"`
	CreatedOn int64 `json:"createdOn"`
}

type BehaviorRuleSummary struct {
	BehaviorRule       string                `json:"behaviorRule" bson:"behaviorRule"`
	Behavior           string                `json:"behavior" bson:"behavior"`
	From               int64                 `json:"from" bson:"from"`
	To                 int64                 `json:"to" bson:"to"`
	Count              int                   `json:"count" bson:"count"`
	AttributeSummaries []*ValueSummmaryEntry `json:"attributeSummaries" bson:"attributeSummaries"`
	Hits               []*BehaviorRuleHit    `json:"hits,omitempty" bson:"hits,omitempty"`
	RiskScore          int                   `json:"riskScore" bson:"riskScore"`
	//scoreSet           *risk.ScoreSet                              `json:"-" bson:"-"`
	//hitMap             map[string][]*eventWatchDao.BehaviorRuleHit `json:"-" bson:"-"`
	fieldMap map[string]*ValueSummmaryEntry `json:"-" bson:"-"`
	Risks    []string                       `json:"risks" bson:"risks"`
}

type BehaviorRuleHit struct {
	Name        string   `json:"name" bson:"name"`
	Description string   `json:"description,omitempty" bson:"description,omitempty"`
	Fields      []string `json:"fields,omitempty" bson:"fields,omitempty"`
	Values      []string `json:"values,omitempty" bson:"values,omitempty"`
	Scope       string   `json:"scope,omitempty" bson:"scope,omitempty"`
	Risks       []string `json:"risks,omitempty" bson:"risks,omitempty"`
}

type ValueSummmaryEntry struct {
	Key    string   `json:"key" bson:"key"`
	Aliase string   `json:"aliase,omitempty" bson:"aliase,omitempty"`
	Values []string `json:"values" bson:"values"`
	//	Value interface{} `json:"value"`
}

type ValueEntry struct {
	Key    string `json:"key"`
	Aliase string `json:"aliase,omitempty"`
	Value  string `json:"value"`
	//	Value interface{} `json:"value"`
}

type BehaviorEvent struct {
	Timestamp    int64              `json:"timestamp"`
	Key          string             `json:"key"`
	KeyType      string             `json:"keyType,omitempty"`
	Title        string             `json:"title,omitempty"`
	BehaviorRule string             `json:"behaviorRule"`
	Behavior     string             `json:"behavior"`
	RiskScore    int                `json:"riskScore"`
	Attributes   []*ValueEntry      `json:"attributes"`
	Risks        []string           `json:"risks,omitempty"`
	RuleRisks    []string           `json:"ruleRisks,omitempty"`
	RuleHits     []*BehaviorRuleHit `json:"ruleHits,omitempty"`
}