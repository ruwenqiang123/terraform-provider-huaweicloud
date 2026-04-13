---
subcategory: "Elastic Load Balance (ELB)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_elb_loadbalancer_topology"
description: |-
  Use this data source to query the topology of a specified load balancer within HuaweiCloud.
---

# huaweicloud_elb_loadbalancer_topology

Use this data source to query the topology of a specified load balancer within HuaweiCloud.

## Example Usage

### Basic Usage

```hcl
variable "loadbalancer_id" {}

data "huaweicloud_elb_loadbalancer_topology" "test" {
  loadbalancer_id = var.loadbalancer_id
}
```

### Filter by listener ID

```hcl
variable "loadbalancer_id" {}
variable "listener_id" {}

data "huaweicloud_elb_loadbalancer_topology" "test" {
  loadbalancer_id = var.loadbalancer_id
  listener_id     = var.listener_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the load balancer topology.
  If omitted, the provider-level region will be used.

* `loadbalancer_id` - (Required, String) Specifies the ID of the load balancer to query topology.

* `listener_id` - (Optional, String) Specifies the listener ID to filter.

* `pool_id` - (Optional, String) Specifies the backend server group ID to filter.

* `listener_name` - (Optional, String) Specifies the listener name to filter.

* `listener_protocol` - (Optional, String) Specifies the listener protocol to filter.

* `listener_protocol_port` - (Optional, Int) Specifies the listener protocol port to filter.

* `pool_name` - (Optional, String) Specifies the backend server group name to filter.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `nodes` - The topology node information.  
  The [nodes](#loadbalancer_topology_nodes) structure is documented below.

* `edges` - The topology edge information.  
  The [edges](#loadbalancer_topology_edges) structure is documented below.

* `labels` - The topology label information.  
  The [labels](#loadbalancer_topology_labels) structure is documented below.

<a name="loadbalancer_topology_nodes"></a>
The `nodes` block supports:

* `loadbalancers` - The load balancer node information.  
  The [loadbalancers](#topology_nodes_loadbalancers) structure is documented below.

* `eips` - The EIP node information.  
  The [eips](#topology_nodes_eips) structure is documented below.

* `listeners` - The listener node information.  
  The [listeners](#topology_nodes_listeners) structure is documented below.

* `pools` - The backend server group node information.  
  The [pools](#topology_nodes_pools) structure is documented below.

<a name="topology_nodes_loadbalancers"></a>
The `loadbalancers` block supports:

* `id` - The ID of the load balancer.

* `name` - The name of the load balancer.

* `guaranteed` - Whether the load balancer is dedicated type.

* `l4_flavor_id` - The Layer 4 flavor ID of the load balancer.

* `l7_flavor_id` - The Layer 7 flavor ID of the load balancer.

* `vip_address` - The IPv4 address of the load balancer.

* `ipv6_vip_address` - The IPv6 address of the load balancer.

* `availability_zone_list` - The availability zone list of the load balancer.

<a name="topology_nodes_eips"></a>
The `eips` block supports:

* `id` - The ID of the EIP.

* `ip_address` - The IP address of the EIP.

* `ip_version` - The IP version of the EIP.  
  The valid values are as follows:
  + **4**: IPv4.
  + **6**: IPv6.

<a name="topology_nodes_listeners"></a>
The `listeners` block supports:

* `id` - The ID of the listener.

* `name` - The name of the listener.

* `protocol` - The protocol of the listener.

* `protocol_port` - The protocol port of the listener.

* `port_ranges` - The port ranges for full port listening.
  The [port_ranges](#topology_nodes_listeners_port_ranges) structure is documented below.

<a name="topology_nodes_listeners_port_ranges"></a>
The `port_ranges` block supports:

* `start_port` - The start port.

* `end_port` - The end port.

<a name="topology_nodes_pools"></a>
The `pools` block supports:

* `id` - The ID of the backend server group.

* `name` - The name of the backend server group.

* `protocol` - The protocol of the backend server group.

* `type` - The type of the backend server group.

* `lb_algorithm` - The load balancer algrithm of the backend server group.

<a name="loadbalancer_topology_edges"></a>
The `edges` block supports:

* `source` - The source node ID of the edge.

* `target` - The target node ID of the edge.

* `source_type` - The source node type of the edge.

* `target_type` - The target node type of the edge.

* `label` - The label of the edge.

* `label_id` - The label ID of the edge.

<a name="loadbalancer_topology_labels"></a>
The `labels` block supports:

* `l7policies` - The load balancer label list.
  The [l7policies](#topology_labels_l7policies) structure is documented below.

<a name="topology_labels_l7policies"></a>
The `l7policies` block supports:

* `id` - The forwarding policy ID.

* `name` - The forwarding policy name.

* `priority` - The forwarding policy priority.

* `action` - The forwarding policy action.

* `rules` - The List of forwarding rules.
  The [rules](#topology_labels_l7policies_l7rules) structure is documented below.

<a name="topology_labels_l7policies_l7rules"></a>
The `rules` block supports:

* `id` - The forwarding rule ID.

* `type` - The forwarding rule type.

* `compare_type` - The forwarding rule matching method.

* `key` - The match content key.

* `value` - The match content value.

* `provisioning_status` - The provisioning status.

* `invert` - Whether reverse matching is supported.

* `conditions` - The name of the load balancer.
  The [conditions](#l7policies_l7rules_conditions) structure is documented below.

* `created_at` - The creation time.

* `updated_at` - The update time.

<a name="l7policies_l7rules_conditions"></a>
The `conditions` block supports:

* `key` - The key of match item.

* `value` - The value of match item.
