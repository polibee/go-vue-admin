package response

import "goravel/app/core/shared/pagination"

type Meta struct {
	RequestID  string           `json:"request_id,omitempty"`
	Pagination *pagination.Meta `json:"pagination,omitempty"`
}

type SuccessEnvelope struct {
	Data any  `json:"data"`
	Meta Meta `json:"meta"`
}

func Success(data any, meta Meta) SuccessEnvelope {
	return SuccessEnvelope{Data: data, Meta: meta}
}
