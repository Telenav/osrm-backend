package osrmconnector

import "github.com/Telenav/osrm-backend/integration/pkg/api/osrm/route"

type RouteRequster interface {
	Request4Route(r *route.Request) <-chan RouteResponse
}
