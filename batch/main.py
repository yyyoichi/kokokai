import src.setting.setting as setting


def main():
    print("Hello world")
    if setting.ENV == "pro":
        print("ENV: 本番環境")
    else:
        print("ENV: 開発環境")


if __name__ == "__main__":
    main()
