package service

func (s *Service) StartListenerManager() error {
	if err := s.mylm.StartListening(); err != nil {
		return err
	}
	return nil
}

func (s *Service) StopListenerManager() error {
	if err := s.mylm.StopListening(); err != nil {
		return nil
	}
	return nil
}
