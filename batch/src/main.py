import MeCab
import minute.main as minute
import setting.setting as setting

print(setting.DICT_PATH)
m = MeCab.Tagger("-d " + setting.DICT_PATH)


def parse(str: str):
    return m.parse(str)


if __name__ == "__main__":
    print(parse("約束のネバーランドが面白い"))
