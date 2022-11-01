import re
import MeCab
from os.path import join, dirname
import minute.main as minute
import setting.setting as setting


def getMecab():
    path = join(dirname(__file__), '..\mecab-ipadic-neologd')
    dict = "-d " + re.sub(r'\\', '/', path)
    return MeCab.Tagger(dict)


m = getMecab()


def parse(str: str):
    return m.parse(str)


if __name__ == "__main__":
    print(parse("約束のネバーランドが面白い"))
