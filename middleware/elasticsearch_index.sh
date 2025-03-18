# 최초 인덱스 생성
curl -u "elastic:elastic" -k -X PUT "https://localhost:30090/event-000001"  -H "Content-Type: application/json" -d '{
    "aliases": {
        "event": {
            "is_write_index": true
        }
    }
}'
