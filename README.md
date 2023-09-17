# PlexAuthZ
Hybrid Authentication API based on RBAC, ABAC, PBAC, and ReBAC.

## Installation
- See https://grpc.io/docs/protoc-installation/
- go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
- go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
- go get github.com/cloudflare/cfssl/cmd/cfssl
- go get github.com/cloudflare/cfssl/cmd/cfssljson
- 
## API Docs
- https://petstore.swagger.io/?url=https://raw.githubusercontent.com/bhatti/PlexAuthZ/main/docs/swagger.yaml


## Key Concepts

### What is an Access Control System
An access control system defines a framework for controlling the accessibility of resources in an organization or system. 
It restricts unauthorized entities from performing actions or accessing data they shouldn’t. 
Access control systems can be physical—controlling who can enter a building, for example—or digital—controlling who can 
access a computer system or network. The access control system consists of following components:

- **Authentication**: It is the process of verifying the identity of a user, application, or system. This process ensures that the entity requesting access is who or what it claims to be. Common methods of authentication include username and password, multi-factor authentication, biometric verification, token-based authentication, and certificate-based authentication.
- **Authorization**: It determines the level of access, or permissions, granted to a legitimately authenticated user or system. Essentially, it answers the question: “What is this authenticated entity allowed to do or see within the system?”.
- **Audit and Monitoring**: It refer to the systematic tracking, recording, and analysis of activities or events within a system or network. These activities often include user actions, system accesses, and operations that affect data and resources. The primary goals are to ensure compliance with established policies, detect unauthorized or abnormal activities, and facilitate the identification of vulnerabilities or weaknesses. Elements often involved in audit and monitoring can include log files, real-time monitoring, alerts and notification, data analytics, and compliance reporting.
- **Policy Management**: It involves the creation, maintenance, and enforcement of rules, guidelines, and standard operating procedures that govern the behavior of users and systems within an organization or environment. These policies may include access policies, security policies, operational policies, compliance policies, change management policies, and policy auditing.

Following are popular mechanisms for enforcing authorization:

- **Role-Based Access Control (RBAC)**: In RBAC, permissions are associated with roles, and users are assigned to these roles. For example, a “Manager” role might have the ability to add or remove employees from a system, while an “Employee” role might only be able to view information. When a user is assigned a role, they inherit all the permissions that come with it.
- **Attribute-Based Access Control (ABAC)**: ABAC is a more flexible and complex system that uses attributes as building blocks in a rule-based approach to control access. These attributes can be associated with the user (e.g., age, department, job role), action (e.g., read, write), resource (e.g., file type, location), or even environmental factors (e.g., time of day, network security level). Policies are then crafted to allow or deny actions based on these attributes.
- **Policy-Based Access Control (PBAC)**: PBAC is similar to ABAC but tends to be more dynamic, incorporating real-time information into its decision-making process. For example, a PBAC system might evaluate current network threat levels or the outcome of a risk assessment to determine whether access should be granted or denied. Policies can be complex, allowing for a high degree of flexibility and context-aware decisions.
- **Access Control Lists (ACLs)**: A list specifying what actions a user or system can or cannot perform.
- **Capabilities**: In a capability-based security model, permissions are attached to tokens (capabilities) rather than to subjects (e.g., users) or objects (e.g., files). These tokens can be passed around between users and systems. Having a token allows a user to access a resource or perform an action. This model decentralizes the control of access, making it flexible but also potentially harder to manage at scale.
- **Permissions**: This is a simple and straightforward model where each object (like a file or database record) has associated permissions that specify which users can perform which types of operations (read, write, delete, etc.). This is often seen in file systems where each file and directory has an associated set of permission flags.
- **Discretionary Access Control (DAC)**: In DAC models, the owner of the resource has the discretion to set its permissions. For example, in many operating systems, the creator of a file can decide who can read or write to that file.
- **Mandatory Access Control (MAC)**: Unlike DAC, where users have some discretion over permissions, in MAC, the system enforces policies that users cannot alter. These policies often use labels or classifications (e.g., Top Secret, Confidential) to determine who can access what.

The approaches to authorization are not mutually exclusive and can be integrated to form hybrid systems. 
For instance, an enterprise might rely on RBAC for broad-based access management, while also employing ABAC or PBAC to 
handle more nuanced or sensitive use-cases. 

Following is a list of high-level data model concepts that are typically used in the authorization systems:

#### Principal

The entity (which could be a user, system, or another service) that is making the request. Principals are often authenticated before they are authorized to perform an action.

#### Subject

Similar to a principal, the subject refers to the entity that is attempting to access a particular resource. In some contexts, a subject may represent a real person that has multiple principal identities for various systems.

#### Permission

An action that a principal is allowed to perform on a particular resource. For example, reading a file, updating a database record, or deleting an account.

#### Claim

A statement made by the principal, usually after authentication, that annotates the principal with specific attributes (e.g., username, roles, permissions). Claims are often used in token-based authentication systems like JWT to carry information about the principal.

#### Role

A named collection of permissions that can be assigned to a principal. Roles simplify the management of permissions by grouping them together under a single label.

#### Group

A collection of principals that are treated as a single unit for the purpose of granting permissions. For example, an “Admins” group might be given a role that allows them to perform administrative actions.

#### Access Policy

A set of rules that define the conditions under which a particular action is allowed or denied. Policies can be simple (“Admins can do anything”) or complex (“Users can edit a document only if they are the creator and the document is in ‘Draft’ status”).

#### Relation

Relations define how different entities are connected. For instance, a “user” can have a “memberOf” relation with a “group”, or a “document” can have an “ownedBy” relation with a “user”.

#### Resource

The object that the principal wants to access (e.g., a file, a database record).

#### Context

Additional situational information (e.g., IP address, time of day) that might influence the authorization decision.

#### Namespace

Each namespace could serve as a container for a set of resources, roles, and permissions.

#### Scope or Realm

This often refers to the level or context in which a permission is granted. For instance, in OAuth, scopes are used to specify what access a token grants the user, like “read-only” access to a particular resource.

#### Rule

A specific condition or criterion in a policy that dictates whether access should be granted or denied.

#### Dynamic Conditions

Dynamic conditions or predicates are expressions that must be evaluated at runtime to determine if access should be granted or denied. Dynamic conditions consists of attributes, operators and values, e.g.,
```
if (principal.role == "employee" AND principal.status == "active") OR (time < "17:00") then ALLOW

if (principal.role == "admin") OR (document.owner == principal.id) then ALLOW

if IP_address in [allowed_IPs] then ALLOW

if time >= 09:00 AND time <= 17:00 then ALLOW

```

## Design Tenets of PlexAuthZ

PlexAuthZ implements a hybrid Authorization system combining Role-Based Access Control (RBAC), 
Attribute-Based Access Control (ABAC), and Relationship-Based Access Control (ReBAC). Following are primary design tenets for PlexAuthZ:

- **Scalability:** Capable of handling a large number of authorization requests per second and expand to accommodate growing numbers of users and resources.
- **Flexibility:** Supports RBAC, ABAC, and ReBAC, allowing for the handling of various scenarios.
- Fine-grained Control: Context aware such as time, location and real-time data, and can decide based on multiple attributes of the subject, object, and environment.
- **Auditing and Monitoring:** Detailed logs for all access attempts and policy changes, and real-time insights into access patterns, possibly with alerting for suspicious activities.
- **Security:** Applies least privilege, and enforces data masking and redaction.
- **Usability:** Easy-to-use interfaces for assigning, changing, and revoking roles.
- **Extensibility:** Comprehensive APIs for integration with other systems and services ability to run custom code during the authorization process.
- **Reliability:** have minimal downtime with backup and recovery.
- **Compliance:** adhere to regulatory requirements like GDPR, HIPAA, etc. and track changes to policies for auditing purposes.
- **Multi-Tenancy**: support multiple services and products with a variety of authorization models under a single unified system.
- **Policy Versioning and Namespacing**: allow multiple versions and namespaces of policies, making it possible to manage complex policy changes.
- **Balance between Expressive and Performance:** provide a good balance with expressive policies offered by OPA and high performance offered by Zanzibar.
- **Policy Validation:** check against invalid, unsafe or ambiguous policies and prevent users from making accidental mistakes.
- **Performance Optimization:** using cache, indexing, parallel processing, lazy evaluation, rule simplifications, automated reasoning, decision trees and other optimization techniques to improve performance of the system.


## PlexAuthZ Data Model
Following data model is defined in [Protocol Buffers](https://protobuf.dev/overview/) definition language based on 
above authorization concepts:

![Data Model](https://weblog.plexobject.com/images/plexauthz-data.png)

### Organization

The Organization abstracts a boundary of authorization data and it can have multiple namespaces for different security realms 
or segments of security domains. Here is the definition of Organization:

```protobuf3
message Organization {
  // ID unique identifier assigned to this organization.
  string id = 1;
  // Version
  int64 version = 2;
  // Name of organization.
  string name = 3;
  // Allowed Namespaces for organization.
  repeated string namespaces = 4;
  // url for organization.
  string url = 5;
  // Optional parent ids.
  repeated string parent_ids = 6;
}
```

### Principal

The Principal abstracts subject who is making an authorization request to perform an action on a target resource 
based on access rules and dynamic conditions. A Principal belongs to an organization and can be associated with groups, 
roles (RBAC), permissions and relationships (ReBAC). The Principal defines following properties:

```protobuf3
message Principal {
    // ID unique identifier assigned to this principal.
    string id = 1;
    // Version
    int64 version = 2;
    // OrganizationId of the principal user.
    string organization_id = 3;
    // Allowed Namespaces for principal, should be subset of namespaces in organization.
    repeated string namespaces = 4;
    // Username of the principal user.
    string username = 5;
    // Email of the principal user.
    string email = 6;
    // Name of the principal user.
    string name = 7;
    // Attributes of principal
    map<string, string> attributes = 8;
    // Groups that the principal belongs to.
    repeated string group_ids = 9;
    // Roles that the principal belongs to.
    repeated string role_ids = 10;
    // Permissions that the principal belongs to.
    repeated string permission_ids = 11;
    // Relationships that the principal belongs to.
    repeated string relation_ids = 12;
}
```

### Resource and ResourceInstance

The Resource represents target object for performing an action and checking an access rules policy. A resource can also 
be used to represent an object with a quota that can be allocated or assigned based on access policies. Here is a definition 
of Resource and ResourceInstance:

```protobuf3
message Resource {
    // ID unique identifier assigned to this resource.
    string id = 1;
    // Version
    int64 version = 2;
    // Namespace for resource.
    string namespace = 3;
    // Name of the resource.
    string name = 4;
    // capacity of resource.
    int32 capacity = 5;
    // Attributes of resource.
    map<string, string> attributes = 6;
    // AllowedActions that can be performed.
    repeated string allowed_actions = 7;
}

enum ResourceState {
    ALLOCATED = 0;
    AVAILABLE = 1;
}

message ResourceInstance {
    // ID unique identifier assigned to this resource instance.
    string id = 1;
    // Version
    int64 version = 2;
    // ResourceID of the resource.
    string resource_id = 3;
    // Namespace for resource.
    string namespace = 4;
    // Principal that is using the resource.
    string principal_id = 5;
    // state of resource instance.
    ResourceState state = 6;
    // Time duration in milliseconds after which instance will expire.
    google.protobuf.Duration expiry = 7;
}
```
### Permission

The Permission defines access policies for a resource including dynamic conditions based 
on [GO Templates](https://pkg.go.dev/text/template) that are evaluated before granting an access:

```protobuf3
enum Effect {
    PERMITTED = 0;
    DENIED = 1;
}

message Permission {
    // ID unique identifier assigned to this permission.
    string id = 1;
    // Version
    int64 version = 2;
    // Namespace for permission.
    string namespace = 3;
    // Scope for permission.
    string scope = 4;
    // Actions that can be performed.
    repeated string actions = 5;
    // Resource for the action.
    string resource_id = 6;
    // Effect Permitted or Denied
    Effect effect = 7;
    // Constraints expression with dynamic properties.
    string constraints = 8;
}
```

### Role

A Principal can be associated with one or more Roles where each Role has a name and can be optionally associated 
with Permissions for implementing RBAC based access control, e.g.,
```protobuf3
message Role {
    // ID unique identifier assigned to this role.
    string id = 1;
    // Version
    int64 version = 2;
    // Namespace for permission.
    string namespace = 3;
    // Name of the role.
    string name = 4;
    // PermissionIDs that can be performed.
    repeated string permission_ids = 5;
    // Optional parent ids
    repeated string parent_ids = 6;
}
```

A Role can also be inherited from multiple other Roles so that common Permissions can be defined in the parent Role(s) 
and specific Permissions are defined in the derived Roles.

### Group

A Principal can be linked to multiple Groups, and each Group can be tied to several Roles. The Principal inherits access 
Permissions not only directly associated with it but also from the Roles it’s part of and the Groups it’s connected to.
Here is the Group definition:

```protobuf3
message Group {
    // ID unique identifier assigned to this group.
    string id = 1;
    // Version
    int64 version = 2;
    // Namespace for permission.
    string namespace = 3;
    // Name of the group.
    string name = 4;
    // RoleIDs that are associated.
    repeated string role_ids = 5;
    // Optional parent ids.
    repeated string parent_ids = 6;
}
```
A Group can also have one or parents similar to Roles so that access rules policies can check membership for groups or 
inherits all permissions that belong to a Group through its association with Roles.

### Relationship

A Principal can define relationships with resources or target objects for performing actions and access policies can 
check for existence of a relationship before permitting an action and implementing ReBAC based policies. Though, 
Relationship seems similar to a Role or a Group but it differs from them because a Relationship directly associate 
between a Principal and a Resource where as a Role can be associated with multiple Principals and is indirectly 
associated with Resource through Permission object. Here is the definition for a Relationship:

```protobuf3
message Relationship {
    // ID unique identifier assigned to this relationship.
    // in:body
    string id = 1;

    // Version
    // in:body
    int64 version = 2;

    // Namespace for permission.
    // in:body
    string namespace = 3;

    // Relation name.
    // in:body
    string relation = 4;

    // PrincipalID for relationship.
    // in:body
    string principal_id = 5;

    // ResourceID for relationship.
    // in:body
    string resource_id = 6;

    // Attributes of relationship.
    // in:body
    map<string, string> attributes = 7;
}
```


## API Specifications for Authorization

The Authorization APIs are grouped into control-plane APIs for managing above data and their relationships 
with Principals and data-plane (behavioral) for Authorizing decisions. Following section defines 
control-plane APIs in [Protocol Buffers](https://protobuf.dev/overview/) definition language for 
managing authorization data and policies:

### Control-Plane APIs for managing Organizations

```protobuf3
service OrganizationsService {
    // Create Organizations swagger:route POST /api/v1/organizations organizations createOrganizationRequest
    // Responses:
    // 200: createOrganizationResponse
    rpc Create (CreateOrganizationRequest) returns (CreateOrganizationResponse);

    // Update Organizations swagger:route PUT /api/v1/organizations/{id} organizations updateOrganizationRequest
    // Responses:
    // 200: updateOrganizationResponse
    rpc Update (UpdateOrganizationRequest) returns (UpdateOrganizationResponse);

    // Get Organization swagger:route GET /api/v1/organizations/{id} organizations getOrganizationRequest
    // Responses:
    // 200: getOrganizationResponse
    rpc Get (GetOrganizationRequest) returns (GetOrganizationResponse);

    // Query Organization swagger:route GET /api/v1/organizations organizations queryOrganizationRequest
    // Responses:
    // 200: queryOrganizationResponse
    rpc Query (QueryOrganizationRequest) returns (stream QueryOrganizationResponse);

    // Delete Organization swagger:route DELETE /api/v1/organizations/{id} organizations deleteOrganizationRequest
    // Responses:
    // 200: deleteOrganizationResponse
    rpc Delete (DeleteOrganizationRequest) returns (DeleteOrganizationResponse);
}
```

Above definition also defines [OpenAPI](https://swagger.io/specification/) specification for REST based APIs so that the 
same behavior can be used by either [gRPC](https://grpc.io/) 
or [REST](https://en.wikipedia.org/wiki/Overview_of_RESTful_API_Description_Languages) API protocols.

### Control-Plane APIs for managing Principals

Following specification defines APIs to manage Principals and add/remove the associations with Roles, Groups, Permissions 
and Relationships:

```protobuf3
service PrincipalsService {
    // Create Principals swagger:route POST /api/v1/{organization_id}/principals principals createPrincipalRequest
    // Responses:
    // 200: createPrincipalResponse
    rpc Create (CreatePrincipalRequest) returns (CreatePrincipalResponse);

    // Update Principals swagger:route PUT /api/v1/{organization_id}/principals/{id} principals updatePrincipalRequest
    // Responses:
    // 200: updatePrincipalResponse
    rpc Update (UpdatePrincipalRequest) returns (UpdatePrincipalResponse);

    // Get Principal swagger:route GET /api/v1/{organization_id}/{namespace}/principals/{id} principals getPrincipalRequest
    // Responses:
    // 200: getPrincipalResponse
    rpc Get (GetPrincipalRequest) returns (GetPrincipalResponse);

    // Query Principal swagger:route GET /api/v1/{organization_id}/principals principals queryPrincipalRequest
    // Responses:
    // 200: queryPrincipalResponse
    rpc Query (QueryPrincipalRequest) returns (stream QueryPrincipalResponse);

    // Delete Principal swagger:route DELETE /api/v1/{organization_id}/principals/{id} principals deletePrincipalRequest
    // Responses:
    // 200: deletePrincipalResponse
    rpc Delete (DeletePrincipalRequest) returns (DeletePrincipalResponse);

    // AddGroups Principal swagger:route PUT /api/v1/{organization_id}/{namespace}/principals/{id}/groups/add principals addGroupsToPrincipalRequest
    // Responses:
    // 200: addGroupsToPrincipalResponse
    rpc AddGroups (AddGroupsToPrincipalRequest) returns (AddGroupsToPrincipalResponse);

    // DeleteGroups Principal swagger:route PUT /api/v1/{organization_id}/{namespace}/principals/{id}/groups/delete principals deleteGroupsToPrincipalRequest
    // Responses:
    // 200: deleteGroupsToPrincipalResponse
    rpc DeleteGroups (DeleteGroupsToPrincipalRequest) returns (DeleteGroupsToPrincipalResponse);

    // AddRoles Principal swagger:route PUT /api/v1/{organization_id}/{namespace}/principals/{id}/roles/add principals addRolesToPrincipalRequest
    // Responses:
    // 200: addRolesToPrincipalResponse
    rpc AddRoles (AddRolesToPrincipalRequest) returns (AddRolesToPrincipalResponse);

    // DeleteRole Principal swagger:route PUT /api/v1/{organization_id}/{namespace}/principals/{id}/roles/delete principals deleteRolesToPrincipalRequest
    // Responses:
    rpc DeleteRoles (DeleteRolesToPrincipalRequest) returns (DeleteRolesToPrincipalResponse);

    // AddPermissions Principal swagger:route PUT /api/v1/{organization_id}/{namespace}/principals/{id}/permissions/add principals addPermissionsToPrincipalRequest
    // Responses:
    // 200: addPermissionsToPrincipalResponse
    rpc AddPermissions (AddPermissionsToPrincipalRequest) returns (AddPermissionsToPrincipalResponse);

    // DeletePermissions Principal swagger:route PUT /api/v1/{organization_id}/{namespace}/principals/{id}/permissions/delete principals deletePermissionsToPrincipalRequest
    // Responses:
    // 200: deletePermissionsToPrincipalResponse
    rpc DeletePermissions (DeletePermissionsToPrincipalRequest) returns (DeletePermissionsToPrincipalResponse);

    // AddRelationships Principal swagger:route PUT /api/v1/{organization_id}/{namespace}/principals/{id}/relations/add principals addRelationshipsToPrincipalRequest
    // Responses:
    rpc AddRelationships (AddRelationshipsToPrincipalRequest) returns (AddRelationshipsToPrincipalResponse);

    // DeleteRelationships Principal swagger:route PUT /api/v1/{organization_id}/{namespace}/principals/{id}/relations/delete principals deleteRelationshipsToPrincipalRequest
    // Responses:
    // 200: deleteRelationshipsToPrincipalResponse
    rpc DeleteRelationships (DeleteRelationshipsToPrincipalRequest) returns (DeleteRelationshipsToPrincipalResponse);
}
```

Above definition defines [OpenAPI](https://swagger.io/specification/) specification for REST based APIs as well for 
providing groups management API using [gRPC](https://grpc.io/) 
or [REST](https://en.wikipedia.org/wiki/Overview_of_RESTful_API_Description_Languages) API protocols.

### Control-Plane APIs for managing Resources
```protobuf3
service ResourcesService {
    // Create Resources swagger:route POST /api/v1/{organization_id}/{namespace}/resources resources createResourceRequest
    // Responses:
    // 200: createResourceResponse
    rpc Create (CreateResourceRequest) returns (CreateResourceResponse);

    // Update Resources swagger:route PUT /api/v1/{organization_id}/{namespace}/resources/{id} resources updateResourceRequest
    // Responses:
    // 200: updateResourceResponse
    rpc Update (UpdateResourceRequest) returns (UpdateResourceResponse);

    // Query Resource swagger:route GET /api/v1/{organization_id}/{namespace}/resources resources queryResourceRequest
    // Responses:
    // 200: queryResourceResponse
    rpc Query (QueryResourceRequest) returns (stream QueryResourceResponse);

    // Delete Resource swagger:route DELETE /api/v1/{organization_id}/{namespace}/resources/{id} resources deleteResourceRequest
    // Responses:
    // 200: deleteResourceResponse
    rpc Delete (DeleteResourceRequest) returns (DeleteResourceResponse);

    // CountResourceInstances Resources swagger:route GET /api/v1/{organization_id}/{namespace}/resources/{id}/instance_count resources countResourceInstancesRequest
    // Responses:
    // 200: countResourceInstancesResponse
    rpc CountResourceInstances (CountResourceInstancesRequest) returns (CountResourceInstancesResponse);

    // QueryResourceInstances Resources swagger:route GET /api/v1/{organization_id}/{namespace}/resources/{id}/instances resources queryResourceInstanceRequest
    // Responses:
    // 200: queryResourceInstanceResponse
    rpc QueryResourceInstances (QueryResourceInstanceRequest) returns (stream QueryResourceInstanceResponse);
}
```

### Control-Plane APIs for managing Groups
```protobuf3
service GroupsService {
    // Create Groups swagger:route POST /api/v1/{organization_id}/{namespace}/groups groups updateGroupRequest
    // Responses:
    // 200: updateGroupResponse
    rpc Create (CreateGroupRequest) returns (CreateGroupResponse);

    // Update Groups swagger:route PUT /api/v1/{organization_id}/{namespace}/groups groups/{id} updateGroupRequest
    // Responses:
    // 200: updateGroupResponse
    rpc Update (UpdateGroupRequest) returns (UpdateGroupResponse);

    // Query Group swagger:route GET /api/v1/{organization_id}/{namespace}/groups groups queryGroupRequest
    // Responses:
    // 200: queryGroupResponse
    rpc Query (QueryGroupRequest) returns (stream QueryGroupResponse);

    // Delete Group swagger:route DELETE /api/v1/{organization_id}/{namespace}/groups/{id} groups deleteGroupRequest
    // Responses:
    // 200: deleteGroupResponse
    rpc Delete (DeleteGroupRequest) returns (DeleteGroupResponse);

    // AddRoles Group swagger:route PUT /api/v1/{organization_id}/{namespace}/groups/{id}/roles/add groups addRolesToGroupRequest
    // Responses:
    // 200: addRolesToGroupResponse
    rpc AddRoles (AddRolesToGroupRequest) returns (AddRolesToGroupResponse);

    // DeleteRoles Group swagger:route PUT /api/v1/{organization_id}/{namespace}/groups/{id}/roles/delete groups deleteRolesToGroupRequest
    // Responses:
    // 200: deleteRolesToGroupResponse
    rpc DeleteRoles (DeleteRolesToGroupRequest) returns (DeleteRolesToGroupResponse);
}
```
### Control-Plane APIs for managing Roles
```protobuf3
service RolesService {
    // Create Roles swagger:route POST /api/v1/{organization_id}/{namespace}/roles roles createRoleRequest
    // Responses:
    // 200: createRoleResponse
    rpc Create (CreateRoleRequest) returns (CreateRoleResponse);

    // Update Roles swagger:route PUT /api/v1/{organization_id}/{namespace}/roles/{id} roles updateRoleRequest
    // Responses:
    // 200: updateRoleResponse
    rpc Update (UpdateRoleRequest) returns (UpdateRoleResponse);

    // Query Role swagger:route GET /api/v1/{organization_id}/{namespace}/roles roles queryRoleRequest
    // Responses:
    // 200: queryRoleResponse
    rpc Query (QueryRoleRequest) returns (stream QueryRoleResponse);

    // Delete Role swagger:route DELETE /api/v1/{organization_id}/{namespace}/roles/{id} roles deleteRoleRequest
    // Responses:
    // 200: deleteRoleResponse
    rpc Delete (DeleteRoleRequest) returns (DeleteRoleResponse);

    // AddPermissions Role swagger:route PUT /api/v1/{organization_id}/{namespace}/roles/{id}/permissions/add roles addPermissionsToRoleRequest
    // Responses:
    // 200: addPermissionsToRoleResponse
    rpc AddPermissions (AddPermissionsToRoleRequest) returns (AddPermissionsToRoleResponse);

    // DeletePermissions Role swagger:route PUT /api/v1/{organization_id}/{namespace}/roles/{id}/permissions/delete roles deletePermissionsToRoleRequest
    // Responses:
    // 200: deletePermissionsToRoleResponse
    rpc DeletePermissions (DeletePermissionsToRoleRequest) returns (DeletePermissionsToRoleResponse);
}
```

### Control-Plane APIs for managing Permissions
```protobuf3
service PermissionsService {
    // Create Permissions swagger:route POST /api/v1/{organization_id}/{namespace}/permissions permissions createPermissionRequest
    // Responses:
    // 200: createPermissionResponse
    rpc Create (CreatePermissionRequest) returns (CreatePermissionResponse);

    // Update Permissions swagger:route PUT /api/v1/{organization_id}/{namespace}/permissions/{id} permissions updatePermissionRequest
    // Responses:
    // 200: updatePermissionResponse
    rpc Update (UpdatePermissionRequest) returns (UpdatePermissionResponse);

    // Query Permission swagger:route GET /api/v1/{organization_id}/{namespace}/permissions permissions queryPermissionRequest
    // Responses:
    // 200: queryPermissionResponse
    rpc Query (QueryPermissionRequest) returns (stream QueryPermissionResponse);

    // Delete Permission swagger:route DELETE /api/v1/{organization_id}/{namespace}/permissions/{id} permissions deletePermissionRequest
    // Responses:
    // 200: deletePermissionResponse
    rpc Delete (DeletePermissionRequest) returns (DeletePermissionResponse);
}
```

### Control-Plane APIs for managing Relationships
```protobuf3
service RelationshipsService {
    // Create Relationships swagger:route POST /api/v1/{organization_id}/{namespace}/relations relationships createRelationshipRequest
    // Responses:
    // 200: createRelationshipResponse
    rpc Create (CreateRelationshipRequest) returns (CreateRelationshipResponse);

    // Update Relationships swagger:route PUT /api/v1/{organization_id}/{namespace}/relations/{id} relationships updateRelationshipRequest
    // Responses:
    // 200: updateRelationshipResponse
    rpc Update (UpdateRelationshipRequest) returns (UpdateRelationshipResponse);

    // Query Relationship swagger:route GET /api/v1/{organization_id}/{namespace}/relations relationships queryRelationshipRequest
    // Responses:
    // 200: queryRelationshipResponse
    rpc Query (QueryRelationshipRequest) returns (stream QueryRelationshipResponse);

    // Delete Relationship swagger:route DELETE /api/v1/{organization_id}/{namespace}/relations/{id} relationships deleteRelationshipRequest
    // Responses:
    // 200: deleteRelationshipResponse
    rpc Delete (DeleteRelationshipRequest) returns (DeleteRelationshipResponse);
}
```

### Data-Plane APIs for Authorization

Following specification defines APIs for authorizing access to resources based on permissions and constraints as 
well operations to allocate and deallocate resources:

```protobuf3
service AuthZService {
    // Authorize swagger:route POST /api/v1/{organization_id}/{namespace}/{principal_id}/auth authz authRequest
    // Responses:
    // 200: authResponse
    rpc Authorize (AuthRequest) returns (AuthResponse);

    // Check swagger:route POST /api/v1/{organization_id}/{namespace}/{principal_id}/auth/constraints authz checkConstraintsRequest
    // Responses:
    // 200: checkConstraintsResponse
    rpc Check (CheckConstraintsRequest) returns (CheckConstraintsResponse);

    // Allocate Resources swagger:route PUT /api/v1/{organization_id}/{namespace}/resources/{id}/allocate/{principal_id} resources allocateResourceRequest
    // Responses:
    // 200: allocateResourceResponse
    rpc Allocate (AllocateResourceRequest) returns (AllocateResourceResponse);

    // Deallocate Resources swagger:route PUT /api/v1/{organization_id}/{namespace}/resources/{id}/deallocate/{principal_id} resources deallocateResourceRequest
    // Responses:
    // 200: deallocateResourceResponse
    rpc Deallocate (DeallocateResourceRequest) returns (DeallocateResourceResponse);
}
```

#### Authorize API

The Authorize API takes AuthRequest as a request that defines Principal-Id, Resource-Name, Action and context attributes 
and checks permissions for granting access:

```protobuf3
message AuthRequest {
    string organization_id = 1;
    string namespace = 2;
    string principal_id = 3;
    string action = 4;
    string resource = 5;
    string scope = 6;
    map<string, string> context = 7;
}
message AuthResponse {
    api.authz.types.Effect effect = 1;
    string message = 2;
}
```

#### Check Constraints API

The Check API allows evaluating dynamic conditions based on [GO Templates](https://pkg.go.dev/text/template) without 
defining Permissions so that you can check for the membership to a group, a role, an existence of a relationship 
or other dynamic properties.

```protobuf3
message CheckConstraintsRequest {
    string organization_id = 1;
    string namespace = 2;
    string principal_id = 3;
    string constraints = 4;
    map<string, string> context = 5;
}
message CheckConstraintsResponse {
    bool matched = 1;
    string output = 2;
}
```

#### Allocate and Deallocate Resources APIs

The Allocate and Deallocate APIs can be used to manage resources that can be assigned based on a quota or a maximum capacity, e.g.:

```protobuf3
message AllocateResourceRequest {
    string organization_id = 1;
    string namespace = 2;
    string resource_id = 3;
    string principal_id = 4;
    string constraints = 5;
    google.protobuf.Duration expiry = 6;
    map<string, string> context = 7;
}
```

## Implementation
PlexAuthZ implements above hybrid authorization APIs. The following diagram illustrates structure of modules for the 
implementing various parts of the Authorization system:

![](https://weblog.plexobject.com/images/plexauthz-pkg.png)

Following are major components in above diagram:

### API Layer

The API layer defines service interfaces and schema for domain model as well request/response objects. The 
interfaces are then implemented by [gRPC](https://grpc.io/) servers 
and [REST](https://en.wikipedia.org/wiki/Overview_of_RESTful_API_Description_Languages) controllers.

### Data Layer and Repositories

The Data layer defines interfaces for storing data in Redis or DynamoDB databases. The Repository layer defines 
interfaces for managing data for each type such as Principal, Organization and Resource.

### Domain Services

The Domain services abstract over Repository layer and implements referential integrity between data objects and 
validation logic before persisting authorization data.

### Authorizer

The Authorizer layer defines interfaces for Authorization decisions. The API layer implements the interface 
based on Casbin for communicating clients and servers. This layer defines a default implementation based on the 
Domain service layer for enforcing authorization decisions based on above APIs.

### Factory and Configuration

The PlexAuthZ makes extensive use of interfaces with different implementations for Datastore, Repositories, 
Authorizer and AuthAdapter. The user can choose different implementations based on the Configuration, which 
are passed to the factory methods when instantiating objects that implement those interfaces.

### AuthAdapter

The AuthAdapter abstracts Data services and Authorizer for interacting with underlying Authorization system. 
AuthAdapter defines a simplified DSL in GO language that understands the relationships between data objects. 
The users can instantiate AuthAdapter that can connect to remote [gRPC](https://grpc.io/) server, 
[REST](https://en.wikipedia.org/wiki/Overview_of_RESTful_API_Description_Languages) controller, or the database directly.

## Usage Examples

In above data model and APIs, Principals, Resources and Relationships can have arbitrary attributes that can 
be checked at runtime for enforcing policies based on attributes. In addition, the request objects for Authorize, 
Check and AllocateResource defines runtime context properties that can be passed along with other attributes when 
evaluating runtime conditions based on [GO Templates](https://pkg.go.dev/text/template). 

Following section defines use-cases for enforcing access policies based 
on [ABAC](https://csrc.nist.gov/Projects/Attribute-Based-Access-Control), [RBAC](https://csrc.nist.gov/projects/role-based-access-control), 
[ReBAC](https://en.wikipedia.org/wiki/Relationship-based_access_control) and [PBAC](https://csrc.nist.gov/glossary/term/policy_based_access_control):

### GO Client Initialization

First, the GO client library will be setup with a selection of the implementation based on the database, 
[gRPC](https://grpc.io/) client or [REST API](https://en.wikipedia.org/wiki/Overview_of_RESTful_API_Description_Languages) 
client, e.g.,
```go
cfg, err := domain.NewConfig("") // omitting error handling here
// config defines mode for access by database, gRPC or REST APIs.
authService, _, err := factory.CreateAuthAdminService(cfg)
authorizer, err := authz.CreateAuthorizer(authz.DefaultAuthorizerKind, cfg, authSvc)

authAdapter := client.New(authorizer, authService)
orgAdapter, err := authAdapter.CreateOrganization(
&types.Organization{
Name:       "xyz-corp",
Namespaces: []string{"marketing", "sales"},
})
namespace := orgAdapter.Organization.Namespaces[0]

```

### Attributes based Access Policies

Following example illustrates implementing attribute-based access policies where three Principals (alice, bob, charlie) 
will define attributes for Department and Rank:
```go
alice, err := orgAdapter.Principals().WithUsername("alice").
WithAttributes("Department", "Engineering", "Rank", "5").Create()
bob, err := orgAdapter.Principals().WithUsername("bob").
WithAttributes("Department", "Engineering", "Rank", "6").Create()
charlie, err := orgAdapter.Principals().WithUsername("charlie").
WithAttributes("Department", "Sales", "Rank", "6").Create()
```

Then a resource for an ios-app and permissions will be defined as follows:
```go
app, err := orgAdapter.Resources(namespace).WithName("ios-app").
WithAttributes("Editors", "alice bob").
WithActions("list", "read", "write", "create", "delete").Create()

rlPerm1, err := orgAdapter.Permissions(namespace).WithResource(app.Resource).
WithConstraints(`
{{or (Includes .Resource.Editors .Principal.Username) (GE .Principal.Rank 6)}}
`).WithActions("read", "list").Create()

wPerm2, err := orgAdapter.Permissions(namespace).WithResource(app.Resource).
WithConstraints(`
{{and (Includes .Resource.Editors .Principal.Username) (GE .Principal.Rank 6)}}
`).WithActions("write").Create()

// assigning permission to all principals
alice.AddPermissions(rlPerm1.Permission, wPerm2.Permission)
bob.AddPermissions(rlPerm1.Permission, wPerm2.Permission))
charlie.AddPermissions(rlPerm1.Permission, wPerm2.Permission)
```

Then check the permissions as follows:
```go
// Alice, Bob, Charlie should be able to read/list since alice/bob belong to
// Editors attribute and Charlie's rank >= 6
require.NoError(t, alice.Authorizer(namespace).WithAction("list").
WithResourceName("ios-app").Check())
require.NoError(t, bob.Authorizer(namespace).WithAction("list").
WithResourceName("ios-app").Check())
require.NoError(t, charlie.Authorizer(namespace).WithAction("list").
WithResourceName("ios-app").Check())

// Only Bob should be able to write because Alice's rank is lower than 6 and
// Charlie doesn't belongto Editors attribute.
require.Error(t, alice.Authorizer(namespace).WithAction("write").
WithResourceName("ios-app").Check())
require.NoError(t, bob.Authorizer(namespace).WithAction("write").
WithResourceName("ios-app").Check())
require.Error(t, charlie.Authorizer(namespace).WithAction("write").
WithResourceName("ios-app").Check())
```

**Note**: The Authorization adapter defines Check method that will invoke the Authorize or Check method of the data-plane 
Authorization API based on parameters.

### Runtime Attributes based on IPAddresses

The GO Templates allow defining custom functions and [PlexAuthz](https://github.com/bhatti/PlexAuthZ) implementation 
includes a number of helper functions to validate IP addresses, Geolocation, Time and other environment factors, e.g.,
```go
rwlPerm, err := orgAdapter.Permissions(namespace).
WithResource(app.Resource).
WithConstraints(`
{{$Loopback := IsLoopback .IPAddress}}
{{$Multicast := IsMulticast .IPAddress}}
{{and (not $Loopback) (not $Multicast) (IPInRange .IPAddress "211.211.211.0/24")}}
`).WithActions("read", "write", "list").Create()
alice.AddPermissions(rwlPerm.Permission)

// The app should be only be accessible if ip-address is not loop-back,
// not multi-cast and within ip-range
require.NoError(t, alice.Authorizer(namespace).WithAction("list").
WithContext("IPAddress", "211.211.211.5").WithResourceName("ios-app").Check())
// But not local ipaddress or multicast
require.Error(t, alice.Authorizer(namespace).WithAction("list").
WithContext("IPAddress", "127.0.0.1").WithResourceName("ios-app").Check())
require.Error(t, alice.Authorizer(namespace).WithAction("list").
WithContext("IPAddress", "224.0.0.1").WithResourceName("ios-app").Check())
```

### RBAC Scenario

The following example will assign roles and groups to Principal objects and then enforce membership before granting the access:
```go
teller, err := orgAdapter.Roles(namespace).WithName("Teller").Create()
manager, err := orgAdapter.Roles(namespace).WithName("Manager").WithParents(teller.Role).Create()
loanOfficer, err := orgAdapter.Roles(namespace).WithName("LoanOfficer").Create()
support, err := orgAdapter.Roles(namespace).WithName("ITSupport").Create()

// Assigning roles
alice.AddRoles(manager.Role)
bob.AddRoles(loanOfficer.Role)
charlie.AddRoles(support.Role)

sales, err := orgAdapter.Groups(namespace).WithName("Sales").Create()
accounting, err := orgAdapter.Groups(namespace).WithName("Accounting").Create()
engineering, err := orgAdapter.Groups(namespace).WithName("Engineering").Create()

// Assigning groups
alice.AddGroups(sales.Group)
bob.AddGroups(accounting.Group)
charlie.AddGroups(engineering.Group)
```

Following snippet illustrates enforcement of roles and groups membership:
```go
require.NoError(t, alice.Authorizer(namespace).WithConstraints(
`{{and (HasRole "Teller") (HasGroup "Sales") (TimeInRange .CurrentTime .StartTime .EndTime)}}`).
WithContext("CurrentTime", "10:00am", "StartTime", "8:00am", "EndTime", "4:00pm").Check())

require.NoError(t, bob.Authorizer(namespace).WithConstraints(
`{{and (HasRole "LoanOfficer") (HasGroup "Accounting") (TimeInRange .CurrentTime .StartTime .EndTime) (GT .Principal.EmploymentLength 1)}}`).
WithContext("CurrentTime", "10:00am", "StartTime", "8:00am", "EndTime", "4:00pm").Check())

require.NoError(t, charlie.Authorizer(namespace).WithConstraints(`
{{and (HasRole "ITSupport") (HasGroup "Engineering") (TimeInRange .CurrentTime .StartTime .EndTime) (GT .Principal.EmploymentLength 1)}}`).
WithContext("CurrentTime", "10:00am", "StartTime", "8:00am", "EndTime", "4:00pm").Check())

// but should fail for bob because ITSupport Role and Engineering Group is required.
require.Error(t, bob.Authorizer(namespace).
WithConstraints(
`{{and (HasRole "ITSupport") (HasGroup "Engineering") (TimeInRange .CurrentTime .StartTime .EndTime) (GT .Principal.EmploymentLength 1)}}`).
WithContext("CurrentTime", "10:00am", "StartTime", "8:00am", "EndTime", "4:00pm").Check())
```
**Note**: The Authorizer adapter will invoke Check API in above use-cases because it’s only using constraints without 
defining permissions.

### ReBAC Scenario

Though, ReBAC systems generally define relationships between actors but you can consider a Principal as a 
subject-actor and a Resource as a target-actor for relationships. Following scenarios illustrates how relationships 
between Principal and Resources can be used to enforce ReBAC based access policies similar to Zanzibar:
```go
smith, err := orgAdapter.Principals().WithUsername("smith").
WithAttributes("UserRole", "Doctor").Create()
john, err := orgAdapter.Principals().WithUsername("john").
WithAttributes("UserRole", "Patient").Create()

medicalRecords, err := orgAdapter.Resources(namespace).WithName("MedicalRecords").
WithAttributes("Year", fmt.Sprintf("%d", time.Now().Year()), "Location", "Hospital").
WithActions("read", "write", "create", "delete").Create()

docRelation, err := smith.Relationships(namespace).WithRelation("AsDoctor").
WithResource(medicalRecords.Resource).
WithAttributes("Location", "Hospital").Create()

patientRelation, err := john.Relationships(namespace).WithRelation("AsPatient").
WithResource(medicalRecords.Resource).Create()

rwPerm, err := orgAdapter.Permissions(namespace).WithResource(medicalRecords.Resource).
WithConstraints(`
{{$CurrentYear := TimeNow "2006"}}
{{and (HasRelation "AsDoctor") (DistanceWithinKM .UserLatLng "46.879967,-121.726906" 100)
(eq .Resource.Year $CurrentYear) (eq .Resource.Location .Location)}}
`).WithActions("read", "write").Create()

rPerm, err := orgAdapter.Permissions(namespace).WithResource(medicalRecords.Resource).
WithScope("john's records").
WithConstraints(`
{{$CurrentYear := TimeNow "2006"}}
{{and (HasRelation "AsPatient") (eq .Resource.Year $CurrentYear) (eq .Resource.Location .Location)}}
`).WithActions("read").Create()

smith.AddPermissions(rwPerm.Permission)
john.AddPermissions(rPerm.Permission)
```

Above snippet defines medical-records as a resource, and Principals for smith and john where smith is assigned a 
relationship for AsDoctor and john is assigned a relationship for AsPatient. The permissions for reading or 
writing medical records enforce the AsDoctor relationship and permissions for reading medical records enforce the 
AsPatient relationship. Then enforcing relationships is defined as follows:

```go
// Dr. Smith should have permission for reading/writing medical records based on constraints
require.NoError(t, smith.Authorizer(namespace).WithAction("write").
WithResource(medicalRecords.Resource).
WithContext("UserLatLng", "47.620422,-122.349358", "Location", "Hospital").Check())

// Patient john should have permission for reading medical records based on constraints
require.NoError(t, john.Authorizer(namespace).WithAction("read").
WithScope("john's records").WithResource(medicalRecords.Resource).
WithContext("Location", "Hospital").Check())

// But Patient john should not write medical records
require.Error(t, john.Authorizer(namespace).WithAction("write").
WithResource(medicalRecords.Resource).
WithContext("Location", "Hospital").Check())
```

Above Snippet also makes use of other functions available in the Template language for enforcing dynamic conditions 
based on Geofencing that permits access only when the doctor is close to the Hospital.

As the Relationships are defined between actors, we can also define a Resource to represent a Doctor and a 
Principal for the patient so that a patient-doctor relationship can be established, e.g.,

```go
// Now treating Doctor as Target Resource for appointment
doctorResource, err := orgAdapter.Resources(namespace).WithName(smith.Principal.Name).
WithAttributes("Year", fmt.Sprintf("%d", time.Now().Year()),
"Location", "Hospital").WithActions("appointment", "consult").Create()

doctorPatientRelation, err := john.Relationships(namespace).WithRelation("Physician").
WithAttributes("StartTime", "8:00am", "EndTime", "4:00pm").
WithResource(doctorResource.Resource).Create()

apptPerm, err := orgAdapter.Permissions(namespace).WithResource(doctorResource.Resource).
WithConstraints(`
{{$CurrentYear := TimeNow "2006"}}
{{and (TimeInRange .AppointmentTime .Relations.Physician.StartTime .Relations.Physician.EndTime)
(HasRelation "Physician") (eq "Patient" .Principal.UserRole) (eq .Resource.Year $CurrentYear) (eq .Resource.Location .Location)}}
`).WithActions("appointment").Create()
john.AddPermissions(apptPerm.Permission)

// Patient john should be able to make appointment within normal Hopspital hours
require.NoError(t, john.Authorizer(namespace).WithAction("appointment").
WithResource(doctorResource.Resource).
WithContext("Location", "Hospital", "AppointmentTime", "10:00am").Check())
```

Above example shows how authorization rules can also limit access between the normal hours of appointments.

### Resources with Quota

[PlexAuthZ](https://github.com/bhatti/PlexAuthZ) supports defining access policies for resources that have 
quota, e.g., an organization may have a fixed set of IDE Licenses to be used by the engineering team or might be 
using a utility based computing resources with a daily budget. Here is an example scenario:

```go
engGroup, err := orgAdapter.Groups(namespace).WithName("Engineering").Create()

alice, err := orgAdapter.Principals().WithUsername("alice").
WithAttributes("Title", "Engineer", "Tenure", "3").Create()

// Assigning groups
alice.AddGroups(engGroup.Group)

// AND with following resources
ideLicences, err := orgAdapter.Resources(namespace).WithName("IDELicence").
WithCapacity(5).WithAttributes("Location", "Chicago").
WithActions("use").Create()

require.NoError(t, ideLicences.
WithConstraints(`and (GT .Principal.Tenure 1) (HasGroup "Engineering") (eq .Resource.Location .Location)`).
WithExpiration(time.Hour).WithContext("Location", "Chicago").Allocate(bob.Principal))
...
// Deallocate after use
require.NoError(t, ideLicences.Deallocate(alice.Principal))
```

Above example demonstrates that the IDE License can only be allocated if the Principal is member of Engineering group, 
has a tenure of more than a year and Location matches Resource Location. In addition, the resource can be allocated 
only for a fixed duration and is automatically deallocated if not allocated explicitly. Both Redis and Dynamo DB 
supports TTL parameters for expiring data so no application logic is required to expire them.

### Resources with Wildcard in the name

[PlexAuthZ](https://github.com/bhatti/PlexAuthZ) supports resources with wildcards in the name so that a user can 
match permissions for all resources that match the wildcard pattern. Here is an example:

```go
alice, err := orgAdapter.Principals().WithUsername("alice").
WithAttributes("Department", "Sales", "Rank", "6").Create()
bob, err := orgAdapter.Principals().WithUsername("bob").
WithAttributes("Department", "Engineering", "Rank", "6").Create()

// Creating a project with wildcard
salesProject, err := orgAdapter.Resources(namespace).
WithName("urn:org-sales-*-project-1000-*").
WithAttributes("SalesYear", fmt.Sprintf("%d", time.Now().Year())).
WithActions("read", "write").Create()

rwlPerm, err := orgAdapter.Permissions(namespace).
WithResource(salesProject.Resource).
WithEffect(types.Effect_PERMITTED).
WithConstraints(`
{{$CurrentYear := TimeNow "2006"}}
{{and (GT .Principal.Rank 5) (eq .Principal.Department "Sales") (IPInRange .IPAddress "211.211.211.0/24") (eq .Resource.SalesYear $CurrentYear)}}
`).WithActions("*").Create()
require.NoError(t, err)

alice.AddPermissions(rwlPerm.Permission))
bob.AddPermissions(rwlPerm1.Permission)

// Alice should be able to access from Sales Department and complete project name
require.NoError(t, alice.Authorizer(namespace).WithAction("read").
WithResourceName("urn:org-sales-abc-project-1000-xyz").
WithContext("IPAddress", "211.211.211.5").Check())
// But bob should not be able to access project because he doesn't belong to the Sales Department
require.Error(t, bob.Authorizer(namespace).WithAction("read").
WithResourceName("urn:org-sales-abc-project-1000-xyz").
WithContext("IPAddress", "211.211.211.5").Check())
```

**Note:** The project name “urn:org-sales-abc-project-1000-xyz” matches the wildcard in resource name and permissions 
also verify attributes of the Resource and Principal.

### Permissions with Scope

[PlexAuthZ](https://github.com/bhatti/PlexAuthZ) allows associating permissions with specific Scope and the 
permission is only granted if the scope in authorization request at runtime matches the scope, e.g.,
```go
alice, err := orgAdapter.Principals().WithUsername("alice").
WithAttributes("Department", "Engineering", "Permanent", "true").Create()
bob, err := orgAdapter.Principals().WithUsername("bob").
WithAttributes("Department", "Sales", "Permanent", "true").Create()

project, err := orgAdapter.Resources(namespace).WithName("nextgen-app").
WithAttributes("Owner", "alice").
WithActions("list", "read", "write", "create", "delete").Create()

rwlPerm, err := orgAdapter.Permissions(namespace).WithResource(project.Resource).
WithScope("Reporting").
WithConstraints(`
{{or (eq .Principal.Username .Resource.Owner) (Not .Private)}}
`).WithActions("read", "write", "list").Create()

alice.AddPermissions(rwlPerm.Permission)
bob.AddPermissions(rwlPerm.Permission)

// Project should be only be accessible by alice as the scope matches and she is the owner.
require.NoError(t, alice.Authorizer(namespace).
WithAction("list").WithScope("Reporting").
WithContext("Private", "true").
WithResourceName("nextgen-app").Check())
// But alice should not be able to access without matching scope.
require.Error(t, alice.Authorizer(namespace).
WithAction("list").WithScope("").
WithContext("Private", "true").
WithResourceName("nextgen-app").Check())

// But bob should not be able to access as project is private and he is not the owner.
require.Error(t, bob.Authorizer(namespace).
WithAction("list").WithScope("Reporting").
WithContext("Private", "true").
WithResourceName("nextgen-app").Check())
```
**Note**: Above example also demonstrates show you can enforce ownership for private resources.

## Summary
-------

PlexAuthZ implements hybrid Authorization system that can support various forms of access policies based on 
on [ABAC](https://csrc.nist.gov/Projects/Attribute-Based-Access-Control), 
[RBAC](https://csrc.nist.gov/projects/role-based-access-control), 
[ReBAC](https://en.wikipedia.org/wiki/Relationship-based_access_control) and 
[PBAC](https://csrc.nist.gov/glossary/term/policy_based_access_control). 
It’s still early in development but feel free to try it and send your feedback or feature requests.