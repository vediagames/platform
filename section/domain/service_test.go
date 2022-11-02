package domain

import "testing"

func TestEditWebsitePlacementsRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		req     EditWebsitePlacementsRequest
		wantErr bool
	}{
		{
			name: "valid",
			req: EditWebsitePlacementsRequest{
				WebsitePlacements: map[Placement]int{
					1: 1,
					2: 2,
					3: 3,
				},
			},
		},
		{
			name: "invalid placement",
			req: EditWebsitePlacementsRequest{
				WebsitePlacements: map[Placement]int{
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
				t.Errorf("EditWebsitePlacementsRequest.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}

			t.Log(err)
		})
	}
}
