import datetime
from src.sentence import daySpeech
from src.mcb import getMecab, Morpheme


class Validation:
    stop_words = [""]

    def __init__(self, morpheme: Morpheme) -> None:
        self.m = morpheme

    def is_stop_word(self):
        return self.m.surface() in self.stop_words

    def is_noun(self):
        return self.m.part_of_speech() == "名詞"

    def is_int(self):
        return self.m.part_of_speech_details_1() == "数"

    def is_asterisk(self):
        return self.m.prototype() == "*"


with open('stop_words.txt', 'r', encoding="utf-8") as f:
    Validation.stop_words = [w.strip() for w in f.readlines()]


def main():
    """
    その日の共起リストをDBに格納する
    """
    m = getMecab()
    date = datetime.datetime.now(datetime.timezone(
        datetime.timedelta(hours=9))) - datetime.timedelta(days=10)
    speech = daySpeech(date.strftime("%Y-%m-%d"))
    sentence = next(speech, None)
    # 名詞の単語リスト
    noun_list = []

    def add_noun_list(mph: Morpheme):
        v = Validation(mph)
        if v.is_noun() and not v.is_asterisk() and not v.is_int() and not v.is_stop_word():
            noun_list.append(mph.prototype())

    while (sentence):
        p = m.parse(sentence)
        morpheme = next(p, None)
        while (morpheme):
            add_noun_list(morpheme)
            morpheme = next(p, None)
        sentence = next(speech, None)


if __name__ == "__main__":
    main()
