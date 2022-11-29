"""
Script to automate server testing.
"""

from requests import get, post
from json import dumps, dump, loads
from sys import argv
import sqlite3

host = "http://localhost:8888/ws"

Query = "update PRODUCTS set MARK = '' where id = '29';"
conn = sqlite3.connect("/home/hz/server/hztech.db")
cursor = conn.cursor()



def QueryServer(CMD, payload0=None) -> None:
	Headers = {'Content-Type': 'application/json'}
	if payload0:
		Pyload =  {'Attributes': {'Command': CMD, **payload0}}
	else:
		Pyload =  {'Attributes': {'Command': CMD}}
	res = post(host, data=dumps(Pyload), headers=Headers)
	childs = res.json()['Childs']
	print(childs)
	for i in childs:
		id, Name = i["Attributes"]["ID"], i["Attributes"]["Name"]
		if not i["Attributes"]["MARK"].strip():
			cursor.execute(f"update PRODUCTS set MARK = 'inconnue' where id = '{id}'")
			print(f"{Name} was fixed! id: {id}")
	print("Done:::::")
		
			
if __name__ == '__main__':
	if len(argv) > 2:
		QueryServer(argv[1], payload0=loads(argv[2]))
	else:
		QueryServer(argv[1])
	
	conn.commit()
	conn.close()



# json = {'FullName': 'sijfsd', 'Email': 'testing....', 'Message': 'Hi'}
# QueryServer("addContact", payload0=json)




