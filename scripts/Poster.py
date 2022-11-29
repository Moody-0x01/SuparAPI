from requests import post

api = "http://localhost:8888/Test"

def Post(data: dict) -> dict:
    return post(api, json=data).json()


def  main():

    data = {
        "Email": "Email@gmail.com",
        "Password": "1234567"
    }
    

    print(Post(data))


if __name__ == "__main__": main()







