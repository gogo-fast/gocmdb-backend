package cloud

type DataPoint struct {
	MetricName string
	InstanceId string
	Timestamps []*float64
	Values     []*float64
}

//type Dimension struct {
//	Name  string
//	Value string
//}
