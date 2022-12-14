import MeCab
import re
from os.path import join, dirname


def getMecab():
    path = join(dirname(__file__), '..\mecab-ipadic-neologd')
    dict = "-d " + re.sub(r'\\', '/', path)
    return Parser(MeCab.Tagger(dict))


class Morpheme:
    def __init__(self, line: str) -> None:
        l = line.split("\t")
        self.s = l[0]
        self.data = l[1].split(",")

    def surface(self) -> str:
        """
        表層形
        """
        return self.s

    def part_of_speech(self) -> str:
        """
        品詞
        """
        return self.data[0]

    def part_of_speech_details_1(self) -> str:
        """
        品詞細分類1
        """
        return self.data[1]

    def part_of_speech_details_2(self) -> str:
        """
        品詞細分類2
        """
        return self.data[2]

    def part_of_speech_details_3(self) -> str:
        """
        品詞細分類3
        """
        return self.data[3]

    def inflection(self) -> str:
        """
        活用型
        (一段など)
        """
        return self.data[4]

    def inflected_form(self) -> str:
        """
        活用形
        (基本形など)
        """
        return self.data[5]

    def prototype(self) -> str:
        """
        原形
        """
        return self.data[6]

    def reading(self) -> str:
        """
        読み
        """
        return self.data[7]

    def pronunciation(self) -> str:
        """
        発音
        """
        return self.data[8]


class Parser:
    def __init__(self, mecab: MeCab.Tagger) -> None:
        self.m = mecab

    def parse(self, str: str):
        c = self.m.parse(str)
        lines = c.split("\n")
        for line in lines:
            if line == "EOS":
                break
            yield Morpheme(line)


if __name__ == "__main__":
    m = getMecab()
    sentence = m.parse("山本洋一郎と申します。")
    morpheme = next(sentence, None)
    while (morpheme):
        print(morpheme.prototype())
        morpheme = next(sentence, None)
