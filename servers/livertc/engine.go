package livertc

// 此处为状态机

// 如果存在业务逻辑，调用这个...

type EngineMedia struct {
}

var DefaultEngine EngineMedia

func Default() *EngineMedia {
	return &DefaultEngine
}

func (engine *EngineMedia) AddPubSession(id string) (*Session, error) {

	return &Session{
		engine: engine,
	}, nil
}

func (engine *EngineMedia) RemovePubSession(session *Session) error {

	return nil
}

func (engine *EngineMedia) AddSubSession(id string) (*Session, error) {

	return &Session{}, nil
}

func (engine *EngineMedia) RemoveSubSession(session *Session) error {

	return nil
}
