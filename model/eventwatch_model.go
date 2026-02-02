package model

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