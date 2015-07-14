[![wercker status](https://app.wercker.com/status/918a2f54ea2bee6f8ec0c1e04c19ca79/m "wercker status")](https://app.wercker.com/project/bykey/918a2f54ea2bee6f8ec0c1e04c19ca79)

# Pez Inventory Service

The inventory capability of Pez.

### What things need to be managed as inventory?

#### Assumptions:
_(Need to validate with team)_

##### Leases:
Leases will be the mechanism required by inventory to reserve and release product(s) associated with a user.

* Binds product(s) to a user for a given slice of time.
* Leases may exist in different states (to be defined).  This list is not meant to be inclusive.
  1. **Pending** - A lease that is requested against a product type that is not currently available or able to be fulfilled
  2. **Active** - A lease that is able to be fulfilled, or has been fulfilled
  3. **Expired** - A lease with an expiration date <= today
  4. **Canceled** - A lease which is administratively canceled (e.g. separation, non-use, etc.)


##### Assumed User Experience (pertaining to inventory):
* User will be presented a list of products and whether they are available.
* User will select product(s) and submit a lease request.
* Inventory will reserve product associated with an active lease
* Inventory will release product associated with a canceled or expired lease
  - What additional tasks must occur before previously used inventory can be considered `available`?

### MVP 1
Manage the availability & consumption of pre-defined (one size fits all) sandboxes:

* [Capacity] The inventory service will know the number of pre-configured sandboxes.
* [Availability] It will know which sandboxes are leased.
* [Reserve] It will associate a sandbox to a lease.
* [Release] It will remove the lease from a sandbox.

### MVP 2
Manage the availability & consumption of virtual components (i.e. compute, memory, storage, & network) in support of variably-sized sandbox creation

### MVP 3
Manage the availability & consumption of physical & virtual component in support of Hybrid Labs offerings.

### MVP Future
* Discover Inventory via vCD APIs
* [Registration] Add inventory
* Support ad-hoc lab creation
* Save a lab assembly for future use
