# Master & Slave (Leader & Follower)

- 執行arangosh，在裡面執行下方腳本
```sh
db._useDatabase("Database");
require("@arangodb/replication").setupReplication({
  endpoint: "tcp://{{domain or IPs}}:8529",
  username: "root",
  password: "",
  autoStart: true,
  autoResync: true,
  autoResyncRetries: 2,
  adaptivePolling: true,
  includeSystem: false,
  requireFromPresent: false,
  idleMinWaitTime: 0.5,
  idleMaxWaitTime: 1.5,
  verbose: false
});
```

- 啟動replication applier並指定紀錄點
```sh
require("@arangodb/replication").applier.start(<tick>);
```

- 每次dump下來都會有記錄點
  <tick> 記錄點

- 停止replication applier
```sh
require("@arangodb/replication").applier.stop();
```

- 檢查replication applier 狀態
```
require("@arangodb/replication").applier.state()
```

- 帶有同步檢查（速度較慢，但更容易）
- 沒有同步檢查（速度更快，但需要正確提供最後一個服務器的記錄點）

- reference
  https://www.arangodb.com/docs/stable/administration-master-slave.html


# Active Failover

## Define
- One ArangoDB Single-Server instance which is read / writable by clients called Leader
- One or more ArangoDB Single-Server instances, which are passive and not writable called Followers, which asynchronously replicate data from the Leader
- At least one Agency acting as a “witness” to determine which server becomes the leader in a failure situation

## Architecture
![Active Failover](https://www.arangodb.com/docs/stable/images/leader-follower.png)

## Limitation
- 兩顆以上follower，active failover會選擇最新follower作為leader(best-effort盡力而為)
- 相比Arangodb cluster(synchronous replication),不能確保資料不會掉
- 機器之間升級版本後，可能發生偶發性錯誤

## 設定方式
1. 使用 ArangoDB starter
2. 手動

## Docker compose
- 包含一個follower與一個leader
- 包含一個agent
- Agent將監聽leader 與 follower 並進行failover （可以kill container測試）
- 如果leader與agent同時fail則follower無法成為leader。除此之外agent都仍持續進行failover
