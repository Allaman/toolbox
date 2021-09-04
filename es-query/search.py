#!/usr/bin/env python3
#-*- coding: utf-8 -*-

from elasticsearch import Elasticsearch, exceptions
import json

host = ""
index = "kubernetes_cluster-*"

client = Elasticsearch(host, retry_on_timeout=True, max_retries=3, timeout=60, request_timeout=60)

body = {
        "_source": ["time","log","kubernetes.pod_name","kubernetes.host"], # only return those keys from search result
        "query": {
            "bool": {
                "must": [
                    {
                        "match_phrase": {
                            "kubernetes.container_name": {
                                "query": "foo-bar"
                            }
                        }
                    },
                    {
                        "range": {
                            "time": {
                                "format": "strict_date_optional_time",
                                "gte": "2021-03-27T23:00:00.000Z",
                                "lte": "2021-03-28T21:59:59.000Z"
                            }
                        }
                    }
                ],
                "filter": [
                    {
                        "match_all": {}
                    }
                ],
                "should": [],
                "must_not": []
            }
        }
    }

def scroll(es, index, body, scroll, size, **kw):
    # print("Starting new scroll")
    page = es.search(index=index, body=body, scroll=scroll, size=size, **kw)
    scroll_id = page['_scroll_id']
    hits = page['hits']['hits']
    while len(hits):
        yield hits
        page = es.scroll(scroll_id=scroll_id, scroll=scroll)
        scroll_id = page['_scroll_id']
        hits = page['hits']['hits']

try:
    info = json.dumps(client.info(), indent=4)
    # print ("Elasticsearch client info():", info)
except exceptions.ConnectionError as err:
    print ("\nElasticsearch info() ERROR:", err)
    print ("\nThe client host:", host, "is invalid or cluster is not running")
    client = None

res = []
if client is not None:
    for hits in scroll(client, index, body, '2m', 10000):
        print(json.dumps(hits, indent=4))
        # res.append(json.dumps(hits))

# with open('out', 'w') as f:
#     f.write(''.join(res))
