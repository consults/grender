# 渲染池
## 配置redis和mongo
redis用于添加任务，在配置文件中配置redis的地址、密码、数据库、redis——key

mongo用于存储渲染结果，配置保存的数据库、存储的表

## 使用python添加任务
```python
from typing import List

import redis
import json
from dataclasses import dataclass


class MyRedis(object):
    def __init__(self, host, port, password, db):
        self.redis_pool = redis.ConnectionPool(host=host, port=port, password=password, db=db,
                                               decode_responses=True)
        self.redis_conn = redis.Redis(connection_pool=self.redis_pool)


@dataclass
class Cook:
    Name: str
    Value: str
    Domain: str


@dataclass
class Task:
    Url: str
    Xpath: str
    TimeOut: int
    Cookies: List[Cook] = None


if __name__ == '__main__':
    host = '117.50.175.64'
    port = "6379"
    password = "jhkdjhkjdhsIUTYURTU_688J6j"
    db = 0
    rds = MyRedis(host=host, port=port, password=password, db=db)
    task = Task(Url="http://httpbin.org/ip", Xpath="//body",
                TimeOut=10)
    c = Cook(Name="sessionid", Value="n1olv6rch6s9hsqafkoubj8lq0l1zbnw", Domain="www.python-spider.com")
    # task.Cookies = [c.__dict__]
    rds.redis_conn.lpush("test", json.dumps(task.__dict__))

```