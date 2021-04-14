package rtmp

type FlvStorage struct {
}

func (this *FlvStorage) Write(*Packet) error {
	return nil
}
