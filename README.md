# mysqlproxy

MySQL Proxy server.

## 前準備

### ワークディレクトリの作成

```
sudo mkdir -p /path/to/workdir
```

### 公開鍵の設置

ワークディレクトリに設置する

### 設定ファイルの設置(必要な場合)

※ ルート(クライアントから接続するためのデーモン)で必要な場合のみ設置

```
vim /path/to/mysqlproxy.toml
["<接続先MySQLユーザー名1>"]
username = "<接続先MySQLユーザー名1>"
password = "<接続先MySQLパスワード1>"
proxyserver = "<接続先プロキシサーバのホスト1>"
["<接続先MySQLユーザー名2>"]
username = "<接続先MySQLユーザー名2>"
password = "<接続先MySQLパスワード2>"
proxyserver = "<接続先プロキシサーバのホスト2>"
…(繰り返し)
```

## Usage

### Connect to MySQL Server via MySQL proxy server

MySQL Proxy サーバーを経由してのMySQL接続方法

#### 設定ファイルを使用しない場合

```
mysql -S /path/to/mysqlproxy.sock -u <MySQLサーバーのユーザー名>(:<MySQLサーバーのパスワード>)@<プロキシサーバーのホスト>:<プロキシサーバーのポート>;<MySQLサーバーのホスト>:<MySQLサーバーのポート>(;<DB名>) -p
Enter password: <MySQLサーバーのパスワード>

※ ポートは3306番であっても必須。
※ パスワード、DB名は省略可能。パスワードがない場合、-pオプションは省略可能。
※ mysqlコマンドからだとユーザー名に文字数制限があるため、rdsのように長いドメインの場合は、
PHP等の各種プログラミング言語のMySQL接続アダプタを介せば接続可能です。
```

#### 設定ファイルを使用する場合

```
mysql -S /path/to/mysqlproxy.sock -u <MySQLサーバーのユーザー名>@<MySQLサーバーのホスト>(:<MySQLサーバーのポート>) -p
Enter password: <MySQLサーバーのパスワード>

※ プロキシサーバのアドレスは設定ファイルに書いているため省略可能ですが、ユーザー名とパスワードはプロキシサーバ自体の認証に必要です。
※ ポートが3306番であれば省略可能。
※ パスワードがない場合、-pオプションは省略可能。
※ mysqlコマンドからだとユーザー名に文字数制限があるため、rdsのように長いドメインの場合は、
PHP等の各種プログラミング言語のMySQL接続アダプタを介せば接続可能です。
```


### Starting MySQL proxy server (root)

クライアントから接続するためのデーモン

```
./mysqlproxy -root -workdir ワークディレクトリのパス -config 設定ファイルのパス
```

### Starting MySQL proxy server

MySQLサーバーに中継するためのデーモン

```
./mysqlproxy -workdir ワークディレクトリのパス
```

### PHP Sample

```php
$link = mysql_connect(
	'/path/to/mysqlproxy.sock',
	'<db user>:<db password>@<proxy host>:<proxy port>;<db host>:<db port>',
	'<db password>',
);

// For example in following Data flow.

// Connect to A
$link = mysql_connect(
	':/path/to/mysqlproxy.sock',
	'user_a:******@192.168.1.1:9696;192.168.1.2:3306',
	'******',
);

// Connect to B
$link = mysql_connect(
	':/path/to/mysqlproxy.sock',
	'user_b:******@192.168.1.1:9696;192.168.1.3:3306',
	'******',
);

// Connect to C
$link = mysql_connect(
	':/path/to/mysqlproxy.sock',
	'user_c:******@192.168.2.1:9696;192.168.2.2:3306',
	'******',
);

// Connect to D
$link = mysql_connect(
	':/path/to/mysqlproxy.sock',
	'user_d:******@192.168.2.1:9696;192.168.2.3:3306',
	'******',
);
```

### Data flow

```
           Unix domain socket   TLS                       TCP
           Connect              Connect                   Connect
+--------+      +-------------+      +------------------+      +------------------+
| mysql  | ---> | mysql proxy | -+-> | mysql proxy      | -+-> | mysql server     |
| client |      | (root)      |  |   |                  |  |   | (A)              |
| (PHP)  |      | localhost   |  |   | 192.168.1.1:9696 |  |   | 192.168.1.2:3306 |
+--------+      +-------------+  |   +------------------+  |   +------------------+
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

