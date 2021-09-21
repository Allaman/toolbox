# ES-Query

There are various methods for searching ES.

[Simple](simple) performs a simple search that is limited to the size of returned documents. Be aware that there is a maximum query size on ES side that you can not affect on the client side. A query that would produce more results will not return all documents!

[Scroll](scroll) performs a scrolling search that basically paginates through all result pages and therefore overcoming the limits of a simple search.
