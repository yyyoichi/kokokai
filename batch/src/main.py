import re
import MeCab
from os.path import join, dirname


def getMecab():
    path = join(dirname(__file__), '..\mecab-ipadic-neologd')
    dict = "-d " + re.sub(r'\\', '/', path)
    return MeCab.Tagger(dict)


if __name__ == "__main__":
    print(getMecab().parse("約束のネバーランドが面白い"))
