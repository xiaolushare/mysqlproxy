# mysqlproxy

MySQL Proxy server for PHP.

mysql command does not work.

## Usage

### Starting MySQL proxy server (root)

```
./mysqlproxy -root
```

### Starting MySQL proxy server

```
./mysqlproxy
```

### PHP Sample

```php
$link = mysql_connect(
	'/path/to/mysqlproxy.sock',
	'<db user>:<db password>@<proxy host>:<proxy port>;<db host>:<db port>',
	'<proxy password>',
);

// For example in following Data flow.

// Connect to A
$link = mysql_connect(
	'/path/to/mysqlproxy.sock',
	'user_a:******@192.168.1.1:9696;192.168.1.2:3306',
	'******',
);

// Connect to B
$link = mysql_connect(
	'/path/to/mysqlproxy.sock',
	'user_b:******@192.168.1.1:9696;192.168.1.3:3306',
	'******',
);

// Connect to C
$link = mysql_connect(
	'/path/to/mysqlproxy.sock',
	'user_c:******@192.168.2.1:9696;192.168.2.2:3306',
	'******',
);

// Connect to D
$link = mysql_connect(
	'/path/to/mysqlproxy.sock',
	'user_d:******@192.168.2.1:9696;192.168.2.3:3306',
	'******',
);
```

### Data flow

```
                 Unix domain socket                TLS                       TCP
                 Connect                           Connect                   Connect
+--------------+      +--------------------------+      +------------------+      +------------------+
| mysql client | ---> | mysql proxy              | -+-> | mysql proxy      | -+-> | mysql server     |
| (PHP)        |      | (root)                   |  |   |                  |  |   | (A)              |
|              |      | /path/to/mysqlproxy.sock |  |   | 192.168.1.1:9696 |  |   | 192.168.1.2:3306 |
+--------------+      +--------------------------+  |   +------------------+  |   +------------------+
                                                    |                         |                      
                                                    |                         |   +------------------+
                                                    |                         +-> | mysql server     |
                                                    |                             | (B)              |
                                                    |                             | 192.168.1.3:3306 |
                                                    |                             +------------------+
                                                    |                                                
                                                    |   +------------------+      +------------------+
                                                    +-> | mysql proxy      | -+-> | mysql server     |
                                                        |                  |  |   | (C)              |
                                                        | 192.168.2.1:9696 |  |   | 192.168.2.2:3306 |
                                                        +------------------+  |   +------------------+
                                                                              |                      
                                                                              |   +------------------+
                                                                              +-> | mysql server     |
                                                                                  | (D)              |
                                                                                  | 192.168.2.3:3306 |
                                                                                  +------------------+
```

