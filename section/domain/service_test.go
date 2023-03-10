package domain

import (
	"testing"
)

func TestEditWebsitePlacementsRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		req     EditPlacedRequest
		wantErr bool
	}{
		{
			name: "valid",
			req: EditPlacedRequest{
				Placements: map[Placement]int{
					1: 1,
					2: 2,
					3: 3,
				},
			},
		},
		{
			name: "invalid placement",
			req: EditPlacedRequest{
				Placements: map[Placement]int{
					1: 1,
					3: 2,
					4: 3,
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.req.Validate()
			if err != nil != tt.wantErr {
				t.Errorf("EditPlacedRequest.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}

			t.Log(err)
		})
	}
}
