package helpers

import "fmt"

// Errotype ...
type (
	Errotype struct{}
)

// GetData ...
func (Errotype) GetData(data interface{}, err error) string {
	return fmt.Sprintf("error get data %v : %v", data, err)
}

// SaveData ...
func (Errotype) SaveData(data interface{}, err error) string {
	return fmt.Sprintf("error save data %v : %v", data, err)
}

// CreateData ...
func (Errotype) CreateData(data interface{}, err error) string {
	return fmt.Sprintf("error create data %v : %v", data, err)
}

// PushData ...
func (Errotype) PushData(data interface{}, err error) string {
	return fmt.Sprintf("error push data %v : %v", data, err)
}

// UpdateData ...
func (Errotype) UpdateData(data interface{}, err error) string {
	return fmt.Sprintf("error update data %v : %v", data, err)
}

// DeleteData ...
func (Errotype) DeleteData(data interface{}, err error) string {
	return fmt.Sprintf("error delete data %v : %v", data, err)
}

// Marshal ...
func (Errotype) Marshal(data interface{}, err error) string {
	return fmt.Sprintf("error marshal %v : %v", data, err)
}

// UnMarshal ...
func (Errotype) UnMarshal(data interface{}, err error) string {
	return fmt.Sprintf("error unmarshal %v : %v", data, err)
}

// BindJSON ...
func (Errotype) BindJSON(data interface{}, err error) string {
	return fmt.Sprintf("error bind json %v : %v", data, err)
}

// BindHead ...
func (Errotype) BindHead(data interface{}, err error) string {
	return fmt.Sprintf("error bind head %v : %v", data, err)
}

// BindQuery ...
func (Errotype) BindQuery(data interface{}, err error) string {
	return fmt.Sprintf("error bind query %v : %v", data, err)
}

// SendData ...
func (Errotype) SendData(data interface{}, err error) string {
	return fmt.Sprintf("error send data %v : %v", data, err)
}

// Decode ...
func (Errotype) Decode(data interface{}, err error) string {
	return fmt.Sprintf("error decoding  %v : %v", data, err)
}

// Encode ...
func (Errotype) Encode(data interface{}, err error) string {
	return fmt.Sprintf("error encoding  %v : %v", data, err)
}

// Parse ...
func (Errotype) Parse(data interface{}, err error) string {
	return fmt.Sprintf("error parsing data  %v : %v", data, err)
}

// Write ...
func (Errotype) Write(data interface{}, err error) string {
	return fmt.Sprintf("error write data  %v : %v", data, err)
}

// Read ...
func (Errotype) Read(data interface{}, err error) string {
	return fmt.Sprintf("error read data  %v : %v", data, err)
}
