package testing

import (
	"bytes"
	"context"
	"fmt"
	"sort"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/influxdata/platform"
	"github.com/influxdata/platform/mock"
)

const (
	dashOneID   = "020f755c3c082000"
	dashTwoID   = "020f755c3c082001"
	dashThreeID = "020f755c3c082002"
)

var dashboardCmpOptions = cmp.Options{
	cmp.Comparer(func(x, y []byte) bool {
		return bytes.Equal(x, y)
	}),
	cmp.Transformer("Sort", func(in []*platform.Dashboard) []*platform.Dashboard {
		out := append([]*platform.Dashboard(nil), in...) // Copy input to avoid mutating it
		sort.Slice(out, func(i, j int) bool {
			return out[i].ID.String() > out[j].ID.String()
		})
		return out
	}),
}

// DashboardFields will include the IDGenerator, and dashboards
type DashboardFields struct {
	IDGenerator   platform.IDGenerator
	Dashboards    []*platform.Dashboard
	Organizations []*platform.Organization
}

// CreateDashboard testing
func CreateDashboard(
	init func(DashboardFields, *testing.T) (platform.DashboardService, func()),
	t *testing.T,
) {
	type args struct {
		dashboard *platform.Dashboard
	}
	type wants struct {
		err        error
		dashboards []*platform.Dashboard
	}

	tests := []struct {
		name   string
		fields DashboardFields
		args   args
		wants  wants
	}{
		{
			name: "create dashboards with empty set",
			fields: DashboardFields{
				IDGenerator: mock.NewIDGenerator(dashOneID, t),
				Dashboards:  []*platform.Dashboard{},
				Organizations: []*platform.Organization{
					{
						Name: "theorg",
						ID:   MustIDFromString(orgOneID),
					},
				},
			},
			args: args{
				dashboard: &platform.Dashboard{
					Name:           "name1",
					OrganizationID: MustIDFromString(orgOneID),
					Cells: []platform.DashboardCell{
						{
							DashboardCellContents: platform.DashboardCellContents{
								Name: "hello",
								X:    10,
								Y:    10,
								W:    100,
								H:    12,
							},
							Visualization: platform.CommonVisualization{
								Query: "SELECT * FROM foo",
							},
						},
					},
				},
			},
			wants: wants{
				dashboards: []*platform.Dashboard{
					{
						Name:           "name1",
						ID:             MustIDFromString(dashOneID),
						OrganizationID: MustIDFromString(orgOneID),
						Organization:   "theorg",
						Cells: []platform.DashboardCell{
							{
								DashboardCellContents: platform.DashboardCellContents{
									ID:   MustIDFromString(dashOneID),
									Name: "hello",
									X:    10,
									Y:    10,
									W:    100,
									H:    12,
								},
								Visualization: platform.CommonVisualization{
									Query: "SELECT * FROM foo",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "basic create dashboard",
			fields: DashboardFields{
				IDGenerator: &mock.IDGenerator{
					IDFn: func() platform.ID {
						return MustIDFromString(dashTwoID)
					},
				},
				Dashboards: []*platform.Dashboard{
					{
						ID:             MustIDFromString(dashOneID),
						Name:           "dashboard1",
						OrganizationID: MustIDFromString(orgOneID),
					},
				},
				Organizations: []*platform.Organization{
					{
						Name: "theorg",
						ID:   MustIDFromString(orgOneID),
					},
					{
						Name: "otherorg",
						ID:   MustIDFromString(orgTwoID),
					},
				},
			},
			args: args{
				dashboard: &platform.Dashboard{
					Name:           "dashboard2",
					OrganizationID: MustIDFromString(orgTwoID),
				},
			},
			wants: wants{
				dashboards: []*platform.Dashboard{
					{
						ID:             MustIDFromString(dashOneID),
						Name:           "dashboard1",
						Organization:   "theorg",
						OrganizationID: MustIDFromString(orgOneID),
					},
					{
						ID:             MustIDFromString(dashTwoID),
						Name:           "dashboard2",
						Organization:   "otherorg",
						OrganizationID: MustIDFromString(orgTwoID),
					},
				},
			},
		},
		{
			name: "basic create dashboard using org name",
			fields: DashboardFields{
				IDGenerator: mock.NewIDGenerator(orgTwoID, t),
				Dashboards: []*platform.Dashboard{
					{
						ID:             MustIDFromString(dashOneID),
						Name:           "dashboard1",
						OrganizationID: MustIDFromString(orgOneID),
					},
				},
				Organizations: []*platform.Organization{
					{
						Name: "theorg",
						ID:   MustIDFromString(orgOneID),
					},
					{
						Name: "otherorg",
						ID:   MustIDFromString(orgTwoID),
					},
				},
			},
			args: args{
				dashboard: &platform.Dashboard{
					Name:         "dashboard2",
					Organization: "otherorg",
				},
			},
			wants: wants{
				dashboards: []*platform.Dashboard{
					{
						ID:             MustIDFromString(dashOneID),
						Name:           "dashboard1",
						Organization:   "theorg",
						OrganizationID: MustIDFromString(orgOneID),
					},
					{
						ID:             MustIDFromString(dashTwoID),
						Name:           "dashboard2",
						Organization:   "otherorg",
						OrganizationID: MustIDFromString(orgTwoID),
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, done := init(tt.fields, t)
			defer done()
			ctx := context.TODO()
			err := s.CreateDashboard(ctx, tt.args.dashboard)
			if (err != nil) != (tt.wants.err != nil) {
				t.Fatalf("expected error '%v' got '%v'", tt.wants.err, err)
			}

			if err != nil && tt.wants.err != nil {
				if err.Error() != tt.wants.err.Error() {
					t.Fatalf("expected error messages to match '%v' got '%v'", tt.wants.err, err.Error())
				}
			}
			defer s.DeleteDashboard(ctx, tt.args.dashboard.ID)

			dashboards, _, err := s.FindDashboards(ctx, platform.DashboardFilter{})
			if err != nil {
				t.Fatalf("failed to retrieve dashboards: %v", err)
			}
			if diff := cmp.Diff(dashboards, tt.wants.dashboards, dashboardCmpOptions...); diff != "" {
				t.Errorf("dashboards are different -got/+want\ndiff %s", diff)
			}
		})
	}
}

// FindDashboardByID testing
func FindDashboardByID(
	init func(DashboardFields, *testing.T) (platform.DashboardService, func()),
	t *testing.T,
) {
	type args struct {
		id platform.ID
	}
	type wants struct {
		err       error
		dashboard *platform.Dashboard
	}

	tests := []struct {
		name   string
		fields DashboardFields
		args   args
		wants  wants
	}{
		{
			name: "basic find dashboard by id",
			fields: DashboardFields{
				Dashboards: []*platform.Dashboard{
					{
						ID:             MustIDFromString(dashOneID),
						OrganizationID: MustIDFromString(orgOneID),
						Name:           "dashboard1",
					},
					{
						ID:             MustIDFromString(dashTwoID),
						OrganizationID: MustIDFromString(orgOneID),
						Name:           "dashboard2",
					},
				},
				Organizations: []*platform.Organization{
					{
						Name: "theorg",
						ID:   MustIDFromString(orgOneID),
					},
				},
			},
			args: args{
				id: MustIDFromString(dashTwoID),
			},
			wants: wants{
				dashboard: &platform.Dashboard{
					ID:             MustIDFromString(dashTwoID),
					OrganizationID: MustIDFromString(orgOneID),
					Organization:   "theorg",
					Name:           "dashboard2",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, done := init(tt.fields, t)
			defer done()
			ctx := context.TODO()

			dashboard, err := s.FindDashboardByID(ctx, tt.args.id)
			if (err != nil) != (tt.wants.err != nil) {
				t.Fatalf("expected errors to be equal '%v' got '%v'", tt.wants.err, err)
			}

			if err != nil && tt.wants.err != nil {
				if err.Error() != tt.wants.err.Error() {
					t.Fatalf("expected error '%v' got '%v'", tt.wants.err, err)
				}
			}

			if diff := cmp.Diff(dashboard, tt.wants.dashboard, dashboardCmpOptions...); diff != "" {
				t.Errorf("dashboard is different -got/+want\ndiff %s", diff)
			}
		})
	}
}

// FindDashboardsByOrganizationName tests.
func FindDashboardsByOrganizationName(
	init func(DashboardFields, *testing.T) (platform.DashboardService, func()),
	t *testing.T,
) {
	type args struct {
		organization string
	}

	type wants struct {
		dashboards []*platform.Dashboard
		err        error
	}
	tests := []struct {
		name   string
		fields DashboardFields
		args   args
		wants  wants
	}{
		{
			name: "find dashboards by organization name",
			fields: DashboardFields{
				Organizations: []*platform.Organization{
					{
						Name: "theorg",
						ID:   MustIDFromString(orgOneID),
					},
					{
						Name: "otherorg",
						ID:   MustIDFromString(orgTwoID),
					},
				},
				Dashboards: []*platform.Dashboard{
					{
						ID:             MustIDFromString(dashOneID),
						OrganizationID: MustIDFromString(orgOneID),
						Name:           "abc",
					},
					{
						ID:             MustIDFromString(dashTwoID),
						OrganizationID: MustIDFromString(orgTwoID),
						Name:           "xyz",
					},
					{
						ID:             MustIDFromString(dashThreeID),
						OrganizationID: MustIDFromString(orgOneID),
						Name:           "123",
					},
				},
			},
			args: args{
				organization: "theorg",
			},
			wants: wants{
				dashboards: []*platform.Dashboard{
					{
						ID:             MustIDFromString(dashOneID),
						OrganizationID: MustIDFromString(orgOneID),
						Organization:   "theorg",
						Name:           "abc",
					},
					{
						ID:             MustIDFromString(dashThreeID),
						OrganizationID: MustIDFromString(orgOneID),
						Organization:   "theorg",
						Name:           "123",
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, done := init(tt.fields, t)
			defer done()
			ctx := context.TODO()
			dashboards, _, err := s.FindDashboardsByOrganizationName(ctx, tt.args.organization)
			if (err != nil) != (tt.wants.err != nil) {
				t.Fatalf("expected errors to be equal '%v' got '%v'", tt.wants.err, err)
			}

			if err != nil && tt.wants.err != nil {
				if err.Error() != tt.wants.err.Error() {
					t.Fatalf("expected error '%v' got '%v'", tt.wants.err, err)
				}
			}

			if diff := cmp.Diff(dashboards, tt.wants.dashboards, dashboardCmpOptions...); diff != "" {
				t.Errorf("dashboards are different -got/+want\ndiff %s", diff)
			}
		})
	}
}

// FindDashboardsByOrganiztionID tests.
func FindDashboardsByOrganizationID(
	init func(DashboardFields, *testing.T) (platform.DashboardService, func()),
	t *testing.T,
) {
	type args struct {
		organizationID platform.ID
	}

	type wants struct {
		dashboards []*platform.Dashboard
		err        error
	}
	tests := []struct {
		name   string
		fields DashboardFields
		args   args
		wants  wants
	}{
		{
			name: "find dashboards by organization id",
			fields: DashboardFields{
				Organizations: []*platform.Organization{
					{
						Name: "theorg",
						ID:   MustIDFromString(orgOneID),
					},
					{
						Name: "otherorg",
						ID:   MustIDFromString(orgTwoID),
					},
				},
				Dashboards: []*platform.Dashboard{
					{
						ID:             MustIDFromString(dashOneID),
						OrganizationID: MustIDFromString(orgOneID),
						Name:           "abc",
					},
					{
						ID:             MustIDFromString(dashTwoID),
						OrganizationID: MustIDFromString(orgTwoID),
						Name:           "xyz",
					},
					{
						ID:             MustIDFromString(dashThreeID),
						OrganizationID: MustIDFromString(orgOneID),
						Name:           "123",
					},
				},
			},
			args: args{
				organizationID: MustIDFromString(orgOneID),
			},
			wants: wants{
				dashboards: []*platform.Dashboard{
					{
						ID:             MustIDFromString(dashOneID),
						OrganizationID: MustIDFromString(orgOneID),
						Organization:   "theorg",
						Name:           "abc",
					},
					{
						ID:             MustIDFromString(dashThreeID),
						OrganizationID: MustIDFromString(orgOneID),
						Organization:   "theorg",
						Name:           "123",
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, done := init(tt.fields, t)
			defer done()
			ctx := context.TODO()
			id := tt.args.organizationID

			dashboards, _, err := s.FindDashboardsByOrganizationID(ctx, id)
			if (err != nil) != (tt.wants.err != nil) {
				t.Fatalf("expected errors to be equal '%v' got '%v'", tt.wants.err, err)
			}

			if err != nil && tt.wants.err != nil {
				if err.Error() != tt.wants.err.Error() {
					t.Fatalf("expected error '%v' got '%v'", tt.wants.err, err)
				}
			}

			if diff := cmp.Diff(dashboards, tt.wants.dashboards, dashboardCmpOptions...); diff != "" {
				t.Errorf("dashboards are different -got/+want\ndiff %s", diff)
			}
		})
	}
}

// FindDashboards testing
func FindDashboards(
	init func(DashboardFields, *testing.T) (platform.DashboardService, func()),
	t *testing.T,
) {
	type args struct {
		ID             platform.ID
		name           string
		organization   string
		organizationID platform.ID
	}

	type wants struct {
		dashboards []*platform.Dashboard
		err        error
	}
	tests := []struct {
		name   string
		fields DashboardFields
		args   args
		wants  wants
	}{
		{
			name: "find all dashboards",
			fields: DashboardFields{
				Organizations: []*platform.Organization{
					{
						Name: "theorg",
						ID:   MustIDFromString(orgOneID),
					},
					{
						Name: "otherorg",
						ID:   MustIDFromString(orgTwoID),
					},
				},
				Dashboards: []*platform.Dashboard{
					{
						ID:             MustIDFromString(dashOneID),
						OrganizationID: MustIDFromString(orgOneID),
						Name:           "abc",
					},
					{
						ID:             MustIDFromString(dashTwoID),
						OrganizationID: MustIDFromString(orgTwoID),
						Name:           "xyz",
					},
				},
			},
			args: args{},
			wants: wants{
				dashboards: []*platform.Dashboard{
					{
						ID:             MustIDFromString(dashOneID),
						OrganizationID: MustIDFromString(orgOneID),
						Organization:   "theorg",
						Name:           "abc",
					},
					{
						ID:             MustIDFromString(dashTwoID),
						OrganizationID: MustIDFromString(orgTwoID),
						Organization:   "otherorg",
						Name:           "xyz",
					},
				},
			},
		},
		{
			name: "find dashboards by organization name",
			fields: DashboardFields{
				Organizations: []*platform.Organization{
					{
						Name: "theorg",
						ID:   MustIDFromString(orgOneID),
					},
					{
						Name: "otherorg",
						ID:   MustIDFromString(orgTwoID),
					},
				},
				Dashboards: []*platform.Dashboard{
					{
						ID:             MustIDFromString(dashOneID),
						OrganizationID: MustIDFromString(orgOneID),
						Name:           "abc",
					},
					{
						ID:             MustIDFromString(dashTwoID),
						OrganizationID: MustIDFromString(orgTwoID),
						Name:           "xyz",
					},
					{
						ID:             MustIDFromString(dashThreeID),
						OrganizationID: MustIDFromString(orgOneID),
						Name:           "123",
					},
				},
			},
			args: args{
				organization: "theorg",
			},
			wants: wants{
				dashboards: []*platform.Dashboard{
					{
						ID:             MustIDFromString(dashOneID),
						OrganizationID: MustIDFromString(orgOneID),
						Organization:   "theorg",
						Name:           "abc",
					},
					{
						ID:             MustIDFromString(dashThreeID),
						OrganizationID: MustIDFromString(orgOneID),
						Organization:   "theorg",
						Name:           "123",
					},
				},
			},
		},
		{
			name: "find dashboards by organization id",
			fields: DashboardFields{
				Organizations: []*platform.Organization{
					{
						Name: "theorg",
						ID:   MustIDFromString(orgOneID),
					},
					{
						Name: "otherorg",
						ID:   MustIDFromString(orgTwoID),
					},
				},
				Dashboards: []*platform.Dashboard{
					{
						ID:             MustIDFromString(dashOneID),
						OrganizationID: MustIDFromString(orgOneID),
						Name:           "abc",
					},
					{
						ID:             MustIDFromString(dashTwoID),
						OrganizationID: MustIDFromString(orgTwoID),
						Name:           "xyz",
					},
					{
						ID:             MustIDFromString(dashThreeID),
						OrganizationID: MustIDFromString(orgOneID),
						Name:           "123",
					},
				},
			},
			args: args{
				organizationID: MustIDFromString(orgOneID),
			},
			wants: wants{
				dashboards: []*platform.Dashboard{
					{
						ID:             MustIDFromString(dashOneID),
						OrganizationID: MustIDFromString(orgOneID),
						Organization:   "theorg",
						Name:           "abc",
					},
					{
						ID:             MustIDFromString(dashThreeID),
						OrganizationID: MustIDFromString(orgOneID),
						Organization:   "theorg",
						Name:           "123",
					},
				},
			},
		},
		{
			name: "find dashboard by id",
			fields: DashboardFields{
				Organizations: []*platform.Organization{
					{
						Name: "theorg",
						ID:   MustIDFromString(orgOneID),
					},
				},
				Dashboards: []*platform.Dashboard{
					{
						ID:             MustIDFromString(dashOneID),
						OrganizationID: MustIDFromString(orgOneID),
						Name:           "abc",
					},
					{
						ID:             MustIDFromString(dashTwoID),
						OrganizationID: MustIDFromString(orgOneID),
						Name:           "xyz",
					},
				},
			},
			args: args{
				ID: MustIDFromString(dashTwoID),
			},
			wants: wants{
				dashboards: []*platform.Dashboard{
					{
						ID:             MustIDFromString(dashTwoID),
						OrganizationID: MustIDFromString(orgOneID),
						Organization:   "theorg",
						Name:           "xyz",
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, done := init(tt.fields, t)
			defer done()
			ctx := context.TODO()

			filter := platform.DashboardFilter{}
			if tt.args.ID.Valid() {
				filter.ID = &tt.args.ID
			}
			if tt.args.organizationID.Valid() {
				filter.OrganizationID = &tt.args.organizationID
			}
			if tt.args.organization != "" {
				filter.Organization = &tt.args.organization
			}

			dashboards, _, err := s.FindDashboards(ctx, filter)
			if (err != nil) != (tt.wants.err != nil) {
				t.Fatalf("expected errors to be equal '%v' got '%v'", tt.wants.err, err)
			}

			if err != nil && tt.wants.err != nil {
				if err.Error() != tt.wants.err.Error() {
					t.Fatalf("expected error '%v' got '%v'", tt.wants.err, err)
				}
			}

			if diff := cmp.Diff(dashboards, tt.wants.dashboards, dashboardCmpOptions...); diff != "" {
				t.Errorf("dashboards are different -got/+want\ndiff %s", diff)
			}
		})
	}
}

// DeleteDashboard testing
func DeleteDashboard(
	init func(DashboardFields, *testing.T) (platform.DashboardService, func()),
	t *testing.T,
) {
	type args struct {
		ID platform.ID
	}
	type wants struct {
		err        error
		dashboards []*platform.Dashboard
	}

	tests := []struct {
		name   string
		fields DashboardFields
		args   args
		wants  wants
	}{
		{
			name: "delete dashboards using exist id",
			fields: DashboardFields{
				Organizations: []*platform.Organization{
					{
						Name: "theorg",
						ID:   MustIDFromString(orgOneID),
					},
				},
				Dashboards: []*platform.Dashboard{
					{
						Name:           "A",
						ID:             MustIDFromString(dashOneID),
						OrganizationID: MustIDFromString(orgOneID),
					},
					{
						Name:           "B",
						ID:             MustIDFromString(dashTwoID),
						OrganizationID: MustIDFromString(orgOneID),
					},
				},
			},
			args: args{
				ID: MustIDFromString(dashOneID),
			},
			wants: wants{
				dashboards: []*platform.Dashboard{
					{
						Name:           "B",
						ID:             MustIDFromString(dashTwoID),
						OrganizationID: MustIDFromString(orgOneID),
						Organization:   "theorg",
					},
				},
			},
		},
		{
			name: "delete dashboards using id that does not exist",
			fields: DashboardFields{
				Organizations: []*platform.Organization{
					{
						Name: "theorg",
						ID:   MustIDFromString(orgOneID),
					},
				},
				Dashboards: []*platform.Dashboard{
					{
						Name:           "A",
						ID:             MustIDFromString(dashOneID),
						OrganizationID: MustIDFromString(orgOneID),
					},
					{
						Name:           "B",
						ID:             MustIDFromString(dashTwoID),
						OrganizationID: MustIDFromString(orgOneID),
					},
				},
			},
			args: args{
				ID: MustIDFromString(dashThreeID),
			},
			wants: wants{
				err: fmt.Errorf("dashboard not found"),
				dashboards: []*platform.Dashboard{
					{
						Name:           "A",
						ID:             MustIDFromString(dashOneID),
						OrganizationID: MustIDFromString(orgOneID),
						Organization:   "theorg",
					},
					{
						Name:           "B",
						ID:             MustIDFromString(dashTwoID),
						OrganizationID: MustIDFromString(orgOneID),
						Organization:   "theorg",
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, done := init(tt.fields, t)
			defer done()
			ctx := context.TODO()
			err := s.DeleteDashboard(ctx, tt.args.ID)
			if (err != nil) != (tt.wants.err != nil) {
				t.Fatalf("expected error '%v' got '%v'", tt.wants.err, err)
			}

			if err != nil && tt.wants.err != nil {
				if err.Error() != tt.wants.err.Error() {
					t.Fatalf("expected error messages to match '%v' got '%v'", tt.wants.err, err.Error())
				}
			}

			filter := platform.DashboardFilter{}
			dashboards, _, err := s.FindDashboards(ctx, filter)
			if err != nil {
				t.Fatalf("failed to retrieve dashboards: %v", err)
			}
			if diff := cmp.Diff(dashboards, tt.wants.dashboards, dashboardCmpOptions...); diff != "" {
				t.Errorf("dashboards are different -got/+want\ndiff %s", diff)
			}
		})
	}
}

// UpdateDashboard testing
func UpdateDashboard(
	init func(DashboardFields, *testing.T) (platform.DashboardService, func()),
	t *testing.T,
) {
	type args struct {
		name      string
		id        platform.ID
		retention int
	}
	type wants struct {
		err       error
		dashboard *platform.Dashboard
	}

	tests := []struct {
		name   string
		fields DashboardFields
		args   args
		wants  wants
	}{
		{
			name: "update name",
			fields: DashboardFields{
				Organizations: []*platform.Organization{
					{
						Name: "theorg",
						ID:   MustIDFromString(orgOneID),
					},
				},
				Dashboards: []*platform.Dashboard{
					{
						ID:             MustIDFromString(dashOneID),
						OrganizationID: MustIDFromString(orgOneID),
						Name:           "dashboard1",
					},
					{
						ID:             MustIDFromString(dashTwoID),
						OrganizationID: MustIDFromString(orgOneID),
						Name:           "dashboard2",
					},
				},
			},
			args: args{
				id:   MustIDFromString(dashOneID),
				name: "changed",
			},
			wants: wants{
				dashboard: &platform.Dashboard{
					ID:             MustIDFromString(dashOneID),
					OrganizationID: MustIDFromString(orgOneID),
					Organization:   "theorg",
					Name:           "changed",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, done := init(tt.fields, t)
			defer done()
			ctx := context.TODO()

			upd := platform.DashboardUpdate{}
			if tt.args.name != "" {
				upd.Name = &tt.args.name
			}

			dashboard, err := s.UpdateDashboard(ctx, tt.args.id, upd)
			if (err != nil) != (tt.wants.err != nil) {
				t.Fatalf("expected error '%v' got '%v'", tt.wants.err, err)
			}

			if err != nil && tt.wants.err != nil {
				if err.Error() != tt.wants.err.Error() {
					t.Fatalf("expected error messages to match '%v' got '%v'", tt.wants.err, err.Error())
				}
			}

			if diff := cmp.Diff(dashboard, tt.wants.dashboard, dashboardCmpOptions...); diff != "" {
				t.Errorf("dashboard is different -got/+want\ndiff %s", diff)
			}
		})
	}
}

// AddDashboardCell tests.
func AddDashboardCell(
	init func(DashboardFields, *testing.T) (platform.DashboardService, func()),
	t *testing.T,
) {
	type args struct {
		dashboardID platform.ID
		cell        *platform.DashboardCell
	}

	type wants struct {
		dashboards []*platform.Dashboard
		err        error
	}
	tests := []struct {
		name   string
		fields DashboardFields
		args   args
		wants  wants
	}{
		{
			name: "add dashboard cell",
			fields: DashboardFields{
				IDGenerator: mock.NewIDGenerator(dashOneID, t),
				Organizations: []*platform.Organization{
					{
						Name: "theorg",
						ID:   MustIDFromString(orgOneID),
					},
				},
				Dashboards: []*platform.Dashboard{
					{
						ID:             MustIDFromString(dashOneID),
						OrganizationID: MustIDFromString(orgOneID),
						Name:           "abc",
					},
				},
			},
			args: args{
				dashboardID: MustIDFromString(dashOneID),
				cell: &platform.DashboardCell{
					DashboardCellContents: platform.DashboardCellContents{
						Name: "hello",
						X:    10,
						Y:    10,
						W:    100,
						H:    12,
					},
					Visualization: platform.CommonVisualization{
						Query: "SELECT * FROM foo",
					},
				},
			},
			wants: wants{
				dashboards: []*platform.Dashboard{
					{
						ID:             MustIDFromString(dashOneID),
						OrganizationID: MustIDFromString(orgOneID),
						Organization:   "theorg",
						Name:           "abc",
						Cells: []platform.DashboardCell{
							{
								DashboardCellContents: platform.DashboardCellContents{
									ID:   MustIDFromString(dashOneID),
									Name: "hello",
									X:    10,
									Y:    10,
									W:    100,
									H:    12,
								},
								Visualization: platform.CommonVisualization{
									Query: "SELECT * FROM foo",
								},
							},
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, done := init(tt.fields, t)
			defer done()
			ctx := context.TODO()
			err := s.AddDashboardCell(ctx, tt.args.dashboardID, tt.args.cell)
			if (err != nil) != (tt.wants.err != nil) {
				t.Fatalf("expected errors to be equal '%v' got '%v'", tt.wants.err, err)
			}

			if err != nil && tt.wants.err != nil {
				if err.Error() != tt.wants.err.Error() {
					t.Fatalf("expected error '%v' got '%v'", tt.wants.err, err)
				}
			}

			dashboards, _, err := s.FindDashboards(ctx, platform.DashboardFilter{})
			if err != nil {
				t.Fatalf("failed to retrieve dashboards: %v", err)
			}
			if diff := cmp.Diff(dashboards, tt.wants.dashboards, dashboardCmpOptions...); diff != "" {
				t.Errorf("dashboards are different -got/+want\ndiff %s", diff)
			}
		})
	}
}

// RemoveDashboardCell tests.
func RemoveDashboardCell(
	init func(DashboardFields, *testing.T) (platform.DashboardService, func()),
	t *testing.T,
) {
	type args struct {
		dashboardID platform.ID
		cellID      platform.ID
	}

	type wants struct {
		dashboards []*platform.Dashboard
		err        error
	}
	tests := []struct {
		name   string
		fields DashboardFields
		args   args
		wants  wants
	}{
		{
			name: "remove dashboard cell",
			fields: DashboardFields{
				IDGenerator: mock.NewIDGenerator(dashOneID, t),
				Organizations: []*platform.Organization{
					{
						Name: "theorg",
						ID:   MustIDFromString(orgOneID),
					},
				},
				Dashboards: []*platform.Dashboard{
					{
						ID:             MustIDFromString(dashOneID),
						OrganizationID: MustIDFromString(orgOneID),
						Name:           "abc",
						Cells: []platform.DashboardCell{
							{
								DashboardCellContents: platform.DashboardCellContents{
									ID:   MustIDFromString(dashOneID),
									Name: "hello",
									X:    10,
									Y:    10,
									W:    100,
									H:    12,
								},
								Visualization: platform.CommonVisualization{
									Query: "SELECT * FROM foo",
								},
							},
							{
								DashboardCellContents: platform.DashboardCellContents{
									ID:   MustIDFromString(dashTwoID),
									Name: "world",
									X:    10,
									Y:    10,
									W:    100,
									H:    12,
								},
								Visualization: platform.CommonVisualization{
									Query: "SELECT * FROM mem",
								},
							},
							{
								DashboardCellContents: platform.DashboardCellContents{
									ID:   MustIDFromString(dashThreeID),
									Name: "bar",
									X:    10,
									Y:    10,
									W:    101,
									H:    12,
								},
								Visualization: platform.CommonVisualization{
									Query: "SELECT thing FROM foo",
								},
							},
						},
					},
				},
			},
			args: args{
				dashboardID: MustIDFromString(dashOneID),
				cellID:      MustIDFromString(dashTwoID),
			},
			wants: wants{
				dashboards: []*platform.Dashboard{
					{
						ID:             MustIDFromString(dashOneID),
						OrganizationID: MustIDFromString(orgOneID),
						Organization:   "theorg",
						Name:           "abc",
						Cells: []platform.DashboardCell{
							{
								DashboardCellContents: platform.DashboardCellContents{
									ID:   MustIDFromString(dashOneID),
									Name: "hello",
									X:    10,
									Y:    10,
									W:    100,
									H:    12,
								},
								Visualization: platform.CommonVisualization{
									Query: "SELECT * FROM foo",
								},
							},
							{
								DashboardCellContents: platform.DashboardCellContents{
									ID:   MustIDFromString(dashThreeID),
									Name: "bar",
									X:    10,
									Y:    10,
									W:    101,
									H:    12,
								},
								Visualization: platform.CommonVisualization{
									Query: "SELECT thing FROM foo",
								},
							},
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, done := init(tt.fields, t)
			defer done()
			ctx := context.TODO()
			err := s.RemoveDashboardCell(ctx, tt.args.dashboardID, tt.args.cellID)
			if (err != nil) != (tt.wants.err != nil) {
				t.Fatalf("expected errors to be equal '%v' got '%v'", tt.wants.err, err)
			}

			if err != nil && tt.wants.err != nil {
				if err.Error() != tt.wants.err.Error() {
					t.Fatalf("expected error '%v' got '%v'", tt.wants.err, err)
				}
			}

			dashboards, _, err := s.FindDashboards(ctx, platform.DashboardFilter{})
			if err != nil {
				t.Fatalf("failed to retrieve dashboards: %v", err)
			}
			if diff := cmp.Diff(dashboards, tt.wants.dashboards, dashboardCmpOptions...); diff != "" {
				t.Errorf("dashboards are different -got/+want\ndiff %s", diff)
			}
		})
	}
}

// ReplaceDashboardCell tests.
func ReplaceDashboardCell(
	init func(DashboardFields, *testing.T) (platform.DashboardService, func()),
	t *testing.T,
) {
	type args struct {
		dashboardID platform.ID
		cell        *platform.DashboardCell
	}

	type wants struct {
		dashboards []*platform.Dashboard
		err        error
	}
	tests := []struct {
		name   string
		fields DashboardFields
		args   args
		wants  wants
	}{
		{
			name: "add dashboard cell",
			fields: DashboardFields{
				Organizations: []*platform.Organization{
					{
						Name: "theorg",
						ID:   MustIDFromString(orgOneID),
					},
				},
				Dashboards: []*platform.Dashboard{
					{
						ID:             MustIDFromString(dashOneID),
						OrganizationID: MustIDFromString(orgOneID),
						Name:           "abc",
						Cells: []platform.DashboardCell{
							{
								DashboardCellContents: platform.DashboardCellContents{
									ID:   MustIDFromString(dashOneID),
									Name: "hello",
									X:    10,
									Y:    10,
									W:    100,
									H:    12,
								},
								Visualization: platform.CommonVisualization{
									Query: "SELECT * FROM foo",
								},
							},
						},
					},
				},
			},
			args: args{
				dashboardID: MustIDFromString(dashOneID),
				cell: &platform.DashboardCell{
					DashboardCellContents: platform.DashboardCellContents{
						ID:   MustIDFromString(dashOneID),
						Name: "what",
						X:    101,
						Y:    102,
						W:    100,
						H:    12,
					},
					Visualization: platform.CommonVisualization{
						Query: "SELECT * FROM thing",
					},
				},
			},
			wants: wants{
				dashboards: []*platform.Dashboard{
					{
						ID:             MustIDFromString(dashOneID),
						OrganizationID: MustIDFromString(orgOneID),
						Organization:   "theorg",
						Name:           "abc",
						Cells: []platform.DashboardCell{
							{
								DashboardCellContents: platform.DashboardCellContents{
									ID:   MustIDFromString(dashOneID),
									Name: "what",
									X:    101,
									Y:    102,
									W:    100,
									H:    12,
								},
								Visualization: platform.CommonVisualization{
									Query: "SELECT * FROM thing",
								},
							},
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, done := init(tt.fields, t)
			defer done()
			ctx := context.TODO()
			err := s.ReplaceDashboardCell(ctx, tt.args.dashboardID, tt.args.cell)
			if (err != nil) != (tt.wants.err != nil) {
				t.Fatalf("expected errors to be equal '%v' got '%v'", tt.wants.err, err)
			}

			if err != nil && tt.wants.err != nil {
				if err.Error() != tt.wants.err.Error() {
					t.Fatalf("expected error '%v' got '%v'", tt.wants.err, err)
				}
			}

			dashboards, _, err := s.FindDashboards(ctx, platform.DashboardFilter{})
			if err != nil {
				t.Fatalf("failed to retrieve dashboards: %v", err)
			}
			if diff := cmp.Diff(dashboards, tt.wants.dashboards, dashboardCmpOptions...); diff != "" {
				t.Errorf("dashboards are different -got/+want\ndiff %s", diff)
			}
		})
	}
}
