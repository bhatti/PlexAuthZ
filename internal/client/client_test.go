package client

import (
	"fmt"
	"github.com/bhatti/PlexAuthZ/api/v1/types"
	"github.com/bhatti/PlexAuthZ/internal/authz"
	"github.com/bhatti/PlexAuthZ/internal/controller"
	"github.com/bhatti/PlexAuthZ/internal/domain"
	"github.com/bhatti/PlexAuthZ/internal/factory"
	"github.com/bhatti/PlexAuthZ/internal/metrics"
	"github.com/bhatti/PlexAuthZ/internal/server"
	"github.com/bhatti/PlexAuthZ/internal/service"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
	"time"
)

func Test_AuthZPermissions(t *testing.T) {
	runTests(t,
		metrics.New(),
		testPermissionsForDepositAccount,
		testPermissionsForProjects,
		testPermissionsWithScopeAndWithinGeoLocation,
		testForDetectingAmbiguousPermission,
		testForPermissionsForFeatureFlagsWithGeoFencing,
		testForPermissionsWithWildcardActions,
		testForPermissionsWithQuota,
		testReBACPermissions,
		testAttributeBasedPermissions,
		testForPermissionsWithAttributes,
		testRBAC,
		testCRUD,
		testForPermissionsWithIPAddresses,
		testPermissionsWithOwners,
	)
}

func testPermissionsWithOwners(
	t *testing.T,
	authAdapter *AuthAdapter,
) {
	// create org
	orgAdapter, err := authAdapter.CreateOrganization(
		&types.Organization{
			Name:       "xyz-corp",
			Namespaces: []string{"marketing", "sales"},
		})
	require.NoError(t, err)
	namespace := orgAdapter.Organization.Namespaces[0]

	// AND with following principals
	alice, err := orgAdapter.Principals().
		WithUsername("alice").
		WithAttributes("Department", "Engineering", "Rank", "5").Create()
	require.NoError(t, err)

	bob, err := orgAdapter.Principals().
		WithUsername("bob").
		WithAttributes("Department", "Engineering", "Rank", "6").Create()
	require.NoError(t, err)

	charlie, err := orgAdapter.Principals().
		WithUsername("charlie").
		WithAttributes("Department", "Sales", "Rank", "6").Create()
	require.NoError(t, err)

	// AND with following resources
	app, err := orgAdapter.Resources(namespace).
		WithName("ios-app").
		WithAttributes("Editors", "alice bob").
		WithActions("list", "read", "write", "create", "delete").Create()
	require.NoError(t, err)

	// AND with following permissions
	rlPerm1, err := orgAdapter.Permissions(namespace).
		WithResource(app.Resource).
		WithConstraints(`
	{{or (Includes .Resource.Editors .Principal.Username) (GE .Principal.Rank 6)}}
`).WithActions("read", "list").Create()
	require.NoError(t, err)

	wPerm2, err := orgAdapter.Permissions(namespace).
		WithResource(app.Resource).
		WithConstraints(`
	{{and (Includes .Resource.Editors .Principal.Username) (GE .Principal.Rank 6)}}
`).WithActions("write").Create()
	require.NoError(t, err)

	// assigning permission to all three principals
	require.NoError(t, alice.AddPermissions(rlPerm1.Permission, wPerm2.Permission))
	require.NoError(t, bob.AddPermissions(rlPerm1.Permission, wPerm2.Permission))
	require.NoError(t, charlie.AddPermissions(rlPerm1.Permission, wPerm2.Permission))

	// Alice, Bob, Charlie should be able to read/list since alice/bob belong to Editors attribute and
	// Charlie's rank >= 6
	require.NoError(t, alice.Authorizer(namespace).
		WithAction("list").
		WithResourceName("ios-app").Check())
	require.NoError(t, bob.Authorizer(namespace).
		WithAction("list").
		WithResourceName("ios-app").Check())
	require.NoError(t, charlie.Authorizer(namespace).
		WithAction("list").
		WithResourceName("ios-app").Check())

	// Only Bob should be able to write because Alice's rank is lower than 6 and Charlie doesn't belong
	// to Editors attribute.
	require.Error(t, alice.Authorizer(namespace).
		WithAction("write").
		WithResourceName("ios-app").Check())
	require.NoError(t, bob.Authorizer(namespace).
		WithAction("write").
		WithResourceName("ios-app").Check())
	require.Error(t, charlie.Authorizer(namespace).
		WithAction("write").
		WithResourceName("ios-app").Check())
}

func testForPermissionsWithIPAddresses(
	t *testing.T,
	authAdapter *AuthAdapter,
) {
	// create org
	orgAdapter, err := authAdapter.CreateOrganization(
		&types.Organization{
			Name:       "xyz-corp",
			Namespaces: []string{"marketing", "sales"},
		})
	require.NoError(t, err)
	namespace := orgAdapter.Organization.Namespaces[0]

	// AND with following principals
	alice, err := orgAdapter.Principals().
		WithUsername("alice").
		WithAttributes("Department", "Engineering", "Permanent", "true").Create()

	// AND with following resources
	project, err := orgAdapter.Resources(namespace).
		WithName("nextgen-app").
		WithActions("list", "read", "write", "create", "delete").Create()
	require.NoError(t, err)

	// AND with following permissions
	rwlPerm, err := orgAdapter.Permissions(namespace).
		WithResource(project.Resource).
		WithEffect(types.Effect_PERMITTED).
		WithConstraints(`
{{$Loopback := IsLoopback .IPAddress}}
{{$Multicast := IsMulticast .IPAddress}}
{{and (not $Loopback) (not $Multicast) (IPInRange .IPAddress "211.211.211.0/24")}}
`).WithActions("read", "write", "list").Create()
	require.NoError(t, err)

	// assigning permission to roles
	require.NoError(t, alice.AddPermissions(rwlPerm.Permission))

	// Project should be only be accessible if ip-address is not loop-back, not multi-cast and within ip-range
	require.NoError(t, alice.Authorizer(namespace).
		WithAction("list").
		WithContext("IPAddress", "211.211.211.5").
		WithResourceName("nextgen-app").Check())
	// But not local ipaddress or multicast
	require.Error(t, alice.Authorizer(namespace).
		WithAction("list").
		WithContext("IPAddress", "127.0.0.1").
		WithResourceName("nextgen-app").Check())
	require.Error(t, alice.Authorizer(namespace).
		WithAction("list").
		WithContext("IPAddress", "224.0.0.1").
		WithResourceName("nextgen-app").Check())
}

func testRBAC(
	t *testing.T,
	authAdapter *AuthAdapter,
) {
	// create org
	orgAdapter, err := authAdapter.CreateOrganization(
		&types.Organization{
			Name:       "bank-abc",
			Namespaces: []string{"Checking", "Loan"},
		})
	require.NoError(t, err)
	namespace := orgAdapter.Organization.Namespaces[0]

	// AND with following principals
	alice, err := orgAdapter.Principals().
		WithUsername("alice").
		WithAttributes("Department", "Personal Banking", "EmploymentLength", "5").Create()
	require.NoError(t, err)

	bob, err := orgAdapter.Principals().
		WithUsername("bob").
		WithAttributes("Department", "Relationship Management", "EmploymentLength", "2").Create()
	require.NoError(t, err)

	charlie, err := orgAdapter.Principals().
		WithUsername("charlie").
		WithAttributes("Department", "Sales", "EmploymentLength", "3").Create()
	require.NoError(t, err)

	teller, err := orgAdapter.Roles(namespace).WithName("Teller").Create()
	require.NoError(t, err)
	manager, err := orgAdapter.Roles(namespace).WithName("Manager").WithParents(teller.Role).Create()
	require.NoError(t, err)
	loanOfficer, err := orgAdapter.Roles(namespace).WithName("LoanOfficer").Create()
	require.NoError(t, err)
	support, err := orgAdapter.Roles(namespace).WithName("ITSupport").Create()
	require.NoError(t, err)

	// Assigning roles
	require.NoError(t, alice.AddRoles(manager.Role))
	require.NoError(t, bob.AddRoles(loanOfficer.Role))
	require.NoError(t, charlie.AddRoles(support.Role))

	sales, err := orgAdapter.Groups(namespace).WithName("Sales").Create()
	require.NoError(t, err)
	accounting, err := orgAdapter.Groups(namespace).WithName("Accounting").Create()
	require.NoError(t, err)
	engineering, err := orgAdapter.Groups(namespace).WithName("Engineering").Create()
	require.NoError(t, err)

	// Assigning groups
	require.NoError(t, alice.AddGroups(sales.Group))
	require.NoError(t, bob.AddGroups(accounting.Group))
	require.NoError(t, charlie.AddGroups(engineering.Group))

	// Test for constraints
	require.NoError(t, alice.Authorizer(namespace).
		WithConstraints(
			`
{{and (HasRole "Teller") (HasGroup "Sales") (TimeInRange .CurrentTime .StartTime .EndTime)}}
`).
		WithContext("CurrentTime", "10:00am", "StartTime", "8:00am", "EndTime", "4:00pm").Check())

	require.NoError(t, bob.Authorizer(namespace).
		WithConstraints(
			`
{{and (HasRole "LoanOfficer") (HasGroup "Accounting") (TimeInRange .CurrentTime .StartTime .EndTime) (GT .Principal.EmploymentLength 1)}}
`).
		WithContext("CurrentTime", "10:00am", "StartTime", "8:00am", "EndTime", "4:00pm").Check())

	require.NoError(t, charlie.Authorizer(namespace).
		WithConstraints(
			`
{{and (HasRole "ITSupport") (HasGroup "Engineering") (TimeInRange .CurrentTime .StartTime .EndTime) (GT .Principal.EmploymentLength 1)}}
`).
		WithContext("CurrentTime", "10:00am", "StartTime", "8:00am", "EndTime", "4:00pm").Check())

	// but should fail for bob because ITSupport Role and Engineering Group is required.
	require.Error(t, bob.Authorizer(namespace).
		WithConstraints(
			`
{{and (HasRole "ITSupport") (HasGroup "Engineering") (TimeInRange .CurrentTime .StartTime .EndTime) (GT .Principal.EmploymentLength 1)}}
`).
		WithContext("CurrentTime", "10:00am", "StartTime", "8:00am", "EndTime", "4:00pm").Check())
}

func testAttributeBasedPermissions(
	t *testing.T,
	authAdapter *AuthAdapter,
) {
	// create org
	orgAdapter, err := authAdapter.CreateOrganization(
		&types.Organization{
			Name:       "bank-abc",
			Namespaces: []string{"Checking", "Loan"},
		})
	require.NoError(t, err)
	namespace := orgAdapter.Organization.Namespaces[0]

	// AND with following principals
	alice, err := orgAdapter.Principals().
		WithUsername("alice").
		WithAttributes("UserRole", "Teller",
			"Department", "Personal Banking", "EmploymentLength", "5").Create()
	require.NoError(t, err)

	bob, err := orgAdapter.Principals().
		WithUsername("bob").
		WithAttributes("UserRole", "CSR",
			"Department", "Relationship Management", "EmploymentLength", "5").Create()
	require.NoError(t, err)

	// AND with following resources
	checkingAccount, err := orgAdapter.Resources(namespace).
		WithName("CheckingAccount").
		WithAttributes("Type", "Checking", "SensitivityLevel", "Normal").
		WithActions("read", "write", "create", "delete").Create()
	require.NoError(t, err)

	secureAccount, err := orgAdapter.Resources(namespace).
		WithName("SecureAccount").
		WithAttributes("Type", "Investment", "SensitivityLevel", "High", "StartTime", "8:00am", "EndTime", "4:00pm").
		WithActions("read", "write", "create", "delete").Create()
	require.NoError(t, err)

	// AND with following permissions
	rwCheckingPerm, err := orgAdapter.Permissions(namespace).
		WithResource(checkingAccount.Resource).
		WithConstraints(`
{{and (eq .Principal.UserRole "Teller") (eq .Resource.SensitivityLevel "Normal")}}
`).WithActions("read", "write").Create()
	require.NoError(t, err)

	allCheckingPerm, err := orgAdapter.Permissions(namespace).
		WithResource(checkingAccount.Resource).
		WithConstraints(`
{{and (eq .Principal.UserRole "CSR") (eq .Resource.SensitivityLevel "Normal")}}
`).WithActions("*").Create()
	require.NoError(t, err)

	rwSecurePerm, err := orgAdapter.Permissions(namespace).
		WithResource(secureAccount.Resource).
		WithConstraints(`
{{and (eq .Principal.UserRole "CSR") (eq .Resource.SensitivityLevel "High") (TimeInRange .CurrentTime .Resource.StartTime .Resource.EndTime)}}
`).WithActions("read", "write").Create()
	require.NoError(t, err)

	allSecurePerm, err := orgAdapter.Permissions(namespace).
		WithResource(secureAccount.Resource).
		WithConstraints(`
{{and (eq .Principal.UserRole "Admin") (eq .Resource.SensitivityLevel "High") (TimeInRange .CurrentTime .Resource.StartTime .Resource.EndTime)}}
`).WithActions("*").Create()
	require.NoError(t, err)
	require.NotNil(t, allSecurePerm)

	// assigning permission to roles
	require.NoError(t, alice.AddPermissions(rwCheckingPerm.Permission))
	require.NoError(t, bob.AddPermissions(allCheckingPerm.Permission))
	require.NoError(t, bob.AddPermissions(rwSecurePerm.Permission))

	// Test for permissions
	require.NoError(t, alice.Authorizer(namespace).
		WithAction("write").
		WithResource(checkingAccount.Resource).Check())

	require.NoError(t, bob.Authorizer(namespace).
		WithAction("create").
		WithResource(checkingAccount.Resource).Check())

	// alice should not create account
	require.Error(t, alice.Authorizer(namespace).
		WithAction("create").
		WithResource(checkingAccount.Resource).Check())

	// Bob should read/write secure account
	require.NoError(t, bob.Authorizer(namespace).
		WithAction("read").
		WithContext("CurrentTime", "10:00am").
		WithResource(secureAccount.Resource).Check())

	// but not create secure account
	require.Error(t, bob.Authorizer(namespace).
		WithAction("create").
		WithContext("CurrentTime", "10:00am").
		WithResource(secureAccount.Resource).Check())

}

func testReBACPermissions(
	t *testing.T,
	authAdapter *AuthAdapter,
) {
	// create org
	orgAdapter, err := authAdapter.CreateOrganization(
		&types.Organization{
			Name:       "hospital",
			Namespaces: []string{"Admissions", "Cardiology", "Oncology", "Radiology"},
		})
	require.NoError(t, err)
	namespace := orgAdapter.Organization.Namespaces[0]

	// AND with following principals
	smith, err := orgAdapter.Principals().
		WithUsername("smith").
		WithName("Dr. Smith").
		WithAttributes("UserRole", "Doctor").Create()
	require.NoError(t, err)

	john, err := orgAdapter.Principals().
		WithUsername("john").
		WithName("John Doe").
		WithAttributes("UserRole", "Patient").Create()
	require.NoError(t, err)

	// AND with following resources
	medicalRecords, err := orgAdapter.Resources(namespace).
		WithName("MedicalRecords").
		WithAttributes("Year", fmt.Sprintf("%d", time.Now().Year()), "Location", "Hospital").
		WithActions("read", "write", "create", "delete").Create()
	require.NoError(t, err)

	docRelation, err := smith.Relationships(namespace).
		WithRelation("AsDoctor").
		WithResource(medicalRecords.Resource).
		WithAttributes("Location", "Hospital").Create()
	require.NoError(t, err)
	require.NotNil(t, docRelation)

	patientRelation, err := john.Relationships(namespace).
		WithRelation("AsPatient").
		WithResource(medicalRecords.Resource).Create()
	require.NoError(t, err)
	require.NotNil(t, patientRelation)

	// AND with following permissions
	rwPerm, err := orgAdapter.Permissions(namespace).
		WithResource(medicalRecords.Resource).
		WithConstraints(`
{{$CurrentYear := TimeNow "2006"}}
{{and (HasRelation "AsDoctor") (DistanceWithinKM .UserLatLng "46.879967,-121.726906" 100) 
(eq .Resource.Year $CurrentYear) (eq .Resource.Location .Location)}}
`).WithActions("read", "write").Create()
	require.NoError(t, err)

	rPerm, err := orgAdapter.Permissions(namespace).
		WithResource(medicalRecords.Resource).
		WithScope("john's records").
		WithConstraints(`
{{$CurrentYear := TimeNow "2006"}}
{{and (HasRelation "AsPatient") (eq .Resource.Year $CurrentYear) (eq .Resource.Location .Location)}}
`).WithActions("read").Create()
	require.NoError(t, err)

	// assigning permission to roles
	require.NoError(t, smith.AddPermissions(rwPerm.Permission))
	require.NoError(t, john.AddPermissions(rPerm.Permission))

	// Test for permissions
	// Dr. Smith should have permission for reading/writing medical records based on constraints
	require.NoError(t, smith.Authorizer(namespace).
		WithAction("write").
		WithResource(medicalRecords.Resource).
		WithContext("UserLatLng", "47.620422,-122.349358", "Location", "Hospital").Check())

	// Patient john should have permission for reading medical records based on constraints
	require.NoError(t, john.Authorizer(namespace).
		WithAction("read").
		WithScope("john's records").
		WithResource(medicalRecords.Resource).
		WithContext("Location", "Hospital").Check())

	// But Patient john should not write medical records
	require.Error(t, john.Authorizer(namespace).
		WithAction("write").
		WithResource(medicalRecords.Resource).
		WithContext("Location", "Hospital").Check())

	// Now treating Doctor as Target Resource for appointment
	doctorResource, err := orgAdapter.Resources(namespace).
		WithName(smith.Principal.Name).
		WithAttributes("Year", fmt.Sprintf("%d", time.Now().Year()),
			"Location", "Hospital").
		WithActions("appointment", "consult").Create()
	require.NoError(t, err)

	doctorPatientRelation, err := john.Relationships(namespace).
		WithRelation("Physician").
		WithAttributes("StartTime", "8:00am", "EndTime", "4:00pm").
		WithResource(doctorResource.Resource).Create()
	require.NoError(t, err)
	require.NotNil(t, doctorPatientRelation)

	apptPerm, err := orgAdapter.Permissions(namespace).
		WithResource(doctorResource.Resource).
		WithConstraints(`
{{$CurrentYear := TimeNow "2006"}}
{{and (TimeInRange .AppointmentTime .Relations.Physician.StartTime .Relations.Physician.EndTime) 
(HasRelation "Physician") (eq "Patient" .Principal.UserRole) (eq .Resource.Year $CurrentYear) (eq .Resource.Location .Location)}}
`).WithActions("appointment").Create()
	require.NoError(t, err)
	require.NoError(t, john.AddPermissions(apptPerm.Permission))

	// Patient john should be able to make appointment within normal Hopspital hours
	require.NoError(t, john.Authorizer(namespace).
		WithAction("appointment").
		WithResource(doctorResource.Resource).
		WithContext("Location", "Hospital", "AppointmentTime", "10:00am").Check())
}

func testForPermissionsWithQuota(
	t *testing.T,
	authAdapter *AuthAdapter,
) {
	// create org
	orgAdapter, err := authAdapter.CreateOrganization(
		&types.Organization{
			Name:       "xyz-corp",
			Namespaces: []string{"marketing", "sales"},
		})
	require.NoError(t, err)
	namespace := orgAdapter.Organization.Namespaces[0]

	engGroup, err := orgAdapter.Groups(namespace).WithName("Engineering").Create()
	require.NoError(t, err)

	// AND with following principals
	alice, err := orgAdapter.Principals().
		WithUsername("alice").
		WithName("Alice").
		WithAttributes("Title", "Engineer", "Tenure", "3").Create()
	require.NoError(t, err)

	bob, err := orgAdapter.Principals().
		WithUsername("bob").
		WithName("Bob").
		WithAttributes("Title", "Manager", "Tenure", "5").Create()
	require.NoError(t, err)

	require.NoError(t, alice.AddGroups(engGroup.Group))
	require.NoError(t, bob.AddGroups(engGroup.Group))

	// AND with following resources
	ideLicences, err := orgAdapter.Resources(namespace).
		WithName("IDELicence").
		WithCapacity(5).
		WithAttributes("Location", "Chicago").
		WithActions("use").Create()
	require.NoError(t, err)

	// allocate IDE licenses
	err = ideLicences.
		WithConstraints(`HasGroup "Engineering"`).
		WithExpiration(time.Hour).
		Allocate(alice.Principal)
	require.NoError(t, err)

	err = ideLicences.
		WithConstraints(`and (GT .Principal.Tenure 1) (HasGroup "Engineering") (eq .Resource.Location .Location)`).
		WithContext("Location", "Chicago").
		WithExpiration(time.Hour).
		Allocate(bob.Principal)
	require.NoError(t, err)

	allocated, err := ideLicences.AllocatedCount()
	require.NoError(t, err)
	require.Equal(t, int32(2), allocated)

	// deallocate IDE licenses
	err = ideLicences.Deallocate(alice.Principal)
	require.NoError(t, err)

	err = ideLicences.Deallocate(bob.Principal)
	require.NoError(t, err)

	allocated, err = ideLicences.AllocatedCount()
	require.NoError(t, err)
	require.Equal(t, int32(0), allocated)
}

func testForPermissionsForFeatureFlagsWithGeoFencing(
	t *testing.T,
	authAdapter *AuthAdapter,
) {
	// create org
	orgAdapter, err := authAdapter.CreateOrganization(
		&types.Organization{
			Name:       "xyz-corp",
			Namespaces: []string{"marketing", "sales"},
		})
	require.NoError(t, err)
	namespace := orgAdapter.Organization.Namespaces[0]

	// AND with following principals
	alice, err := orgAdapter.Principals().
		WithUsername("alice").
		WithEmail("alice@xyz.com").
		WithAttributes("Title", "Engineer", "Permanent", "true").Create()
	require.NoError(t, err)

	bob, err := orgAdapter.Principals().
		WithUsername("bob").
		WithAttributes("Title", "Manager", "Permanent", "true").Create()
	require.NoError(t, err)

	// AND with following resources
	features, err := orgAdapter.Resources(namespace).
		WithName("UIFeatures").
		WithActions("read", "write", "create", "delete").Create()
	require.NoError(t, err)

	// AND with following permissions
	readPerms, err := orgAdapter.Permissions(namespace).
		WithResource(features.Resource).
		WithScope("beta").
		WithConstraints(`
	{{$CurrentYear := TimeNow "2006"}}
	{{and (eq .Principal.Permanent "true") (DistanceWithinKM .UserLatLng "46.879967,-121.726906" 100)}}
`).WithActions("read", "write").Create()
	require.NoError(t, err)
	allPerms, err := orgAdapter.Permissions(namespace).
		WithResource(features.Resource).
		WithScope("*").
		WithConstraints(`
	{{$CurrentYear := TimeNow "2006"}}
	{{and (eq .Principal.Permanent "true") (DistanceWithinKM .UserLatLng "46.879967,-121.726906" 100)}}
`).WithActions("read", "write").Create()
	require.NoError(t, err)

	// AND with following roles
	employee, err := orgAdapter.Roles(namespace).
		WithName("Employee").
		Create()
	require.NoError(t, err)
	customer, err := orgAdapter.Roles(namespace).
		WithName("Customer").
		Create()
	require.NoError(t, err)

	err = alice.AddRoles(employee.Role)
	require.NoError(t, err)
	err = alice.AddRoles(customer.Role)
	require.NoError(t, err)

	// assigning permission to roles
	require.NoError(t, employee.AddPermissions(allPerms.Permission))
	require.NoError(t, customer.AddPermissions(readPerms.Permission))

	// add role to principals
	require.NoError(t, alice.AddRoles(employee.Role))
	require.NoError(t, bob.AddRoles(customer.Role))

	// Test for permissions for features
	// Alice should access permissions
	require.NoError(t, alice.Authorizer(namespace).
		WithAction("write").
		WithScope("svc1").
		WithResourceName("UIFeatures").
		WithContext("UserLatLng", "47.620422,-122.349358").Check())

	// But should have only access to beta read
	require.NoError(t, bob.Authorizer(namespace).
		WithAction("read").
		WithScope("beta").
		WithResourceName("UIFeatures").
		WithContext("UserLatLng", "47.620422,-122.349358").Check())

	// Bob should not have access for any other scope
	require.Error(t, bob.Authorizer(namespace).
		WithAction("read").
		WithScope("prod").
		WithResourceName("UIFeatures").
		WithContext("UserLatLng", "47.620422,-122.349358").Check())
	require.NoError(t, employee.DeletePermissions(allPerms.Permission))
	require.NoError(t, customer.DeletePermissions(readPerms.Permission))
}

func testForPermissionsWithWildcardActions(
	t *testing.T,
	authAdapter *AuthAdapter,
) {
	// create org
	orgAdapter, err := authAdapter.CreateOrganization(
		&types.Organization{
			Name:       "xyz-corp",
			Namespaces: []string{"marketing", "sales"},
		})
	require.NoError(t, err)
	namespace := orgAdapter.Organization.Namespaces[0]

	// AND with following principals
	alice, err := orgAdapter.Principals().
		WithUsername("alice").
		WithAttributes("Department", "Sales", "Rank", "6").Create()

	bob, err := orgAdapter.Principals().
		WithUsername("bob").
		WithAttributes("Department", "Engineering", "Rank", "6").Create()

	// Creating a project with wildcard
	salesProject, err := orgAdapter.Resources(namespace).
		WithName("urn:org-sales-*-project-1000-*").
		WithAttributes("SalesYear", fmt.Sprintf("%d", time.Now().Year())).
		WithActions("list", "read", "write", "create", "delete").Create()
	require.NoError(t, err)

	// AND with following permissions
	rwlPerm, err := orgAdapter.Permissions(namespace).
		WithResource(salesProject.Resource).
		WithEffect(types.Effect_PERMITTED).
		WithConstraints(`
	{{$CurrentYear := TimeNow "2006"}}
	{{and (GT .Principal.Rank 5) (eq .Principal.Department "Sales") (IPInRange .IPAddress "211.211.211.0/24") (eq .Resource.SalesYear $CurrentYear)}}
`).WithActions("*").Create()
	require.NoError(t, err)

	// assigning permission to roles
	require.NoError(t, alice.AddPermissions(rwlPerm.Permission))
	require.NoError(t, bob.AddPermissions(rwlPerm.Permission))

	// Test for list/read/write permissions
	// Alice should be able to access from Sales Department and complete project name
	require.NoError(t, alice.Authorizer(namespace).
		WithAction("list").
		WithResourceName("urn:org-sales-abc-project-1000-xyz").
		WithContext("IPAddress", "211.211.211.5").Check())
	require.NoError(t, alice.Authorizer(namespace).
		WithAction("read").
		WithResourceName("urn:org-sales-abc-project-1000-xyz").
		WithContext("IPAddress", "211.211.211.5").Check())
	require.NoError(t, alice.Authorizer(namespace).
		WithAction("write").
		WithResourceName("urn:org-sales-abc-project-1000-xyz").
		WithContext("IPAddress", "211.211.211.5").Check())
	// But bob should not be able to access project because he doesn't belong to the Sales Department
	require.Error(t, bob.Authorizer(namespace).
		WithAction("read").
		WithResourceName("urn:org-sales-abc-project-1000-xyz").
		WithContext("IPAddress", "211.211.211.5").Check())

	// But unknown action should fail
	require.Error(t, alice.Authorizer(namespace).
		WithAction("unknown").
		WithResourceName("urn:org-sales-abc-project-1000-xyz").
		WithContext("IPAddress", "211.211.211.5").Check())
}

func testForPermissionsWithAttributes(
	t *testing.T,
	authAdapter *AuthAdapter,
) {
	// create org
	orgAdapter, err := authAdapter.CreateOrganization(
		&types.Organization{
			Name:       "xyz-corp",
			Namespaces: []string{"marketing", "sales"},
		})
	require.NoError(t, err)
	namespace := orgAdapter.Organization.Namespaces[0]

	// AND with following principals
	alice, err := orgAdapter.Principals().
		WithUsername("alice").
		WithAttributes("Department", "Engineering", "Permanent", "true").Create()

	bob, err := orgAdapter.Principals().
		WithUsername("bob").
		WithAttributes("Department", "Sales", "Permanent", "true").Create()

	// AND with following resources
	project, err := orgAdapter.Resources(namespace).
		WithName("nextgen-app").
		WithAttributes("Owner", "alice").
		WithActions("list", "read", "write", "create", "delete").Create()
	require.NoError(t, err)

	// AND with following permissions
	rwlPerm, err := orgAdapter.Permissions(namespace).
		WithResource(project.Resource).
		WithEffect(types.Effect_PERMITTED).
		WithScope("Reporting").
		WithConstraints(`
	{{or (eq .Principal.Username .Resource.Owner) (Not .Private)}}
`).WithActions("read", "write", "list").Create()
	require.NoError(t, err)

	// assigning permission to roles
	require.NoError(t, alice.AddPermissions(rwlPerm.Permission))
	require.NoError(t, bob.AddPermissions(rwlPerm.Permission))

	// Project should be only be accessible by alice as the scope matches and she is the owner.
	require.NoError(t, alice.Authorizer(namespace).
		WithAction("list").
		WithScope("Reporting").
		WithContext("Private", "true").
		WithResourceName("nextgen-app").Check())

	// But alice should not be able to access without matching scope.
	require.Error(t, alice.Authorizer(namespace).
		WithAction("list").
		WithScope("").
		WithContext("Private", "true").
		WithResourceName("nextgen-app").Check())

	// But bob should not be able to access as project is private and he is not the owner.
	require.Error(t, bob.Authorizer(namespace).
		WithAction("list").
		WithScope("Reporting").
		WithContext("Private", "true").
		WithResourceName("nextgen-app").Check())

	// Project should be accessible by both if private is false
	require.NoError(t, alice.Authorizer(namespace).
		WithAction("list").
		WithScope("Reporting").
		WithContext("Private", "false").
		WithResourceName("nextgen-app").Check())
	require.NoError(t, bob.Authorizer(namespace).
		WithAction("list").
		WithScope("Reporting").
		WithContext("Private", "false").
		WithResourceName("nextgen-app").Check())
}

func testForDetectingAmbiguousPermission(
	t *testing.T,
	authAdapter *AuthAdapter,
) {
	// create org
	orgAdapter, err := authAdapter.CreateOrganization(
		&types.Organization{
			Name:       "xyz-corp",
			Namespaces: []string{"marketing", "sales"},
		})
	require.NoError(t, err)
	namespace := orgAdapter.Organization.Namespaces[0]

	// AND with following principals
	alice, err := orgAdapter.Principals().
		WithUsername("alice").
		WithAttributes("Title", "Engineer", "Permanent", "true").Create()

	// AND with following resources
	project1, err := orgAdapter.Resources(namespace).
		WithName("urn:org-sales-*-project-1000-*").
		WithAttributes("Year", fmt.Sprintf("%d", time.Now().Year())).
		WithActions("list", "read", "write", "create", "delete").Create()
	require.NoError(t, err)

	// AND with following permissions
	rwlPerm1, err := orgAdapter.Permissions(namespace).
		WithResource(project1.Resource).
		WithEffect(types.Effect_PERMITTED).
		WithConstraints(`
	{{$CurrentYear := TimeNow "2006"}}
	{{and (eq .Principal.Permanent "true") (eq .CurrentLocation "Chicago") (eq .Resource.Year $CurrentYear)}}
`).WithActions("read", "write", "list").Create()
	require.NoError(t, err)

	rwlPerm2, err := orgAdapter.Permissions(namespace).
		WithResource(project1.Resource).
		WithEffect(types.Effect_DENIED).
		WithConstraints(`
	{{and (eq .Principal.Permanent "true") (eq .CurrentLocation "Chicago")}}
`).WithActions("read", "write").Create()
	require.NoError(t, err)

	// assigning permission to roles
	require.NoError(t, alice.AddPermissions(rwlPerm1.Permission, rwlPerm2.Permission))

	// Test for list should be fine as it's only defined once
	require.NoError(t, alice.Authorizer(namespace).
		WithAction("list").
		WithResourceName("urn:org-sales-abc-project-1000-xyz").
		WithContext("CurrentLocation", "Chicago").Check())

	// But test for read permissions should fail with ambiguous error
	err = alice.Authorizer(namespace).
		WithAction("read").
		WithResourceName("urn:org-sales-abc-project-1000-xyz").
		WithContext("CurrentLocation", "Chicago").Check()
	require.Error(t, err)
	require.Contains(t, err.Error(), "EC100452")
}

func testPermissionsWithScopeAndWithinGeoLocation(
	t *testing.T,
	authAdapter *AuthAdapter,
) {
	// create org
	orgAdapter, err := authAdapter.CreateOrganization(
		&types.Organization{
			Name:       "xyz-corp",
			Namespaces: []string{"marketing", "sales"},
		})
	require.NoError(t, err)
	namespace := orgAdapter.Organization.Namespaces[0]

	// AND with following principals
	alice, err := orgAdapter.Principals().
		WithUsername("alice").
		WithAttributes("Title", "Engineer", "Permanent", "true").Create()
	require.NoError(t, err)

	bob, err := orgAdapter.Principals().
		WithUsername("bob").
		WithAttributes("Title", "Manager", "Permanent", "true").Create()
	require.NoError(t, err)

	// AND with following resources
	project1, err := orgAdapter.Resources(namespace).
		WithName("urn:org-sales-*-project-1000-*").
		WithAttributes("Year", fmt.Sprintf("%d", time.Now().Year())).
		WithActions("list", "read", "write", "create", "delete").Create()
	require.NoError(t, err)

	// AND with following permissions
	rwlPerm1, err := orgAdapter.Permissions(namespace).
		WithResource(project1.Resource).
		WithScope("svc1").
		WithConstraints(`
	{{$CurrentYear := TimeNow "2006"}}
	{{and (eq .Principal.Permanent "true") (DistanceWithinKM .UserLatLng "46.879967,-121.726906" 100) (eq .Resource.Year $CurrentYear)}}
`).WithActions("read", "write", "list").Create()
	require.NoError(t, err)

	rwlPerm2, err := orgAdapter.Permissions(namespace).
		WithResource(project1.Resource).
		WithScope("svc2").
		WithConstraints(`
	{{$CurrentYear := TimeNow "2006"}}
	{{and (eq .Principal.Permanent "true") (DistanceWithinKM .UserLatLng "46.879967,-121.726906" 100) 
(eq .Resource.Year $CurrentYear)}}
`).WithActions("read", "write", "list").Create()
	require.NoError(t, err)

	// assigning permission to roles
	require.NoError(t, alice.AddPermissions(rwlPerm1.Permission))
	require.NoError(t, bob.AddPermissions(rwlPerm2.Permission))

	// Test for permissions with scope
	// Alice should have permission for svc1 scope
	require.NoError(t, alice.Authorizer(namespace).
		WithAction("list").
		WithScope("svc1").
		WithResourceName("urn:org-sales-abc-project-1000-xyz").
		WithContext("UserLatLng", "47.620422,-122.349358").Check())

	// But not svc2
	require.Error(t, alice.Authorizer(namespace).
		WithAction("list").
		WithScope("svc2").
		WithResourceName("urn:org-sales-abc-project-1000-xyz").
		WithContext("UserLatLng", "47.620422,-122.349358").Check())

	// Bob should have permission for svc2
	require.NoError(t, bob.Authorizer(namespace).
		WithAction("list").
		WithScope("svc2").
		WithResourceName("urn:org-sales-abc-project-1000-xyz").
		WithContext("UserLatLng", "47.620422,-122.349358").Check())
}

func testPermissionsForProjects(
	t *testing.T,
	authAdapter *AuthAdapter,
) {
	// create org
	orgAdapter, err := authAdapter.CreateOrganization(
		&types.Organization{
			Name:       "xyz-corp",
			Namespaces: []string{"marketing", "sales"},
		})
	require.NoError(t, err)
	namespace := orgAdapter.Organization.Namespaces[0]

	// AND with following principals
	alice, err := orgAdapter.Principals().
		WithUsername("alice").
		WithAttributes("Title", "Engineer", "Permanent", "true").Create()
	require.NoError(t, err)

	bob, err := orgAdapter.Principals().
		WithUsername("bob").
		WithAttributes("Title", "Manager", "Permanent", "true").Create()
	require.NoError(t, err)

	charlie, err := orgAdapter.Principals().
		WithUsername("charlie").
		WithAttributes("Title", "PM", "Permanent", "true").Create()
	require.NoError(t, err)

	david, err := orgAdapter.Principals().
		WithUsername("david").
		WithAttributes("Title", "Contractor", "Permanent", "false").Create()
	require.NoError(t, err)

	// AND with following groups
	admin, err := orgAdapter.Groups(namespace).WithName("Admin").Create()
	require.NoError(t, err)

	// AND with following resources
	project1, err := orgAdapter.Resources(namespace).
		WithName("urn:org-sales-*-project-1000-*").
		WithAttributes("Year", fmt.Sprintf("%d", time.Now().Year())).
		WithActions("list", "read", "write", "create", "delete").Create()
	require.NoError(t, err)
	project2, err := orgAdapter.Resources(namespace).
		WithName("urn:org-eng-*-project-2000-*").
		WithAttributes("Year", fmt.Sprintf("%d", time.Now().Year())).
		WithActions("list", "read", "write", "create", "delete").Create()
	require.NoError(t, err)

	// AND with following permissions
	listPerm1, err := orgAdapter.Permissions(namespace).
		WithResource(project1.Resource).
		WithConstraints(`
	{{$CurrentYear := TimeNow "2006"}}
	{{and (eq .Principal.Permanent "true") (eq .CurrentLocation "Chicago") (eq .Resource.Year $CurrentYear)}}
`).WithActions("list").Create()
	require.NoError(t, err)

	rwPerm1, err := orgAdapter.Permissions(namespace).
		WithResource(project1.Resource).
		WithConstraints(`
	{{$CurrentYear := TimeNow "2006"}}
	{{and (eq .Principal.Permanent "true") (eq .CurrentLocation "Chicago") (eq .Resource.Year $CurrentYear)}}
`).WithActions("read", "write").Create()
	require.NoError(t, err)

	cdPerm1, err := orgAdapter.Permissions(namespace).
		WithResource(project1.Resource).
		WithConstraints(`
{{$CurrentYear := TimeNow "2006"}}
{{and (eq .Principal.Permanent "true") (eq .CurrentLocation "Chicago") (HasGroup "Admin") (eq .Resource.Year $CurrentYear)}} 
`).WithActions("create", "delete").Create()
	require.NoError(t, err)

	listPerm2, err := orgAdapter.Permissions(namespace).
		WithResource(project2.Resource).
		WithConstraints(`
	{{$CurrentYear := TimeNow "2006"}}
	{{and (eq .Principal.Permanent "true") (eq .CurrentLocation "Chicago") (eq .Resource.Year $CurrentYear)}}
`).WithActions("list").Create()
	require.NoError(t, err)

	rwPerm2, err := orgAdapter.Permissions(namespace).
		WithResource(project2.Resource).
		WithConstraints(`
	{{$CurrentYear := TimeNow "2006"}}
	{{and (eq .Principal.Permanent "true") (eq .CurrentLocation "Chicago") (eq .Resource.Year $CurrentYear)}}
`).WithActions("read", "write").Create()
	require.NoError(t, err)

	cdPerm2, err := orgAdapter.Permissions(namespace).
		WithResource(project2.Resource).
		WithConstraints(`
{{$CurrentYear := TimeNow "2006"}}
{{and (eq .Principal.Permanent "true") (eq .CurrentLocation "Chicago") (HasGroup "Admin") (eq .Resource.Year $CurrentYear)}} 
`).WithActions("create", "delete").Create()
	require.NoError(t, err)

	// assigning permission to roles
	require.NoError(t, alice.AddPermissions(listPerm1.Permission, listPerm2.Permission, rwPerm1.Permission))
	require.NoError(t, bob.AddPermissions(listPerm2.Permission, rwPerm2.Permission, cdPerm2.Permission))
	require.NoError(t, charlie.AddPermissions(listPerm1.Permission, listPerm2.Permission, rwPerm1.Permission, rwPerm2.Permission, cdPerm1.Permission, cdPerm2.Permission))
	require.NoError(t, charlie.AddGroups(admin.Group))
	// No permissions for david

	// Test for list permissions
	// Alice should have permission to list both projects and read/write project1
	require.NoError(t, alice.Authorizer(namespace).
		WithAction("list").
		WithResourceName("urn:org-sales-abc-project-1000-xyz").
		WithContext("CurrentLocation", "Chicago").Check())

	require.NoError(t, alice.Authorizer(namespace).
		WithAction("list").
		WithResourceName("urn:org-eng-abc-project-2000-ijk").
		WithContext("CurrentLocation", "Chicago").Check())

	require.NoError(t, alice.Authorizer(namespace).
		WithAction("read").
		WithResourceName("urn:org-sales-abc-project-1000-xyz").
		WithContext("CurrentLocation", "Chicago").Check())

	require.NoError(t, alice.Authorizer(namespace).
		WithAction("write").
		WithResourceName("urn:org-sales-abc-project-1000-xyz").
		WithContext("CurrentLocation", "Chicago").Check())

	// but not project2
	require.Error(t, alice.Authorizer(namespace).
		WithAction("read").
		WithResourceName("urn:org-eng-abc-project-2000-ijk").
		WithContext("CurrentLocation", "Chicago").Check())

	// Bob should have permission for only project2
	require.Error(t, bob.Authorizer(namespace).
		WithAction("list").
		WithResourceName("urn:org-sales-abc-project-1000-xyz").
		WithContext("CurrentLocation", "Chicago").Check())

	require.NoError(t, bob.Authorizer(namespace).
		WithAction("list").
		WithResourceName("urn:org-eng-abc-project-2000-ijk").
		WithContext("CurrentLocation", "Chicago").Check())

	require.NoError(t, bob.Authorizer(namespace).
		WithAction("read").
		WithResourceName("urn:org-eng-abc-project-2000-ijk").
		WithContext("CurrentLocation", "Chicago").Check())

	require.NoError(t, bob.Authorizer(namespace).
		WithAction("write").
		WithResourceName("urn:org-eng-abc-project-2000-ijk").
		WithContext("CurrentLocation", "Chicago").Check())

	// but not create / delete
	require.Error(t, bob.Authorizer(namespace).
		WithAction("create").
		WithResourceName("urn:org-eng-abc-project-2000-ijk").
		WithContext("CurrentLocation", "Chicago").Check())

	// Charlie should have permission for both project1 and project2 for all operations
	require.NoError(t, charlie.Authorizer(namespace).
		WithAction("list").
		WithResourceName("urn:org-sales-abc-project-1000-xyz").
		WithContext("CurrentLocation", "Chicago").Check())

	require.NoError(t, charlie.Authorizer(namespace).
		WithAction("read").
		WithResourceName("urn:org-sales-abc-project-1000-xyz").
		WithContext("CurrentLocation", "Chicago").Check())

	require.NoError(t, charlie.Authorizer(namespace).
		WithAction("write").
		WithResourceName("urn:org-sales-abc-project-1000-xyz").
		WithContext("CurrentLocation", "Chicago").Check())

	require.NoError(t, charlie.Authorizer(namespace).
		WithAction("create").
		WithResourceName("urn:org-sales-abc-project-1000-xyz").
		WithContext("CurrentLocation", "Chicago").Check())

	require.NoError(t, charlie.Authorizer(namespace).
		WithAction("delete").
		WithResourceName("urn:org-sales-abc-project-1000-xyz").
		WithContext("CurrentLocation", "Chicago").Check())

	require.NoError(t, charlie.Authorizer(namespace).
		WithAction("list").
		WithResourceName("urn:org-eng-abc-project-2000-ijk").
		WithContext("CurrentLocation", "Chicago").Check())

	require.NoError(t, charlie.Authorizer(namespace).
		WithAction("read").
		WithResourceName("urn:org-eng-abc-project-2000-ijk").
		WithContext("CurrentLocation", "Chicago").Check())

	require.NoError(t, charlie.Authorizer(namespace).
		WithAction("write").
		WithResourceName("urn:org-eng-abc-project-2000-ijk").
		WithContext("CurrentLocation", "Chicago").Check())

	require.NoError(t, charlie.Authorizer(namespace).
		WithAction("create").
		WithResourceName("urn:org-eng-abc-project-2000-ijk").
		WithContext("CurrentLocation", "Chicago").Check())

	require.NoError(t, charlie.Authorizer(namespace).
		WithAction("delete").
		WithResourceName("urn:org-eng-abc-project-2000-ijk").
		WithContext("CurrentLocation", "Chicago").Check())

	// But not outside Chicago
	require.Error(t, charlie.Authorizer(namespace).
		WithAction("list").
		WithResourceName("urn:org-sales-abc-project-1000-xyz").
		WithContext("CurrentLocation", "Seattle").Check())
	require.Error(t, charlie.Authorizer(namespace).
		WithAction("read").
		WithResourceName("urn:org-sales-abc-project-1000-xyz").
		WithContext("CurrentLocation", "Seattle").Check())

	require.Error(t, charlie.Authorizer(namespace).
		WithAction("write").
		WithResourceName("urn:org-sales-abc-project-1000-xyz").
		WithContext("CurrentLocation", "Seattle").Check())

	require.Error(t, charlie.Authorizer(namespace).
		WithAction("create").
		WithResourceName("urn:org-sales-abc-project-1000-xyz").
		WithContext("CurrentLocation", "Seattle").Check())

	require.Error(t, charlie.Authorizer(namespace).
		WithAction("delete").
		WithResourceName("urn:org-sales-abc-project-1000-xyz").
		WithContext("CurrentLocation", "Seattle").Check())

	// David should have no permissions
	require.Error(t, david.Authorizer(namespace).
		WithAction("list").
		WithResourceName("urn:org-sales-abc-project-1000-xyz").
		WithContext("CurrentLocation", "Seattle").Check())

	require.Error(t, david.Authorizer(namespace).
		WithAction("read").
		WithResourceName("urn:org-sales-abc-project-1000-xyz").
		WithContext("CurrentLocation", "Seattle").Check())

	require.Error(t, david.Authorizer(namespace).
		WithAction("write").
		WithResourceName("urn:org-sales-abc-project-1000-xyz").
		WithContext("CurrentLocation", "Seattle").Check())

	require.Error(t, david.Authorizer(namespace).
		WithAction("create").
		WithResourceName("urn:org-sales-abc-project-1000-xyz").
		WithContext("CurrentLocation", "Seattle").Check())

	require.Error(t, david.Authorizer(namespace).
		WithAction("delete").
		WithResourceName("urn:org-sales-abc-project-1000-xyz").
		WithContext("CurrentLocation", "Seattle").Check())

	require.Error(t, david.Authorizer(namespace).
		WithAction("list").
		WithResourceName("urn:org-eng-abc-project-2000-ijk").
		WithContext("CurrentLocation", "Seattle").Check())

	require.Error(t, david.Authorizer(namespace).
		WithAction("read").
		WithResourceName("urn:org-eng-abc-project-2000-ijk").
		WithContext("CurrentLocation", "Seattle").Check())

	require.Error(t, david.Authorizer(namespace).
		WithAction("write").
		WithResourceName("urn:org-eng-abc-project-2000-ijk").
		WithContext("CurrentLocation", "Seattle").Check())

	require.Error(t, david.Authorizer(namespace).
		WithAction("create").
		WithResourceName("urn:org-eng-abc-project-2000-ijk").
		WithContext("CurrentLocation", "Seattle").Check())

	require.Error(t, david.Authorizer(namespace).
		WithAction("delete").
		WithResourceName("urn:org-eng-abc-project-2000-ijk").
		WithContext("CurrentLocation", "Seattle").Check())
}

func testPermissionsForDepositAccount(
	t *testing.T,
	authAdapter *AuthAdapter,
) {
	// create org
	orgAdapter, err := authAdapter.CreateOrganization(
		&types.Organization{
			Name:       "abc-bank",
			Namespaces: []string{"finance", "loan", "checking"},
		})
	require.NoError(t, err)
	namespace := orgAdapter.Organization.Namespaces[0]

	// AND with following principals
	tom, err := orgAdapter.Principals().
		WithUsername("tom").
		WithAttributes("Region", "Midwest").Create()
	require.NoError(t, err)

	cassy, err := orgAdapter.Principals().
		WithUsername("cassy").
		WithAttributes("Region", "Midwest").Create()
	require.NoError(t, err)

	ali, err := orgAdapter.Principals().
		WithUsername("ali").
		WithAttributes("Region", "Midwest").Create()
	require.NoError(t, err)

	mike, err := orgAdapter.Principals().
		WithUsername("mike").
		WithAttributes("Region", "Midwest").Create()
	require.NoError(t, err)

	larry, err := orgAdapter.Principals().
		WithUsername("larry").
		WithAttributes("Region", "Midwest").Create()
	require.NoError(t, err)

	depositAccount, err := orgAdapter.Resources(namespace).WithName("DepositAccount").
		WithAttributes("AccountType", "Checking").
		WithActions("balance", "withdraw", "deposit", "open", "close").Create()
	require.NoError(t, err)
	loanAccount, err := orgAdapter.Resources(namespace).WithName("LoanAccount").
		WithAttributes("Rate", "4.5", "MaxLoanBalance", "5000").
		WithActions("create", "delete", "read", "write").Create()
	require.NoError(t, err)
	generalLedger, err := orgAdapter.Resources(namespace).WithName("GeneralLedger").
		WithAttributes("LedgerYear", fmt.Sprintf("%d", time.Now().Year())).
		WithActions("create", "delete", "read", "write").Create()
	require.NoError(t, err)
	postingRules, err := orgAdapter.Resources(namespace).WithName("GeneralLedgerPostingRules").
		WithAttributes("PostingYear", fmt.Sprintf("%d", time.Now().Year())).
		WithActions("read", "post").Create()
	require.NoError(t, err)

	// AND with roles
	employee, err := orgAdapter.Roles(namespace).WithName("Employee").Create()
	require.NoError(t, err)
	teller, err := orgAdapter.Roles(namespace).WithName("Teller").
		WithParents(employee.Role).Create()
	require.NoError(t, err)
	csr, err := orgAdapter.Roles(namespace).WithName("CSR").
		WithParents(teller.Role).Create()
	require.NoError(t, err)
	accountant, err := orgAdapter.Roles(namespace).WithName("Accountant").
		WithParents(employee.Role).Create()
	require.NoError(t, err)
	accountantMgr, err := orgAdapter.Roles(namespace).WithName("AccountingManager").
		WithParents(accountant.Role).Create()
	require.NoError(t, err)
	loanOfficer, err := orgAdapter.Roles(namespace).WithName("LoanOfficer").
		WithParents(accountantMgr.Role).Create()
	require.NoError(t, err)

	// AND with following permissions
	listPerm, err := orgAdapter.Permissions(namespace).
		WithResource(depositAccount.Resource).
		WithActions("balance").Create()
	require.NoError(t, err)
	depositPerm, err := orgAdapter.Permissions(namespace).
		WithResource(depositAccount.Resource).
		WithConstraints(`and (eq .Principal.Region "Midwest") (eq .CurrentLocation "Chicago")`).
		WithActions("deposit", "withdraw").Create()
	require.NoError(t, err)
	openClosePerm, err := orgAdapter.Permissions(namespace).
		WithResource(depositAccount.Resource).
		WithConstraints(`and (eq .Principal.Region "Midwest") (eq .CurrentLocation "Chicago")`).
		WithActions("open", "close").Create()
	require.NoError(t, err)

	cdLoanPerm, err := orgAdapter.Permissions(namespace).
		WithResource(loanAccount.Resource).
		WithConstraints(`and (eq .Principal.Region "Midwest") (eq .CurrentLocation "Chicago") 
(lt .CurrentBalance .Resource.MaxLoanBalance)`).
		WithActions("create", "delete").Create()
	require.NoError(t, err)
	rwLoanPerm, err := orgAdapter.Permissions(namespace).
		WithResource(loanAccount.Resource).
		WithConstraints(`and (eq .Principal.Region "Midwest") (eq .CurrentLocation "Chicago") 
(lt .CurrentBalance .Resource.MaxLoanBalance)`).
		WithActions("read", "write").Create()
	require.NoError(t, err)

	cdLedgerPerm, err := orgAdapter.Permissions(namespace).
		WithResource(generalLedger.Resource).
		WithConstraints(`
{{$CurrentYear := TimeNow "2006"}}
{{and (eq .Principal.Region "Midwest") (eq .CurrentLocation "Chicago") (eq .Resource.LedgerYear $CurrentYear) 
(ActionIncludes "create" "delete") (HasRole "AccountingManager")}}
`).
		WithActions("create", "delete").Create()
	require.NoError(t, err)
	rwLedgerPerm, err := orgAdapter.Permissions(namespace).
		WithResource(generalLedger.Resource).
		WithConstraints(
			`
	{{$CurrentYear := TimeNow "2006"}}
	{{and (eq .Principal.Region "Midwest") (eq .CurrentLocation "Chicago") (eq .Resource.LedgerYear $CurrentYear)}}
`).
		WithActions("read", "write").Create()
	require.NoError(t, err)
	allGlprPerm, err := orgAdapter.Permissions(namespace).
		WithResource(postingRules.Resource).
		WithConstraints(
			`
	{{$CurrentYear := TimeNow "2006"}}
	{{and (eq .Principal.Region "Midwest") (eq .CurrentLocation "Chicago") (eq .Resource.PostingYear $CurrentYear)
(ActionIncludes "post" "read") (HasRole "LoanOfficer")}}
`).
		WithActions("post", "read").Create()
	require.NoError(t, err)

	// assigning permission to roles
	require.NoError(t, employee.AddPermissions(listPerm.Permission))
	require.NoError(t, teller.AddPermissions(depositPerm.Permission))
	require.NoError(t, csr.AddPermissions(openClosePerm.Permission))
	require.NoError(t, accountant.AddPermissions(rwLoanPerm.Permission))
	require.NoError(t, accountant.AddPermissions(rwLedgerPerm.Permission))
	require.NoError(t, accountantMgr.AddPermissions(cdLoanPerm.Permission))
	require.NoError(t, accountantMgr.AddPermissions(cdLedgerPerm.Permission))
	require.NoError(t, loanOfficer.AddPermissions(allGlprPerm.Permission))

	// WHEN assigning creating roles
	require.NoError(t, tom.AddRoles(teller.Role))
	require.NoError(t, cassy.AddRoles(csr.Role))
	require.NoError(t, ali.AddRoles(accountant.Role))
	require.NoError(t, mike.AddRoles(accountantMgr.Role))
	require.NoError(t, larry.AddRoles(loanOfficer.Role))

	// Test for DepositAccount
	// tom with teller should access balance of DepositAccount
	require.NoError(t, tom.Authorizer(namespace).
		WithAction("balance").
		WithResource(depositAccount.Resource).Check())

	// tom, the teller should not be able to access deposit without CurrentLocation == Chicago
	require.Error(t, tom.Authorizer(namespace).
		WithAction("deposit").
		WithResource(depositAccount.Resource).Check())

	// tom, the teller should fail with CurrentLocation == Seattle
	require.Error(t, tom.Authorizer(namespace).
		WithAction("deposit").
		WithResource(depositAccount.Resource).
		WithContext("CurrentLocation", "Seattle").Check())
	// tom, the teller should succeed now with CurrentLocation == Chicago
	require.NoError(t, tom.Authorizer(namespace).
		WithAction("deposit").
		WithResource(depositAccount.Resource).
		WithContext("CurrentLocation", "Chicago").Check())
	// tom, the teller should not succeed with open action
	require.Error(t, tom.Authorizer(namespace).
		WithAction("open").
		WithResource(depositAccount.Resource).
		WithContext("CurrentLocation", "Chicago").Check())
	// cassy, the csr should succeed with open action
	require.NoError(t, cassy.Authorizer(namespace).
		WithAction("open").
		WithResource(depositAccount.Resource).
		WithContext("CurrentLocation", "Chicago").Check())

	// ali, the accountant should not be able to access deposit because accountant role does not extend from teller
	require.Error(t, ali.Authorizer(namespace).
		WithAction("deposit").
		WithResource(depositAccount.Resource).
		WithContext("CurrentLocation", "Chicago").Check())

	// mike, the account manager should not be able to access withdraw because accountantMgr role does not extend from teller
	require.Error(t, mike.Authorizer(namespace).
		WithAction("withdraw").
		WithResource(depositAccount.Resource).
		WithContext("CurrentLocation", "Chicago").Check())

	// Test for LoanOfficer
	// ali, the accountant should be able to read LoanAccount
	require.NoError(t, ali.Authorizer(namespace).
		WithAction("read").
		WithResource(loanAccount.Resource).
		WithContext("CurrentLocation", "Chicago", "CurrentBalance", "4000").Check())

	// ali, the accountant should not be able to delete LoanAccount
	require.Error(t, ali.Authorizer(namespace).
		WithAction("delete").
		WithResource(loanAccount.Resource).
		WithContext("CurrentLocation", "Chicago", "CurrentBalance", "4000").Check())

	// mike, the accountant-manager should be able to delete LoanAccount
	require.NoError(t, mike.Authorizer(namespace).
		WithAction("delete").
		WithResource(loanAccount.Resource).
		WithContext("CurrentLocation", "Chicago", "CurrentBalance", "4000").Check())

	// Test for GeneralLedger
	// ali, the accountant should be able to read GeneralLedger
	require.NoError(t, ali.Authorizer(namespace).
		WithAction("read").
		WithResource(generalLedger.Resource).
		WithContext("CurrentLocation", "Chicago").Check())

	// ali, the accountant should not be able to delete GeneralLedger
	require.Error(t, ali.Authorizer(namespace).
		WithAction("delete").
		WithResource(generalLedger.Resource).
		WithContext("CurrentLocation", "Chicago").Check())

	// mike, the account-manager should be able to delete GeneralLedger
	require.NoError(t, mike.Authorizer(namespace).
		WithAction("delete").
		WithResource(generalLedger.Resource).
		WithContext("CurrentLocation", "Chicago").Check())

	// Test GeneralLedgerPostingRules
	// mike, the account-manager should not be able to read GeneralLedgerPostingRules
	require.Error(t, mike.Authorizer(namespace).
		WithAction("read").
		WithResource(postingRules.Resource).
		WithContext("CurrentLocation", "Chicago").Check())

	// larry, the loan officer should be able to read and post GeneralLedgerPostingRules
	require.NoError(t, larry.Authorizer(namespace).
		WithAction("read").
		WithResource(postingRules.Resource).
		WithContext("CurrentLocation", "Chicago").Check())
	require.NoError(t, larry.Authorizer(namespace).
		WithAction("post").
		WithResource(postingRules.Resource).
		WithContext("CurrentLocation", "Chicago").Check())
}

func testCRUD(
	t *testing.T,
	authAdapter *AuthAdapter,
) {
	// create org
	orgAdapter, err := authAdapter.CreateOrganization(
		&types.Organization{
			Name:       "bank-abc",
			Namespaces: []string{"Checking", "Loan"},
		})
	require.NoError(t, err)
	namespace := orgAdapter.Organization.Namespaces[0]

	_, err = authAdapter.GetOrganization(orgAdapter.Organization.Id)
	require.NoError(t, err)

	require.NoError(t, orgAdapter.Update())

	// AND with following principals
	alice, err := orgAdapter.Principals().
		WithUsername("alice").
		WithAttributes("Department", "Personal Banking", "EmploymentLength", "5").Create()
	require.NoError(t, err)
	require.NoError(t, alice.Update())
	require.NoError(t, alice.Get(alice.Principal.Id))

	resource, err := orgAdapter.Resources(namespace).WithName("buffet").
		WithActions("all").Create()
	require.NoError(t, err)
	require.NoError(t, resource.Update())
	require.NoError(t, resource.Get(resource.Resource.Id))

	rel, err := alice.Relationships(namespace).
		WithResource(resource.Resource).WithRelation("relation").Create()
	require.NoError(t, err)
	require.NoError(t, rel.Update())
	require.NoError(t, rel.Get(rel.Relationship.Id))

	perm, err := orgAdapter.Permissions(namespace).
		WithResource(resource.Resource).WithActions("all").Create()
	require.NoError(t, err)
	require.NoError(t, perm.Update())
	require.NoError(t, perm.Get(perm.Permission.Id))

	role, err := orgAdapter.Roles(namespace).WithName("role").Create()
	require.NoError(t, err)
	crole, err := orgAdapter.Roles(namespace).WithName("crole").WithParents(role.Role).Create()
	require.NoError(t, err)
	require.NoError(t, crole.Update())
	require.NoError(t, crole.Get(crole.Role.Id))
	require.NoError(t, alice.AddRoles(crole.Role))
	require.NoError(t, alice.DeleteRoles(crole.Role))

	require.NoError(t, alice.AddPermissions(perm.Permission))
	require.NoError(t, alice.DeletePermissions(perm.Permission))

	group, err := orgAdapter.Groups(namespace).WithName("group").Create()
	require.NoError(t, err)
	cgroup, err := orgAdapter.Groups(namespace).WithName("cgroup").WithParents(group.Group).Create()
	require.NoError(t, err)
	require.NoError(t, alice.AddGroups(cgroup.Group))
	require.NoError(t, alice.DeleteGroups(cgroup.Group))
	require.NoError(t, cgroup.Update())
	require.NoError(t, cgroup.Get(cgroup.Group.Id))
	require.NoError(t, cgroup.AddRoles(crole.Role))
	require.NoError(t, cgroup.DeleteRoles(crole.Role))

	require.NoError(t, alice.AddRelations(rel.Relationship))
	require.NoError(t, alice.DeleteRelations(rel.Relationship))

	require.NoError(t, alice.Delete())
	require.NoError(t, orgAdapter.Delete())
}

func runTests(
	t *testing.T,
	registry *metrics.Registry,
	fns ...func(t *testing.T, authAdapter *AuthAdapter),
) {
	_ = os.Setenv("CONFIG_DIR", "../../config")
	cfg, err := domain.NewConfig("")
	require.NoError(t, err)
	cfg.GrpcSasl = true

	_, webTeardown := controller.SetupWebServerForTesting(t, cfg, nil)
	_, grpcTeardown := server.SetupGrpcServerForTesting(t, cfg, domain.RootClientType, nil)

	// Create auth-service based on redis database
	cfg.PersistenceProvider = domain.RedisPersistenceProvider
	cfg.AuthServiceProvider = domain.DatabaseAuthServiceProvider
	redisAuthService, redisCC, err := factory.CreateAuthAdminService(cfg, registry, domain.RootClientType, "")
	require.NoError(t, err)

	// Create auth-service based on Dynamo DB database
	cfg.PersistenceProvider = domain.DynamoDBPersistenceProvider
	ddbAuthService, ddbCC, err := factory.CreateAuthAdminService(cfg, registry, domain.RootClientType, "")
	require.NoError(t, err)

	// Create auth-service based on gRPC client -- assuming the grpc server is running (which will be started in setup methods)
	cfg.AuthServiceProvider = domain.GrpcAuthServiceProvider
	grpcAuthService, grpcCC, err := factory.CreateAuthAdminService(cfg, registry, domain.RootClientType, cfg.GrpcListenPort)
	require.NoError(t, err)

	// Create auth-service based on http client -- assuming the web server is running (which will be started in setup methods)
	cfg.AuthServiceProvider = domain.HttpAuthServiceProvider
	httpAuthService, httpCC, err := factory.CreateAuthAdminService(cfg, registry, domain.RootClientType, "http://"+cfg.HttpListenPort)
	require.NoError(t, err)

	authServices := []service.AuthAdminService{redisAuthService, ddbAuthService, grpcAuthService, httpAuthService}

	// Go through all auth-service implementation to test each function
	for _, authSvc := range authServices {
		authorizer, err := authz.CreateAuthorizer(authz.DefaultAuthorizerKind, cfg, authSvc)
		require.NoError(t, err)

		for _, fn := range fns {
			authAdapter := New(authorizer, authSvc)
			fn(t, authAdapter)
		}
	}

	_ = redisCC.Close()
	_ = ddbCC.Close()
	_ = grpcCC.Close()
	_ = httpCC.Close()
	webTeardown()
	grpcTeardown()
}
