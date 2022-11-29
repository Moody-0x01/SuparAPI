from requests import post, get

u = "http://localhost:8888/query"

res = get(u)

print(res.json())







