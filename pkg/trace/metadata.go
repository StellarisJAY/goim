package trace

import "google.golang.org/grpc/metadata"

type MetaDataReadWriter struct {
	metadata.MD
}

func (m *MetaDataReadWriter) ForeachKey(handler func(key string, val string) error) error {
	for k, vs := range m.MD {
		for _, v := range vs {
			if err := handler(k, v); err != nil {
				return err
			}
		}
	}
	return nil
}

func (m *MetaDataReadWriter) Set(key, val string) {
	m.MD.Set(key, val)
}
