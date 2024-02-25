# Entity Relationship Diagram (ERD) for Merchant and Product Database

This ERD shows the structure and relationships of the Merchant and Product tables in a relational database.

## Design Choice

- A relational database is chosen for this problem because it ensures data integrity and consistency, and support complex queries and joins.
- The data is normalized into two tables: Merchant and Product.
- The merchant ID and the SKU ID are used as the primary keys for the Merchant and Product tables respectively, because they are unique identifiers for each entity.
- The merchant ID primary key in the Merchant table is used as an index to improve query performance for the large database of merchants.
- A Merchant ID foreign key is included in the Product table, to improve the performance of queries that filter by merchant ID. This is because there could be a large number of products for each merchant, and indexing can speed up the search.
- A foreign key constraint is created between the merchant ID in the Product table and the merchant ID in the Merchant table, to enable joining the tables.

## Relationships and Cardinalities

- **Merchant** and **Product**: This is a one-to-many relationship. One merchant can sell many products, but one product can only be sold by one merchant. The cardinality of this relationship is 1:N.

## Diagram

The following diagram shows the ERD of the Merchant and Product database
````
+-----------------+     +-----------------+
| Merchant        |     | Product         |
+-----------------+     +-----------------+
| merchant_id (PK)|-----| merchant_id (FK)|
| name            |     | sku_id (PK)     |
| address         |     | name            |
| phone           |     | description     |
| email           |     | price           |
+-----------------+     | created_at      |
                        +-----------------+
