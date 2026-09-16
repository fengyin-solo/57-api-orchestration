package service

func (s *Service) ExportAll() map[string]interface{} {
	return s.store.ExportAll()
}
