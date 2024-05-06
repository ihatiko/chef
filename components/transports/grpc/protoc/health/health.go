package healthz

import "context"

type Health struct {
}

func (h *Health) Check(ctx context.Context, request *HealthCheckRequest) (*HealthCheckResponse, error) {
	response := new(HealthCheckResponse)
	response.Status = HealthCheckResponse_SERVING
	return response, nil
}
