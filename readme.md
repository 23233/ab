#### ab

* page 控制页码 page_size 控制条数
    * 最大均为100 100页 100条
* order(asc) order_desc
* search搜索 __会被替换为% search=__赵日天 会替换为 %赵日天
* filter_[字段名] 进行过滤 filter_id=1 最长64位请注意

#### 限制
* 目前不支持header为json的请求 只能是form 受限于iris解析

#### new version
* support read write split , use mysql storage and redis read!
* custom set cache time , hot point set 1 hour cache time .

#### process
* read
```
fetch -> redis -> has -> return
fetch -> redis -> not has -> mysql -> save to redis -> return
```
* update
```
req -> delete redis item -> mysql -> lazy(500ms) delete redis item
```
* delete
```
req -> mysql -> delete redis item
```
