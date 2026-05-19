package dto

type CreateReportItemRequest struct {
	ItemID   uint `json:"item_id"`
	Quantity int  `json:"quantity"`
}

type CreateReportRequest struct {
	VehicleID    uint                      `json:"vehicle_id"`
	Odometer     int                       `json:"odometer"`
	Complaint    string                    `json:"complaint"`
	InitialPhoto string                    `json:"initial_photo"`
	Items        []CreateReportItemRequest `json:"items"`
}

type CompleteReportRequest struct {
	ProofPhoto string `json:"proof_photo"`
}
