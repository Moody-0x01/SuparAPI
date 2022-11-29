from requests import post

api = "http://localhost:8888/login"

def Post(data: dict) -> bytes:
    return post(api, json=data).json()


def  main():
    data = {
        "Password": "1234",
        "Email": "ex@gmail.com"
    }

    print(Post(data))


if __name__ == "__main__": main()







