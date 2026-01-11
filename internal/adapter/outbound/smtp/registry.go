package smtp_outbound_adapter

import outbound_port "prabogo/internal/port/outbound"

type adapter struct{}

func NewRegistryAdapter() outbound_port.EmailPort {
	return NewAdapter()
}