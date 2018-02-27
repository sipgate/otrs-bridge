package types

type TicketStateUpdateEvent struct {
	OldTicketData struct {
		State                  string `json:"State"`
		Type                   string `json:"Type"`
		Changed                string `json:"Changed"`
		Responsible            string `json:"Responsible"`
		EscalationTime         string `json:"EscalationTime"`
		PriorityID             string `json:"PriorityID"`
		CustomerID             string `json:"CustomerID"`
		ServiceID              string `json:"ServiceID"`
		EscalationResponseTime string `json:"EscalationResponseTime"`
		Age                    int    `json:"Age"`
		CustomerUserID         string `json:"CustomerUserID"`
		UntilTime              int    `json:"UntilTime"`
		EscalationUpdateTime   string `json:"EscalationUpdateTime"`
		Lock                   string `json:"Lock"`
		ChangeBy               string `json:"ChangeBy"`
		TicketNumber           string `json:"TicketNumber"`
		StateID                string `json:"StateID"`
		Owner                  string `json:"Owner"`
		UnlockTimeout          string `json:"UnlockTimeout"`
		Title                  string `json:"Title"`
		OwnerID                string `json:"OwnerID"`
		SLAID                  string `json:"SLAID"`
		ArchiveFlag            string `json:"ArchiveFlag"`
		Priority               string `json:"Priority"`
		LockID                 string `json:"LockID"`
		TicketID               string `json:"TicketID"`
		TypeID                 string `json:"TypeID"`
		RealTillTimeNotUsed    string `json:"RealTillTimeNotUsed"`
		StateType              string `json:"StateType"`
		EscalationSolutionTime string `json:"EscalationSolutionTime"`
		CreateTimeUnix         string `json:"CreateTimeUnix"`
		ResponsibleID          string `json:"ResponsibleID"`
		QueueID                string `json:"QueueID"`
		Created                string `json:"Created"`
		GroupID                string `json:"GroupID"`
		Queue                  string `json:"Queue"`
		CreateBy               string `json:"CreateBy"`
	} `json:"OldTicketData"`
}