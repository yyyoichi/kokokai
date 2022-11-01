import MeCab
import minute.main as minute


m = MeCab.Tagger("-Ochasen")


def parse(str: str):
    return m.parse(str)


if __name__ == "__main__":
    print(parse("すもももももももものうち"))
