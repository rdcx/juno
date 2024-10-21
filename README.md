# Juno

*Note: Juno is not yet production-ready, and many features are still a work in progress.*

Juno is a community-driven distributed search engine, where independent node operators power crawling and query operations. It offers programmatic search capabilities using the Monkey programming language, executed securely across 100,000 shards.

![Tokens Screenshot](juno-shot.png "Tokens Screenshot")

## System Overview

### Node Operations
- **Crawling**: Nodes are designed to handle crawling operations as simply as possible. Each node manages its own page and version databases, without direct communication with other nodes or dependencies on external APIs, except for the initialization of the shard list. 
- **Storage**: Each node guarantees a minimum storage capacity of 10GB, representing 1/100,000 of the expected 1 billion pages.

- **Query**: Queries can be of two types:
  - **Basic Query**: A straightforward request for information across the network.
  - **Script Query**: Allows users to execute custom scripts using the Monkey programming language for more complex data analysis and extraction.
  
  Every shard in the network is queried through a hierarchy of **query aggregators** and **range aggregators**. These aggregators perform a map-reduce-like operation, consolidating query results and efficiently piping the aggregated data back to the API. This enables users to receive responses from a vast distributed dataset seamlessly.

### Sharding Mechanism
- **Shards**: The Juno network is built on a foundation of 100,000 fixed shards, with each shard consisting of a minimum of one node, ideally three for redundancy and performance.
- **Node Assignment**: Nodes can self-assign to multiple shards, allowing flexibility and scalability, especially during the initial stages when data per shard is sparse. For example, a node may join with a shard range of [0, 999], managing 1,000 shards, or 1,000/100,000 of the expected data set. As the network grows, nodes can reallocate themselves from a broader range to a smaller one, enabling more efficient scaling.

### Automated Scaling
- **Dynamic Reallocation**: Nodes can dynamically reduce their shard coverage as the data volume increases. For instance, a node initially covering [0, 999] could reduce its range to [0, 499], freeing up 500 shards. When this happens, HTML content, page, and version metadata can either be removed or migrated to other nodes to ensure data availability.

### Architecture
- **Service Organization**: Juno employs a Domain-Driven Design (DDD) approach for organizing services. Each major function, like `pkg/api/node` or `pkg/balancer/crawl`, is modularized, with its own `domain.go` file. This structure enables easy migration of services into independent instances if needed, due to minimal coupling between components.
- **Repository Options**: Services support both in-memory and MySQL repository options, providing flexibility for different deployment scenarios.

### Testing Strategy
- **Test Improvements**: Current tests can be further refined to reduce some dependencies. However, the complexity of mocking certain services was deemed too high due to frequently changing requirements, so direct service testing was prioritized.

## API Operations

Juno provides an API for managing nodes, balancers, transactions, and tokens. Transactions and tokens are used to track the expenditure of customer queries and the earnings of node operators. 

### Nodes Management
- **Register Node**: Allows a node to register itself with the network, specifying its storage capacity and desired shard range.
- **Update Node Status**: Enables nodes to update their status, such as their current shard range or availability for new data.
- **Deallocate Shard**: Nodes can request to reduce their shard range, triggering the migration of data to other available nodes.

### Balancer Management
- **Add Balancer**: Adds a new balancer to manage traffic within a specified range of shards.
- **Update Balancer**: Allows updating the settings of a balancer, such as its shard responsibilities.
- **Remove Balancer**: Removes a balancer from the network, redistributing its responsibilities among other available balancers.

### Transactions and Tokens
- **Create Transaction**: Initiates a transaction to track a customer’s query expenditure.
- **Token Issuance**: Issues tokens to node operators as compensation for their services, based on the query processing they have completed.
- **Query Balance**: Allows customers and node operators to query their current balance of tokens.
- **Redeem Tokens**: Node operators can redeem tokens for payouts or other forms of compensation.

These API operations enable efficient management of nodes and balancers, ensuring that the Juno network remains scalable and responsive while providing transparent tracking of expenditures and earnings through transactions and tokens.
