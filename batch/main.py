import src.setting as setting
import src.db as db
import src.main as mecab


def main():
    print("Hello world")
    if setting.ENV == "pro":
        print("ENV: 本番環境")
    else:
        print("ENV: 開発環境")
    d = db.DB()
    d.test()
    d.close()
    m = mecab.getMecab()
    m.parse("鬼滅の刃")


if __name__ == "__main__":
    main()
