package metrics

import (
	"fmt"
	"reflect"
	"strings"
)

func MustInit(metrics interface{}, factory Factory) {
	if err := Init(metrics, factory); err != nil {
		panic(err.Error())
	}
}

func Init(m interface{}, factory Factory) error {
	// Allow user to opt out of reporting metrics by passing in nil.
	if factory == nil {
		factory = NullFactory
	}

	counterPtrType := reflect.TypeOf((*Counter)(nil)).Elem()
	gaugePtrType := reflect.TypeOf((*Gauge)(nil)).Elem()
	//timerPtrType := reflect.TypeOf((*Timer)(nil)).Elem()
	histogramPtrType := reflect.TypeOf((*Histogram)(nil)).Elem()

	v := reflect.ValueOf(m).Elem()
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		tags := make(map[string]string)
		var labels []string
		//for k, v := range globalTags {
		//	tags[k] = v
		//}
		//var buckets []float64
		field := t.Field(i)
		metric := field.Tag.Get("metric")
		if metric == "" {
			return fmt.Errorf("field %s is missing a tag 'metric'", field.Name)
		}
		if tagString := field.Tag.Get("tags"); tagString != "" {
			tagPairs := strings.Split(tagString, ",")
			for _, tagPair := range tagPairs {
				tag := strings.Split(tagPair, "=")
				if len(tag) != 2 {
					return fmt.Errorf(
						"Field [%s]: Tag [%s] is not of the form key=value in 'tags' string [%s]",
						field.Name, tagPair, tagString)
				}
				tags[tag[0]] = tag[1]
			}
		}
		if labelString := field.Tag.Get("labels"); labelString != "" {
			labels = append(labels, strings.Split(labelString, ",")...)
		}
		//if bucketString := field.Tag.Get("buckets"); bucketString != "" {
		//	if field.Type.AssignableTo(timerPtrType) {
		//		// TODO: Parse timer duration buckets
		//		return fmt.Errorf(
		//			"Field [%s]: Buckets are not currently initialized for timer metrics",
		//			field.Name)
		//	} else if field.Type.AssignableTo(histogramPtrType) {
		//		bucketValues := strings.Split(bucketString, ",")
		//		for _, bucket := range bucketValues {
		//			b, err := strconv.ParseFloat(bucket, 64)
		//			if err != nil {
		//				return fmt.Errorf(
		//					"Field [%s]: Bucket [%s] could not be converted to float64 in 'buckets' string [%s]",
		//					field.Name, bucket, bucketString)
		//			}
		//			buckets = append(buckets, b)
		//		}
		//	} else {
		//		return fmt.Errorf(
		//			"Field [%s]: Buckets should only be defined for Timer and Histogram metric types",
		//			field.Name)
		//	}
		//}
		help := field.Tag.Get("help")
		var obj interface{}
		if field.Type.AssignableTo(counterPtrType) {
			obj = factory.Counter(Options{
				Name:   metric,
				Tags:   tags,
				Labels: labels,
				Help:   help,
			})
		} else if field.Type.AssignableTo(gaugePtrType) {
			obj = factory.Gauge(Options{
				Name:   metric,
				Tags:   tags,
				Labels: labels,
				Help:   help,
			})
			//} else if field.Type.AssignableTo(timerPtrType) {
			//	// TODO: Add buckets once parsed (see TODO above)
			//	obj = factory.Timer(TimerOptions{
			//		Name: metric,
			//		Tags: tags,
			//		Help: help,
			//	})
		} else if field.Type.AssignableTo(histogramPtrType) {
			obj = factory.Histogram(HistogramOptions{
				Name:   metric,
				Tags:   tags,
				Labels: labels,
				Help:   help,
				//Buckets: buckets,
			})
		} else {
			return fmt.Errorf(
				"Field %s is not a pointer to timer, gauge, or counter",
				field.Name)
		}
		v.Field(i).Set(reflect.ValueOf(obj))
	}
	return nil
}
