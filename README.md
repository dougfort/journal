# journal
storage for gm-control transactions

## synopsis

This package defines an interface for backing up transactions to the
gm-control-api back end. 

Transactions are stored in an append only manner using a flat file, 
gm-data or some other permanent storage.

https://en.wikipedia.org/wiki/Write-ahead_logging

https://engineering.linkedin.com/distributed-systems/log-what-every-software-engineer-should-know-about-real-time-datas-unifying