package mock

import "github.com/Telenav/osrm-backend/integration/api/osrm/route"

// NewOSRMRouteNoLeg creates a new OSRM Route mock object but no leg on the route.
func NewOSRMRouteNoLeg() *route.Route {

	return &route.Route{
		Distance:   18377,
		Duration:   1099.1,
		Geometry:   "iw`cFtjtgVFrBdCF]iK|LL?us@pJtA|DyDdi@pAbaCsDdj@jBrX_B~FuQvEgXl@_eA|Hik@_Fg}ExB{c@@wy@tJid@`@eMaBkNcMm]{VgYoD}IoBsZbCyHhHsFxWoH|XaUhD`A~JwGuG_RgYgd@tBsBzAzB",
		Weight:     1099.1,
		WeightName: "routability",
	}
}

// NewOSRMRouteOneEmptyLeg creates a new OSRM Route mock object with one empty leg on the route. It's a invalid leg actually.
func NewOSRMRouteOneEmptyLeg() *route.Route {

	return &route.Route{
		Distance:   18377,
		Duration:   1099.1,
		Geometry:   "iw`cFtjtgVFrBdCF]iK|LL?us@pJtA|DyDdi@pAbaCsDdj@jBrX_B~FuQvEgXl@_eA|Hik@_Fg}ExB{c@@wy@tJid@`@eMaBkNcMm]{VgYoD}IoBsZbCyHhHsFxWoH|XaUhD`A~JwGuG_RgYgd@tBsBzAzB",
		Weight:     1099.1,
		WeightName: "routability",
		Legs: []*route.Leg{
			nil,
		},
	}
}

// NewOSRMRouteNoAnnotation creates a new OSRM Route mock object that contains one leg but no annotations on the leg.
func NewOSRMRouteNoAnnotation() *route.Route {

	return &route.Route{
		Distance:   18377,
		Duration:   1099.1,
		Geometry:   "iw`cFtjtgVFrBdCF]iK|LL?us@pJtA|DyDdi@pAbaCsDdj@jBrX_B~FuQvEgXl@_eA|Hik@_Fg}ExB{c@@wy@tJid@`@eMaBkNcMm]{VgYoD}IoBsZbCyHhHsFxWoH|XaUhD`A~JwGuG_RgYgd@tBsBzAzB",
		Weight:     1099.1,
		WeightName: "routability",
		Legs: []*route.Leg{
			&route.Leg{
				Distance: 18377,
				Duration: 1099.1,
				Weight:   1099.1,
				Summary:  "",
			},
		},
	}
}

// NewOSRMRouteNormal creates a new OSRM Route mock object: only contains one leg with valid annotation, no step.
func NewOSRMRouteNormal() *route.Route {

	// GET /route/v1/driving/-121.955678,37.350982;-121.952011,37.350982?&annotations=true
	r := route.Route{
		Distance:   558.8,
		Duration:   60.3,
		Geometry:   "wc~bFzjjgVeFzB?cIMqEhHmD",
		Weight:     60.3,
		WeightName: "routability",
		Legs: []*route.Leg{
			&route.Leg{
				Distance: 558.8,
				Duration: 60.3,
				Weight:   60.3,
				Summary:  "",
				Annotation: &route.Annotation{
					Distance:    []float64{15.936241, 32.829292, 40.00651, 50.602957, 22.899744, 21.308256, 35.896895, 19.00944, 17.594784, 26.5248, 23.01505, 8.841599, 33.671633, 27.767681, 44.431011, 69.128936, 57.603303, 11.706858},
					Duration:    []float64{1.8, 3.7, 4.5, 5.7, 1.6, 1.5, 2.5, 1.3, 1.2, 1.9, 1.6, 0.6, 2.4, 1.9, 5, 7.7, 6.4, 1.3},
					DataSources: []int{00, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
					Nodes:       []int64{8512423140001101, 898879622102, 49943035102, 1167110288102, 877950111102, 12109141850001101, 1152042734102, 12109141860001101, 12109141860002101, 877950113102, 925839364102, 925839365102, 8880790140001101, 925839366102, 877950115102, 237046430001101, 49943051102, 237046420001101, 49943052102},
					Ways:        []int64{851242314, 851242315, 1234704366, 1234704367, 1210914185, 1210914185, 1210914186, 1210914186, 1210914186, 888079012, 888079013, 888079014, 888079014, 888079015, -23704643, -23704643, -23704642, -23704642},
					Weight:      []float64{1.8, 3.7, 4.5, 5.7, 1.6, 1.5, 2.5, 1.3, 1.2, 1.9, 1.6, 0.6, 2.4, 1.9, 5, 7.7, 6.4, 1.3},
					Speed:       []float64{8.9, 8.9, 8.9, 8.9, 14.3, 14.2, 14.4, 14.6, 14.7, 14, 14.4, 14.7, 14, 14.6, 8.9, 9, 9, 9},
					Metadata:    &route.Metadata{DataSourceNames: []string{"lua profile"}},
				},
			},
		},
	}

	return &r
}
