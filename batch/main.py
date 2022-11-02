from src.main import daySpeech
import src.setting as setting
import src.db as db


def main():
    date = "2022-03-01"
    it = daySpeech(date)
    last = "1"
    s = next(it, None)
    while (s):
        last = s
        s = next(it, None)
    print(last)


if __name__ == "__main__":
    main()
